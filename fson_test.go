package fson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func testDataTypeToString(t *testing.T, data []byte, expect string) {
	out := New(data)
	if fmt.Sprintf("%v", out) != expect {
		t.Errorf("Error dumping data type to string\nExpected: '%v'\nGot: '%v'\n", expect, out)
	}
}

func TestString(t *testing.T) {

	data0 := []byte("{\"boo\":\"baz\"}")
	data1 := []byte("{\"boo\":10}")
	data2 := []byte("{\"boo\":true}")
	data3 := []byte("{\"boo\":[false]}")
	data4 := []byte("{\"boo\":{\"baz\":false}}")

	testDataTypeToString(t, data0, "{\"boo\":\"baz\"}")
	testDataTypeToString(t, data1, "{\"boo\":10}")
	testDataTypeToString(t, data2, "{\"boo\":true}")
	testDataTypeToString(t, data3, "{\"boo\":[false]}")
	testDataTypeToString(t, data4, "{\"boo\":{\"baz\":false}}")

}

func TestSet(t *testing.T) {
	data := []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\"}, \"baz\": [400,2,3]}")
	out := New(data)

	if val, err := out.Get([]string{"boo"}); err != nil {
		t.Errorf("Error retrieving bool value: %v", val)
	} else {
		if val.(bool) != true {
			t.Errorf("error parsing bool value")
		}
	}

	out.Set([]string{"once", "twice", "third"}, 100, false)

	if val, err := out.Get([]string{"once", "twice", "third"}); err != nil {
		t.Errorf("error retrieving value")
	} else {
		if val != 100 {
			t.Errorf("error converting value to 100")
		}
	}

	out.Set([]string{"once", "twice", "third"}, 200, true)

	if val, err := out.Get([]string{"once", "twice", "third"}); err != nil {
		t.Errorf("error retreiving value")
	} else {
		switch val.(type) {
		case []interface{}:
			if len(val.([]interface{})) != 2 {
				t.Errorf("array has unexpected len")
			}

		default:
			t.Errorf("should be an array")
		}
	}
}

func TestFmap(t *testing.T) {
	data := []byte("{\"foo\": 1, \"foo2\": {\"bar\": 1, \"baz\": [1,1,1]}}")
	out := New(data)

	out.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return v.(float64) + 1
		default:
			return v.(int) + 1
		}
	})

	if val, err := out.Get([]string{"foo"}); err != nil {
		t.Errorf("error getting foo value")
	} else {
		switch val.(type) {
		case float64:
			if val.(float64) != 2 {
				t.Errorf("error")
			}
		}
	}
}

func TestFmap2(t *testing.T) {
	data := []byte("{\"foo\": 1, \"foo2\": {\"bar\": 1, \"baz\": {\"v\": 1, \"vv\": {\"vvv\": 1}}}}")
	out := New(data)

	out.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return v.(float64) + 1
		default:
			return v.(int) + 1
		}
	})

	if val, err := out.Get([]string{"foo2", "baz", "vv", "vvv"}); err != nil {
		t.Errorf("error getting foo value")
	} else {
		switch val.(type) {
		case float64:
			if val.(float64) != 2 {
				t.Errorf("error")
			}
		}
	}
}

type TestStruct struct {
	Name    string
	ObjData int
	ObjName string
}

func TestJSONMarshal(t *testing.T) {
	ts := TestStruct{"name", 10, "obj name"}
	tsBytes := new(bytes.Buffer)
	json.NewEncoder(tsBytes).Encode(ts)

	out := New(tsBytes.Bytes())

	if b, err := json.Marshal(out); err != nil {
		t.Errorf("error marshaling json")
	} else {
		if b2, err := json.Marshal(ts); err != nil {
			t.Errorf("error marshling test struct")
		} else {
			if len(b) != len(b2) {
				t.Errorf("error matching marshals")
			}
		}
	}
}
