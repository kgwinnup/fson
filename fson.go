// Package fson provides a way to interact with arbitrary JSON fields as well as use JSON(b) within Postgresql
// Fson supports the json interfaces MarshalJSON/UnmarshalJSON as well as the Scan interface within the db packages
package fson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Fson struct is the core structure, no exported members
type Fson struct {
	Data map[string]interface{}
}

// MarshalJSON implementation for JSON encoding
func (f *Fson) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Data)
}

// UnmarshalJSON impelmentation for JSON decoding to Fson
func (f *Fson) UnmarshalJSON(b []byte) error {
	return f.Loads(b)
}

// Scan will take a pointer interface and attempt to decode as Fson structure
func (f *Fson) Scan(src interface{}) error {
	return f.Loads(src.([]byte))
}

// Bytes will return a byte slice of the underlying JSON data. If there is an
// error an empty byte slice is returned
func (f *Fson) Bytes() []byte {
	if b, err := json.Marshal(f.Data); err == nil {
		return b
	}
	return []byte{}
}

// String converts the Fson object into a string
func (f *Fson) String() string {
	return string(f.Bytes())
}

// New initializes a new Fson object given a JSON byte slice
func New(b []byte) *Fson {
	if b == nil {
		return &Fson{make(map[string]interface{})}
	} else {
		f := &Fson{}
		json.Unmarshal(b, &f.Data)
		return f
	}
}

// Loads will take nil or a []byte array and create a Fson object with it.
// Generally nil is used if you want to build your JSON object from scratch with the Set method
// This will overwrite any current data that is in the Fson object
func (f *Fson) Loads(b []byte) error {
	f.Data = make(map[string]interface{})
	if b == nil {
		return nil
	}
	if err := json.Unmarshal(b, &f.Data); err != nil {
		return err
	}
	return nil
}

// FromFile will load a JSON object from a file
func FromFile(path string) (*Fson, error) {
	data, err := ioutil.ReadFile("/tmp/dat")
	if err != nil {
		return New(data), nil
	}

	return nil, err
}

// ParseJSON takes some bytes and parses it into an Fson object
func ParseJSON(b []byte) (*Fson, error) {
	if b == nil {
		return nil, fmt.Errorf("no bytes provided to parse as JSON")
	}

	f := &Fson{}
	if err := json.Unmarshal(b, &f.Data); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Fson) set(path []string, value interface{}, cur map[string]interface{}, appendList bool) {
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
		f.set(path[1:], value, cur[path[0]].(map[string]interface{}), appendList)
	}
}

// Set will set a value to a specified path. The path is defined as a variadic
// parameter. This will overwrite any value that is located at the provided key
func (f *Fson) Set(value interface{}, path ...string) {
	f.set(path, value, f.Data, false)
}

// SetA is the same as set except that it will append to a list. If no list
// exists, it will create a list. If an element currently exists, a new list
// will be created with the existing element as the head of the list.
func (f *Fson) SetA(value interface{}, path ...string) {
	f.set(path, value, f.Data, true)
}

func (f *Fson) get(path []string, cur map[string]interface{}) interface{} {
	if len(path) == 1 {
		if _, ok := cur[path[0]]; ok {
			return cur[path[0]]
		} else {
			return nil
		}
	}

	if _, ok := cur[path[0]].(map[string]interface{}); ok {
		return f.get(path[1:], cur[path[0]].(map[string]interface{}))
	}

	return nil
}

// Get works like set in that the path to the key is specified as a list of
// string values which represent the orderd nested object keys
func (f *Fson) Get(path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}

	v := f.get(path, f.Data)
	if v == nil {
		return nil, false
	}

	return v, true
}

// GetObject will work the same as Get except this will automatically encode the
// result into an Fson object. If the value retrieved is not a JSON object,
// this will return nil, false
func (f *Fson) GetObject(path ...string) (*Fson, bool) {
	if len(path) == 0 {
		return nil, false
	}

	v := f.get(path, f.Data)
	if v == nil {
		return nil, false
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}

	return New(b), true
}

// Exists will return true if the key exists in the JSON
func (f *Fson) Exists(path ...string) bool {
	_, ok := f.Get(path...)
	return ok
}

