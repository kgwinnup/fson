
# Basic Usage


```go
import (
    "fmt"
    "github.com/kgwinnup/fson"
)

func main() {
	data := []byte("{\"foo\": 1, \"foo2\": {\"bar\": 1, \"baz\": [1,1,1]}}")

	fsonobj := New()
	if err := fsonobj.Loads(&data); err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsonobj)
	fsonobj.Fmap(func(v interface{}) interface{} {
		switch v.(type) {
		case float64:
			return v.(float64) + 1
		default:
			return v
		}
	})
	fmt.Println(fsonobj)

	data = []byte("{\"boo\": true, \"hello\": \"world\", \"obj\": {\"foo\": \"bar\"}, \"baz\": [400,2,3]}")
	if err := fsonobj.Loads(&data); err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsonobj)

	fsonobj.Insert([]string{"new", "field"}, "hello world", false)
	fmt.Println(fsonobj)

	newfield, _ := fsonobj.Get([]string{"new", "field"})
	fmt.Println(newfield.(string))

}

```

## output

```
{"foo":1,"foo2":{"bar":1,"baz":[1,1,1]}}
{"foo":2,"foo2":{"bar":2,"baz":[2,2,2]}}
{"foo":2,"foo2":{"bar":2,"baz":[2,2,2]},"boo":true,"hello":"world","obj":{"foo":"bar"},"baz":[400,2,3]}
{"foo":2,"foo2":{"bar":2,"baz":[2,2,2]},"boo":true,"hello":"world","obj":{"foo":"bar"},"baz":[400,2,3],"new":{"field":"hello world"}}
hello world
```
