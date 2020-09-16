package object

import (
	"bytes"
)

// NewEnvironment function creates new environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{
		store: s,
		outer: nil,
	}
}

// NewEnclosedEnvironment function creates new environment with outer environment
/*
	Enclose identifier of parameters to current environment: outer
	Put parameters to new environment
*/
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Environment struct
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get method of Environment struct
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set method of Environment struct
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// GetOuter method of Environment struct : for analysis
func (e *Environment) GetOuter() *Environment {
	return e.outer
}

// GetInner method of Environment struct : for analysis
func (e *Environment) GetInner() map[string]Object {
	return e.store
}

// InspectEnvironment function : for analysis
func InspectEnvironment(e *Environment) string {
	var out bytes.Buffer
	if e == nil {
		//fmt.Println("nil desuyo")
		return out.String()
	}
	if len(e.GetInner()) > 0 {
		for name, value := range e.GetInner() {
			out.WriteString(name + " " + value.Inspect() + "\n")
		}
	}
	out.WriteString("outer env\n")
	out.WriteString(InspectEnvironment(e.GetOuter()))
	return out.String()
}
