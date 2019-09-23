package main

import (
	"fmt"
	"testing"
)

func testDataTypeToString(t *testing.T, data []byte, expect string) {
	out := New()
	if err := out.Loads(&data); err != nil {
		t.Errorf("%v", err)
	}

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

func TestInsertion(t *testing.T) {
	data := []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\"}, \"baz\": [400,2,3]}")

	out := New()
	if err := out.Loads(&data); err != nil {
		t.Errorf("%v", err)
	}

	if val, err := out.Get([]string{"boo"}); err != nil {
		t.Errorf("Error retrieving bool value: %v", val)
	} else {
		if val.(bool) != true {
			t.Errorf("error parsing bool value")
		}
	}

	out.Insert([]string{"once", "twice", "third"}, 100, false)

	if val, err := out.Get([]string{"once", "twice", "third"}); err != nil {
		t.Errorf("error retrieving value")
	} else {
		if val != 100 {
			t.Errorf("error converting value to 100")
		}
	}

	out.Insert([]string{"once", "twice", "third"}, 200, true)

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

	out := New()
	if err := out.Loads(&data); err != nil {
		t.Errorf("%v", err)
	}

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
				fmt.Println(val)
				t.Errorf("error")
			}
		}
	}
}
