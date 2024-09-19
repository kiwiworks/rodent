package mixins

import "entgo.io/ent"

var All = []ent.Mixin{SoftDeletable{}, Timestamped{}, Resource{}}
