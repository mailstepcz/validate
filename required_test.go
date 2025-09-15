package validate

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Person struct {
	Name Required[string] `json:"name"`
	Age  Required[int]    `json:"age"`
	Type Required[string] `json:"type" enums:"ADMIN,USER"`
}

func TestRequired(t *testing.T) {
	req := require.New(t)

	var p Person
	err := json.Unmarshal([]byte(`{"name":"Saoirse","type":"ADMIN"}`), &p)
	req.Nil(err)
	err = Struct(&p)
	req.NotNil(err)
	req.Equal("field 'Age' in 'validate.Person' is required", err.Error())

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

	req.Equal(p.Name.Ptr().(*string), (*string)(p.Name.UnsafePtr()))
	req.Equal(reflect.ValueOf(p.Name.Ptr()).UnsafePointer(), p.Name.UnsafePtr())
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
