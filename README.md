
# Fson

[![godoc for kgwinnup/fson][godoc-badge]][godoc-url]

Fson is a simple library for working with arbitrary JSON data of unknown
structure. One of the key reasons for creating this were to add two critical
features for my use cases. First was an Fmap function for easily applying a
function to the JSON data structure and secondly, provide a interface type for
dealing with Postgresql's JSON(b) types and that types encoding to JSON for our
REST api.

# Usage

Import package

```go
import (
	"github.com/kgwinnup/fson"
)
```

Basic example of interacting with a JSON structure

```go
package main

import (
	"fmt"
	"github.com/kgwinnup/fson"
)

func main() {
	//receive some byte array and create an Fson object
	data := []byte(`{
        "foo": 1, 
        "foo2": { 
                    "bar": 1, 
                    "baz": [1,1,1]
                }
        }`)
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

