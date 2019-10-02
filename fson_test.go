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

func TestGetF(t *testing.T) {
	data := []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\"}, \"baz\": [400,2,3]}")
	out := New(data)

	if obj, ok := out.GetF("obj"); !ok {
		t.Errorf("error retrieving object with GetF")
	} else {
		if _, ok := obj.Get("foo"); !ok {
			t.Errorf("error getting value from Fson object returned by GetF")
		}
	}
}

func TestSet(t *testing.T) {
	data := []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\"}, \"baz\": [400,2,3]}")
	out := New(data)

	if val, ok := out.Get("boo"); !ok {
		t.Errorf("Error retrieving bool value: %v", val)
	} else {
		if val.(bool) != true {
			t.Errorf("error parsing bool value")
		}
	}

	out.Set(100, "once", "twice", "third")

	if val, ok := out.Get("once", "twice", "third"); !ok {
		t.Errorf("error retrieving value")
	} else {
		if val != 100 {
			t.Errorf("error converting value to 100")
		}
	}

	out.SetA(200, "once", "twice", "third")

	if val, ok := out.Get("once", "twice", "third"); !ok {
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

func TestGet(t *testing.T) {
	data := []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\"}, \"baz\": [400,2,3]}")
	out := New(data)

	if d, ok := out.Get("obj", "foo"); !ok {
		t.Errorf("error getting value")
	} else if d != "bar" {
		t.Errorf("invalid key fetched with GetP")
	}

	if d, ok := out.Get("obj", "foo"); !ok {
		t.Errorf("error getting value")
	} else if d != "bar" {
		t.Errorf("invalid key fetched with GetP")
	}

}

func TestDel(t *testing.T) {
	data := []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\", \"foo2\": \"bar2\"}, \"baz\": [400,2,3]}")
	out := New(data)

	out.Del([]string{"obj", "foo"})

	if _, ok := out.Get("obj", "foo"); ok {
		t.Errorf("error deleting key")
	}
}

func TestFmap(t *testing.T) {
	data := []byte(`{
		"foo": 1, 
		"foo2": {
			"bar": 1, 
			"baz": [1,1,1]
			}
		}`)
	out := New(data)

	out.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return v.(float64) + 1
		default:
			return v.(int) + 1
		}
	})

	if val, ok := out.Get("foo"); !ok {
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

func TestFilter(t *testing.T) {
	data := []byte(`{
		"foo": 2, 
		"foo2": {
			"bar": 1, 
			"baz": [2,1,1]
			}
		}`)
	out := New(data)

	out.Filter(func(v interface{}) bool {
		switch v.(type) {
		case float64:
			if v.(float64) > 1 {
				return false
			} else {
				return true
			}
		default:
			return true
		}
	})

	if d, ok := out.Get("foo2", "baz"); !ok {
		t.Errorf("error getting value")
	} else if len(d.([]interface{})) != 2 {
		t.Errorf("Filter did not reduce the list values")
	}

}

func TestFmap2(t *testing.T) {
	data := []byte("{\"foo\": 1, \"foo2\": {\"bar\": 1, \"baz\": {\"v\": 1, \"vv\": {\"vvv\": 1}}}}")
	out := New(data)

	out.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return v.(float64) + 1
		}
		return v
	})

	if val, ok := out.Get("foo2", "baz", "vv", "vvv"); !ok {
		t.Errorf("error getting foo value")
	} else {
		switch val.(type) {
		case float64:
			if val.(float64) != 2 {
				t.Errorf("error")
			}
		}
	}

	out.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return fmt.Sprintf("%v", v.(float64))
		}
		return v
	})

	if val, ok := out.Get("foo2", "baz", "vv", "vvv"); !ok {
		t.Errorf("error getting foo value")
	} else {
		switch val.(type) {
		case string:
			if val.(string) != "2" {
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

func TestMerge(t *testing.T) {
	data := []byte(`{
		"foo": 1, 
		"foo2": {
			"bar": 1, 
			"baz": [1,1,1]
			}
		}`)
	data2 := []byte(`{
		"boo": 1, 
		"foo2": {
			"bar": 1, 
			"faz": "blah"
			}
		}`)

	out := New(data)

	out.Merge(New(data2))

	if val, ok := out.Get("boo"); !ok {
		t.Errorf("error getting foo value")
	} else {
		switch val.(type) {
		case float64:
			if val.(float64) != 1 {
				t.Errorf("error")
			}
		}
	}

	if val, ok := out.Get("foo2", "faz"); !ok {
		t.Errorf("error getting foo2/faz value")
	} else {
		switch val.(type) {
		case string:
			if val.(string) != "blah" {
				t.Errorf("error")
			}
		}
	}

}

func TestMerge2(t *testing.T) {
	data := []byte(`{
		"foo": 1, 
		"foo2": {
			"bar": 1, 
			"faz": [1,1,1]
			}
		}`)
	data2 := []byte(`{
		"boo": 1, 
		"foo2": {
			"bar": 1, 
			"faz": "blah"
			}
		}`)

	out := New(data)

	out.Merge(New(data2))

	if val, ok := out.Get("boo"); !ok {
		t.Errorf("error getting foo value")
	} else {
		switch val.(type) {
		case float64:
			if val.(float64) != 1 {
				t.Errorf("error")
			}
		default:
			t.Errorf("wat")
		}
	}

	if val, ok := out.Get("foo2", "faz"); !ok {
		t.Errorf("error getting foo2/faz value")
	} else {
		switch val.(type) {
		case string:
			if val.(string) != "blah" {
				t.Errorf("error")
			}
		default:
			t.Errorf("wat")
		}
	}

}
