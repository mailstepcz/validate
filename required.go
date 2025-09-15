// Package validate provides a special decorative type which indicates that a field in a request is required,
// that is, it must have a value in the incoming request (which might be `null` for some types but must not be omitted).
//
// The type supports JSON unmarshalling provided the underlying type can be unmarshalled from a slice of byte
// into an instance.
package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unsafe"
)

// Required is a decorative type which signifies that a field in a structure is required.
// It should only be used in structures that represent request bodies.
// The type supports custom unmarshalling from JSON.
// Furthermore the keyvalue copier can handle this type provided it figures in the source.
type Required[T any] struct {
	value T
	valid bool
}

func (r *Required[T]) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &r.value); err != nil {
		return err
	}
	r.valid = true
	return nil
}

func (r *Required[T]) String() string {
	if r.valid {
		return fmt.Sprintf("%v", r.value)
	}
	return "N/A"
}

// HasValue returns true if the underlying value has been unmarshalled into.
func (r *Required[T]) HasValue() bool { return r.valid }

// Value returns the underlying value.
func (r *Required[T]) Value() interface{} { return r.value }

// Ptr returns the pointer to the underlying value.
func (r *Required[T]) Ptr() interface{} { return &r.value }

// UnsafePtr returns the unsafe pointer to the underlying value.
func (r *Required[T]) UnsafePtr() unsafe.Pointer { return unsafe.Pointer(&r.value) }

// RequiredType returns the type of the underlying value.
func (r *Required[T]) RequiredType() reflect.Type {
	return reflect.TypeFor[T]()
}

// SetValid marks the instance as valid, that is, containing a value.
func (r *Required[T]) SetValid(v bool) { r.valid = true }

// SettableValue returns the settable (reflection) value of the underlying value.
func (r *Required[T]) SettableValue() reflect.Value { return reflect.ValueOf(&r.value).Elem() }

// RequiredIface is the interface without type parameters providing access to the [Required] type constructor.
type RequiredIface interface {
	HasValue() bool
	Value() interface{}
	Ptr() interface{}
	UnsafePtr() unsafe.Pointer
	RequiredType() reflect.Type
	SetValid(bool)
	SettableValue() reflect.Value
}

var (
	// RequiredIfaceType is the type of [RequiredIface].
	RequiredIfaceType = reflect.TypeFor[RequiredIface]()
	// ErrBadType indicates that the provided argument is ill-typed.
	ErrBadType = errors.New("bad type")

	_ json.Unmarshaler = (*Required[int])(nil)
	_ RequiredIface    = (*Required[int])(nil)
)

// Struct validates the provided argument which must be a pointer to a structure.
// Any fields whose type is [Required] are checked.
// Any fields which has enums tag, then value is checked.
// The returned error is a multi-error containing the errors emitted for all misbehaving fields.
//
// Struct panics if the argument is ill-typed.
func Struct(x interface{}) error {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("%w: %T", ErrBadType, x)
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %T", ErrBadType, x)
	}
	var errs error
	for _, f := range reflect.VisibleFields(v.Type()) {
		if !f.IsExported() {
			continue
		}
		fv := v.FieldByIndex(f.Index)
		if x, ok := fv.Addr().Interface().(RequiredIface); ok {
			if !x.HasValue() {
				errs = errors.Join(errs, fmt.Errorf("field '%s' in '%s' is required", f.Name, v.Type()))
			}

			if t := x.RequiredType(); t.Kind() == reflect.Struct {
				errs = errors.Join(errs, Struct(x.Ptr()))
			} else if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
				errs = errors.Join(errs, Struct(x.Value()))
			} else if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Pointer && t.Elem().Elem().Kind() == reflect.Struct {
				v := reflect.ValueOf(x.Value())
				for i := 0; i < v.Len(); i++ {
					errs = errors.Join(errs, Struct(v.Index(i).Interface()))
				}
			} else if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct {
				v := reflect.ValueOf(x.Value())
				for i := 0; i < v.Len(); i++ {
					errs = errors.Join(errs, Struct(v.Index(i).Addr().Interface()))
				}
			}
		} else if f.Type.Kind() == reflect.Struct {
			errs = errors.Join(errs, Struct(fv.Addr().Interface()))
		} else if f.Type.Kind() == reflect.Pointer && f.Type.Elem().Kind() == reflect.Struct {
			errs = errors.Join(errs, Struct(fv.Interface()))
		} else if f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Pointer && f.Type.Elem().Elem().Kind() == reflect.Struct {
			for i := 0; i < fv.Len(); i++ {
				errs = errors.Join(errs, Struct(fv.Index(i).Interface()))
			}
		} else if f.Type.Kind() == reflect.Slice && f.Type.Elem().Kind() == reflect.Struct {
			for i := 0; i < fv.Len(); i++ {
				errs = errors.Join(errs, Struct(fv.Index(i).Addr().Interface()))
			}
		}

	}
	return errs
}

func validateEnumValue(tag, value, filedName string) error {
	allowedValues := strings.Split(tag, ",")
	if !slices.Contains(allowedValues, value) {
		return fmt.Errorf("invalid enum value '%s' for field '%s', allowed values are '%s'", value, filedName, tag)
	}
	return nil
}
