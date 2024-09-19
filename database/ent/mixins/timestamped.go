package mixins

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Timestamped struct {
	mixin.Schema
}

func (Timestamped) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Comment("The date the resource was created at"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("The date the resource was deleted at"),
	}
}

func (Timestamped) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}