// GetArray will retrieve an array item from a specific key in the Fson object.
// The JSON array will be returned as a slice of interface{}
func (f *Fson) GetArray(path ...string) ([]interface{}, error) {
	data, ok := f.Get(path...)
	if !ok {
		return nil, fmt.Errorf("key does not exist")
	}

	switch data.(type) {
	case []interface{}:
		return data.([]interface{}), nil
	}

	return nil, fmt.Errorf("Data fetched is not a list")
}

// FilterFn is the interface returning a boolean value of whether to include
// this value. This will change the JSON structure
type FilterFn func(interface{}) bool

func (f *Fson) filter(fn FilterFn, value interface{}) interface{} {
	switch value.(type) {
	case []interface{}:
		lst := make([]interface{}, 0, 0)
		for _, item := range value.([]interface{}) {
			if fn(item) {
				lst = append(lst, item)
			}
		}
		return lst
	case map[string]interface{}:
		mp := make(map[string]interface{})
		for k, val := range value.(map[string]interface{}) {
			if fn(val) {
				mp[k] = f.filter(fn, val)
			}
		}
		return mp
	default:
		return value
	}
}

// Filter will filter out values from the JSON where the f "filterFn" returns
// false for that value
func (f *Fson) Filter(fn FilterFn) {
	mp := make(map[string]interface{})

	for k, v := range f.Data {
		if fn(v) {
			mp[k] = f.filter(fn, v)
		}
	}

	f.Data = mp
}

// FmapFn will transform a value within the JSON into a new value, leaving the
// JSON structure alone
type FmapFn func(interface{}) interface{}

func (f *Fson) fmap(fn FmapFn, value interface{}) interface{} {
	switch value.(type) {
	case []interface{}:
		lst := make([]interface{}, 0, 0)
		for _, item := range value.([]interface{}) {
			item = f.fmap(fn, item)
			lst = append(lst, item)
		}
		return lst
	case map[string]interface{}:
		mp := make(map[string]interface{})
		for k, val := range value.(map[string]interface{}) {
			mp[k] = f.fmap(fn, val)
		}
		return mp
	default:
		return fn(value)
	}
}

// Fmap applys a function to every value in the JSON structure, mutating them in
// place
func (f *Fson) Fmap(fn FmapFn) {
	mp := make(map[string]interface{})

	for k, v := range f.Data {
		mp[k] = f.fmap(fn, v)
	}

	f.Data = mp
}

func (f *Fson) del(path []string, cur interface{}) interface{} {
	switch cur.(type) {
	case map[string]interface{}:
		mp := make(map[string]interface{})
		if len(path) == 1 {
			for key, val := range cur.(map[string]interface{}) {
				if key != path[0] {
					mp[key] = val
				}
			}
		} else {
			for key, val := range cur.(map[string]interface{}) {
				mp[key] = f.del(path[1:], val)
			}
		}
		return mp
	default:
		return cur
	}
}

// Del will delete a key from the JSON object
func (f *Fson) Del(path []string) {
	mp := make(map[string]interface{})

	for k, v := range f.Data {
		if len(path) == 1 {
			if k != path[0] {
				mp[k] = v
			}
		} else {
			if k == path[0] {
				mp[k] = f.del(path[1:], v)
			} else {
				mp[k] = v
			}
		}
	}

	f.Data = mp
}

func (f *Fson) merge(path []string, obj interface{}) {
	switch obj.(type) {
	case map[string]interface{}:
		for k, v := range obj.(map[string]interface{}) {
			newpath := append(path, k)
			if _, ok := f.Get(newpath...); ok {
				f.merge(newpath, v)
			} else {
				f.Set(v, newpath...)
			}
		}
	default:
		f.Set(obj, path...)
	}
}

// Merge will take two an Fson object and merge it with an existing Fson
// object. The new object will overwrite any matching key values
func (f *Fson) Merge(obj *Fson) {
	for k, v := range obj.Data {
		if _, ok := f.Data[k]; ok {
			f.merge([]string{k}, v)
		} else {
			f.Data[k] = v
		}
	}
}

// Pretty returns a prety printed string of the underlying JSON
func (f *Fson) Pretty() string {
	var out bytes.Buffer
	json.Indent(&out, f.Bytes(), "", "    ")
	return out.String()
}
