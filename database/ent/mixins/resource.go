package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Resource struct {
	mixin.Schema
}

func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("public_id", uuid.UUID{}).
			Default(uuid.New).
			Comment("A resource public ID which will be exposed"),
	}
}

func (Resource) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),
	}
}
