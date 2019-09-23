
# Fson

Fson is a simple library for working with arbitrary JSON data of unknown
structure. Additionally, the JSON and Postgres's JSON(b) interface types are
supported.


# Basic Usage

```go
package main

import (
	"fmt"
	"github.com/kgwinnup/fson"
)

func main() {
	//receive some byte array and create an Fson object
	data := []byte("{\"foo\": 1, \"foo2\": {\"bar\": 1, \"baz\": [1,1,1]}}")
	fsonobj := fson.New(data)

	// print the JSON string
	fmt.Println(fsonobj)

	// Get a value from the JSON
	if val, err := fsonobj.Get([]string{"foo2", "bar"}); err != nil {
		fmt.Println(err)
	} else {

		//add 100 to our new value
		newVal := val.(float64) + 100

		//set the new value back in the JSON structure
		// the last bool parameter, if true, will turn the value in the key to an array
		fsonobj.Set([]string{"foo2", "bar"}, newVal, false)
		fmt.Println(fsonobj)
	}

	// -------------------------------------------------------------------------------
	// -------------------------------------------------------------------------------

	// lets incremenet all number values by 1
	fsonobj.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return v.(float64) + 1
		default:
			return v
		}
	})
	fmt.Println(fsonobj)

}
```

