package main

import "fmt"
import "encoding/json"

type Fson struct {
	data map[string]interface{}
}

func New() *Fson {
	data := make(map[string]interface{})
	return &Fson{data}
}

func (self *Fson) Loads(b *[]byte) error {
	if err := json.Unmarshal(*b, &self.data); err != nil {
		return err
	}
	return nil
}

func (self *Fson) toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return "\"" + value.(string) + "\""
	case bool:
		return fmt.Sprintf("%v", value.(bool))
	case []interface{}:
		ret := "["
		for i, item := range value.([]interface{}) {
			ret += self.toString(item)
			if i < len(value.([]interface{}))-1 {
				ret += ","
			}
		}
		ret += "]"
		return ret
	case map[string]interface{}:
		ret := "{"
		size := len(value.(map[string]interface{}))
		count := 0
		for k, v := range value.(map[string]interface{}) {
			ret += "\"" + k + "\""
			ret += ":"
			ret += self.toString(v)
			if count < size-1 {
				ret += ","
			}
			count += 1
		}
		ret += "}"
		return ret
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (self Fson) String() string {
	ret := "{"
	size := len(self.data)
	count := 0
	for k, v := range self.data {
		ret += "\"" + k + "\":"
		ret += self.toString(v)
		if count < size-1 {
			ret += ","
		}
		count += 1
	}
	ret += "}"
	return ret
}

func (self *Fson) insert(path []string, value interface{}, cur map[string]interface{}, appendList bool) {
	// check if we are where the insert should happen
	if len(path) == 1 {
		if appendList {
			//check if something is already there
			if _, ok := cur[path[0]]; ok {
				switch cur[path[0]].(type) {
				case []interface{}:
					// add new value to end of current list
					cur[path[0]] = append(cur[path[0]].([]interface{}), value)
				default:
					// it is some singleton value, create a list of it and append new value to it
					temp := cur[path[0]]
					cur[path[0]] = make([]interface{}, 0, 2)
					cur[path[0]] = append(cur[path[0]].([]interface{}), temp)
					cur[path[0]] = append(cur[path[0]].([]interface{}), value)
				}
			}
		} else {
			cur[path[0]] = value
		}
	} else {
		if _, ok := cur[path[0]]; !ok {
			cur[path[0]] = make(map[string]interface{})
		}
		self.insert(path[1:], value, cur[path[0]].(map[string]interface{}), appendList)
	}
}

func (self *Fson) Insert(path []string, value interface{}, appendList bool) {
	self.insert(path, value, self.data, appendList)
}

func (self *Fson) get(path []string, cur map[string]interface{}) interface{} {
	if len(path) == 1 {
		if _, ok := cur[path[0]]; ok {
			return cur[path[0]]
		} else {
			return nil
		}
	}

	if _, ok := cur[path[0]].(map[string]interface{}); ok {
		return self.get(path[1:], cur[path[0]].(map[string]interface{}))
	} else {
		return nil
	}
}

func (self *Fson) Get(path []string) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("No path specified")
	}

	if v := self.get(path, self.data); v == nil {
		return nil, fmt.Errorf("path not found: %v", path)
	} else {
		return v, nil
	}
}

func main() {

}
