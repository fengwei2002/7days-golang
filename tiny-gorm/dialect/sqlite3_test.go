package dialect

import (
	"reflect"
	"testing"
)

// 测试是否可以正常通过反射得到 value
func TestDataTypeOf(t *testing.T) {
	dial := &sqlite3{}
	cases := []struct {
		Value interface{}
		Type  string
	}{
		{"Tom", "text"},
		{123, "integer"},
		{1.2, "real"},
		{[]int{1, 2, 3}, "blob"},
	}

	for _, c := range cases {
		if typ := dial.DataTypeOf(reflect.ValueOf(c.Value)); typ != c.Type {
			t.Fatalf("expect %s, but got %s", c.Type, typ)
		}
	}
}
