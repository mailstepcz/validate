package validate

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/mailstepcz/testutils/testcond"
	"github.com/stretchr/testify/require"
)

type Person struct {
	Name Required[string] `json:"name"`
	Age  Required[int]    `json:"age"`
	Type Required[string] `json:"type" enums:"ADMIN,USER"`
}

type Address struct {
	ZipCode Required[string] `json:"zipCode"`
}

type PersonNested struct {
	Address Required[*Address] `json:"address"`
}

type PersonNested2 struct {
	Address *Address `json:"address"`
}

type PersonNested3 struct {
	Address Address `json:"address"`
}

type PersonNested4 struct {
	Address Required[Address] `json:"address"`
}

type PersonNested5 struct {
	Addresses Required[[]*Address] `json:"addresses"`
}

type PersonNested6 struct {
	Addresses Required[[]Address] `json:"addresses"`
}

type PersonNested7 struct {
	Addresses []*Address `json:"addresses"`
}

type PersonNested8 struct {
	Addresses []Address `json:"addresses"`
}

func TestRequired(t *testing.T) {
	req := require.New(t)

	var p Person
	err := json.Unmarshal([]byte(`{"name":"Saoirse","type":"ADMIN"}`), &p)
	req.Nil(err)
	err = Struct(&p)
	req.NotNil(err)
	testcond.Equal(t, "field 'Age' in 'validate.Person' is required", err.Error())

	err = json.Unmarshal([]byte(`{"name":"Saoirse","age":25,"type":"ADMIN"}`), &p)
	req.Nil(err)
	err = Struct(&p)
	req.Nil(err)
}

func TestRequiredPtr(t *testing.T) {
	req := require.New(t)

	var p Person
	err := json.Unmarshal([]byte(`{"name":"Saoirse","type":"ADMIN"}`), &p)
	req.NoError(err)

	testcond.Equal(t, p.Name.Ptr().(*string), (*string)(p.Name.UnsafePtr()))
	testcond.Equal(t, reflect.ValueOf(p.Name.Ptr()).UnsafePointer(), p.Name.UnsafePtr())
}

func TestRequiredNestedSuccess(t *testing.T) {
	req := require.New(t)

	var p PersonNested
	err := json.Unmarshal([]byte(`{"address": {"zipCode": "111222"}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Address.Value().(*Address).ZipCode.value)
}

func TestRequiredNested2Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested2
	err := json.Unmarshal([]byte(`{"address": {"zipCode": "111222"}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Address.ZipCode.value)
}

func TestRequiredNested3Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested3
	err := json.Unmarshal([]byte(`{"address": {"zipCode": "111222"}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Address.ZipCode.value)
}

func TestRequiredNested4Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested4
	err := json.Unmarshal([]byte(`{"address": {"zipCode": "111222"}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Address.Value().(Address).ZipCode.value)
}

func TestRequiredNested5Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested5
	err := json.Unmarshal([]byte(`{"addresses": [{"zipCode": "111222"}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Addresses.Value().([]*Address)[0].ZipCode.value)
}

func TestRequiredNested6Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested6
	err := json.Unmarshal([]byte(`{"addresses": [{"zipCode": "111222"}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Addresses.Value().([]Address)[0].ZipCode.value)
}

func TestRequiredNested7Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested7
	err := json.Unmarshal([]byte(`{"addresses": [{"zipCode": "111222"}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Addresses[0].ZipCode.value)
}

func TestRequiredNested8Success(t *testing.T) {
	req := require.New(t)

	var p PersonNested8
	err := json.Unmarshal([]byte(`{"addresses": [{"zipCode": "111222"}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Nil(err)

	testcond.Equal(t, "111222", p.Addresses[0].ZipCode.value)
}

func TestRequiredNestedError(t *testing.T) {
	req := require.New(t)

	var p PersonNested
	err := json.Unmarshal([]byte(`{"address": {}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested2Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested2
	err := json.Unmarshal([]byte(`{"address": {}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested3Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested3
	err := json.Unmarshal([]byte(`{"address": {}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested4Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested4
	err := json.Unmarshal([]byte(`{"address": {}}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested5Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested5
	err := json.Unmarshal([]byte(`{"addresses": [{}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested6Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested6
	err := json.Unmarshal([]byte(`{"addresses": [{}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested7Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested7
	err := json.Unmarshal([]byte(`{"addresses": [{}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredNested8Error(t *testing.T) {
	req := require.New(t)

	var p PersonNested8
	err := json.Unmarshal([]byte(`{"addresses": [{}]}`), &p)
	req.Nil(err)

	err = Struct(&p)
	req.Error(err)
}

func TestRequiredEnums(t *testing.T) {
	req := require.New(t)

	var p Person
	err := json.Unmarshal([]byte(`{"name":"Saoirse","age":25,"type":"ADMIN"}`), &p)
	req.NoError(err)
	err = Struct(&p)
	req.NoError(err)

	err = json.Unmarshal([]byte(`{"name":"Saoirse","age":25,"type":"SUPER_ADMIN"}`), &p)
	req.NoError(err)
	err = Struct(&p)
	req.Error(err)
	req.Contains(err.Error(), "invalid enum value 'SUPER_ADMIN' for field 'Type', allowed values are 'ADMIN,USER'")

}

var gr interface{}

func BenchmarkWithValidation(b *testing.B) {
	var lr interface{}
	bs := []byte(`{"name":"Saoirse","age":25,"type":"ADMIN"}`)
	for i := 0; i < b.N; i++ {
		var p Person
		if err := json.Unmarshal(bs, &p); err != nil {
			b.Fatal(err)
		}
		if err := Struct(&p); err != nil {
			b.Fatal(err)
		}
		lr = &p
	}
	gr = lr
}

func BenchmarkWithoutValidation(b *testing.B) {
	var lr interface{}
	bs := []byte(`{"name":"Saoirse","age":25,"type":"ADMIN"}`)
	for i := 0; i < b.N; i++ {
		var p Person
		if err := json.Unmarshal(bs, &p); err != nil {
			b.Fatal(err)
		}
		lr = &p
	}
	gr = lr
}
