package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PythonDistro holds the schema definition for the PythonDistro entity.
type PythonDistro struct {
	ent.Schema
}

// Fields of the PythonDistro.
func (PythonDistro) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
	}
}

// Edges of the PythonDistro.
func (PythonDistro) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("packages", PythonPackage.Type),
	}
}
