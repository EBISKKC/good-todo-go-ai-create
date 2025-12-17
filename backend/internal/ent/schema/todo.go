package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Immutable(),
		field.String("tenant_id").NotEmpty().Immutable(),
		field.String("user_id").NotEmpty().Immutable(),
		field.String("title").NotEmpty(),
		field.Text("description").Optional().Default(""),
		field.Bool("completed").Default(false),
		field.Bool("is_public").Default(false),
		field.Time("due_date").Optional().Nillable(),
		field.Time("completed_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenant.Type).
			Ref("todos").
			Field("tenant_id").
			Required().
			Unique().
			Immutable(),
		edge.From("user", User.Type).
			Ref("todos").
			Field("user_id").
			Required().
			Unique().
			Immutable(),
	}
}
