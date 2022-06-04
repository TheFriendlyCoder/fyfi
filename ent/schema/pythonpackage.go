package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// PythonPackage holds the schema definition for the PythonPackage entity.
type PythonPackage struct {
	ent.Schema
}

// Fields of the PythonPackage.
func (PythonPackage) Fields() []ent.Field {
	// URL           string
	// filename      string
	// pythonVersion string
	// checksum      string
	return []ent.Field{
		field.String("url").NotEmpty(),
		field.String("filename"),
		field.String("pythonVersion"),
		field.String("checksum"),
	}
}

// Edges of the PythonPackage.
func (PythonPackage) Edges() []ent.Edge {
	return nil
}
