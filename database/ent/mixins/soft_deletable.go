package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type SoftDeletable struct {
	mixin.Schema
}

func (SoftDeletable) Fields() []ent.Field {
	return []ent.Field{
		field.Time("deleted_at").
			Optional().
			Comment("Represent the time at which the resource was soft-deleted"),
	}
}

func (SoftDeletable) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}

func (SoftDeletable) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{}
}
