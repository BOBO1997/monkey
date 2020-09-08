package object

import (
	"fmt"
)

// ObjectType type represents type of Object
type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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

// null

// Null struct
type Null struct{}

// Inspect method of Boolean struct
func (n *Null) Inspect() string {
	return fmt.Sprintf("null")
}

// Type method of Boolean struct
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}
