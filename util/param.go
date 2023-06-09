package util

import (
	// "log"
	"reflect"
)

// Params constructs function parameters similar to Python **kwargs.
// It is used to add multiple optional parameters to a function.
type Params struct {
	params map[string]interface{}
}

// NewParams creates a new Params.
func NewParams(values map[string]interface{}) *Params {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &Params{params: values}
}

// Pop pops(removes) and returns parameter with key from Params.
// It returns nil if value not found.
func (p *Params) Pop(key string) (value interface{}) {
	if val, ok := p.params[key]; ok {
		value = val
		delete(p.params, key)
	}

	return value
}

// Param returns value of parameter corresponding to key.
// It returns nil if not found.
func (p *Params) Param(key string) (value interface{}) {
	if val, ok := p.params[key]; ok {
		return val
	}

	return nil
}

// Has returns whether having a not-nil value corresponding to given key.
func (p *Params) Has(key string) bool {
	val, ok := p.params[key]
	if !ok {
		return false
	}

	if IsNil(val) {
		return false
	}

	return true
}

// Get returns a parameter value by its key.
// If returns nil/default value if parameter type or value is nil.
func (p *Params) Get(key string, defaultValueOpt ...interface{}) (val interface{}) {
	var defaultVal interface{} = nil
	if len(defaultValueOpt) > 0 {
		defaultVal = defaultValueOpt[0]
	}
	if val, ok := p.params[key]; !ok {
		return defaultVal
	} else {
		if IsNil(val) {
			return defaultVal
		} else {
			return val
		}
	}
}

// Copy (shallow) copies parameter from one Params to other Params
func (p *Params) Copy(params *Params, key string, newKeyOpt ...string) {
	newKey := key
	if len(newKeyOpt) > 0 {
		newKey = newKeyOpt[0]
	}
	val := params.Get(key)
	p.Set(newKey, val)
}

// DeepCopy copies a param with given name.
func (p *Params) DeepCopy(params *Params, key string, newKeyOpt ...string) {
	newKey := key
	if len(newKeyOpt) > 0 {
		newKey = newKeyOpt[0]
	}

	newVal := params.deepCopy(key)
	p.Set(newKey, newVal)
}

// Clone clones (deep copy) all parameters to new Params.
func (p *Params) Clone() *Params {
	out := NewParams(nil)
	for k := range p.Values() {
		newVal := p.deepCopy(k)
		out.Set(k, newVal)
	}

	return out
}

func (p *Params) deepCopy(key string) interface{} {
	v := p.Get(key)
	if v == nil {
		return nil
	}

	typ := reflect.TypeOf(v).String()
	switch typ {
	case "*util.Params":
		vals := v.(*Params)
		out := NewParams(nil)
		for key := range vals.Values() {
			newVal := vals.deepCopy(key)
			out.Set(key, newVal)
		}
		return out
	default:
		return v
	}
}

// Select selects a subset of parameters
func (p *Params) Select(keys []string) *Params {
	params := NewParams(nil)
	for _, k := range keys {
		v := p.Get(k)
		params.Set(k, v)
	}

	return params
}

// Keys returns a slice of parameter names.
func (p *Params) Keys() []string {
	keys := []string{}
	for k := range p.params {
		keys = append(keys, k)
	}

	return keys
}

func (p *Params) Delete(key string) {
	v := p.Get(key)
	if v == nil {
		delete(p.Values(), key)
		return
	}

	delete(p.Values(), key)
}

func (p *Params) DeleteAll() {
	for k := range p.Values() {
		p.Delete(k)
	}
}

// Values returns a map of parameters.
func (p *Params) Values() map[string]interface{} {
	return p.params
}

// Len returns number of parameters.
func (p *Params) Len() int {
	return len(p.params)
}

// Set adds/updates a new parameter to Params.
func (p *Params) Set(key string, value interface{}) {
	p.params[key] = value
}

type ParamOption func(*Params)

// WithParam adds a parameter option to Params.
func WithParam(key string, value interface{}) ParamOption {
	return func(p *Params) {
		p.params[key] = value
	}
}

func WithParams(p *Params) []ParamOption {
	var opts []ParamOption
	params := p.Values()
	for k, v := range params {
		opt := WithParam(k, v)
		opts = append(opts, opt)
	}

	return opts
}

// DefaultParams creates Params with default values.
func DefaultParams() *Params {
	panic("NotImplementedError. Should be implemented by end user or struct that embeds Params struct.")
}

// IsNil checks whether is a variable of type interface{}
// is nil.
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.ValueOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
