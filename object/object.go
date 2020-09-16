package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/BOBO1997/monkey/ast"
)

// ObjectType type represents type of Object
type ObjectType string

// object name
const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

/* ====== object definitions ====== */

// integer

// Object interface requires Type() method and Inspect() method
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer struct
type Integer struct {
	Value int64
}

// Inspect method of Integer struct
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type method of Integer struct
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// boolean

// Boolean struct
type Boolean struct {
	Value bool
}

// Inspect method of Boolean struct
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type method of Boolean struct
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// string

// String struct
type String struct {
	Value string
}

// Inspect method of String struct
func (s *String) Inspect() string {
	return s.Value
}

// Type method of String struct
func (s *String) Type() ObjectType {
	return STRING_OBJ
}

// null

// Null struct
type Null struct{}

// Inspect method of Null struct
func (n *Null) Inspect() string {
	return fmt.Sprintf("null")
}

// Type method of Null struct
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// ReturnValue struct
type ReturnValue struct {
	Value Object
}

// Inspect method of ReturnValue struct
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Type method of ReturnValue struct
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

// Error struct
type Error struct {
	Message string
}

// Inspect method of Error struct
func (err *Error) Inspect() string {
	return "ERROR: " + err.Message
}

// Type method of Error struct
func (err *Error) Type() ObjectType {
	return ERROR_OBJ
}

// Function struct
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Inspect method of Function struct
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.Value)
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Type method of Function struct
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

// builtin function

// BuiltinFunction type
type BuiltinFunction func(args ...Object) Object

// Builtin struct
type Builtin struct {
	Fn BuiltinFunction
}

// Inspect method of Builtin struct
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Type method of Builtin struct
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

// Array struct
type Array struct {
	Elements []Object
}

// Inspect method of Array struct
func (arr *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, element := range arr.Elements {
		elements = append(elements, element.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Type method of Array struct
func (arr *Array) Type() ObjectType {
	return ARRAY_OBJ
}

// HashKey struct
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// HashKey method of Boolean struct
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{
		Type:  b.Type(),
		Value: value,
	}
}

// HashKey method of Integer struct
func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

// HashKey method of String struct
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value)) // ? not []rune
	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}
