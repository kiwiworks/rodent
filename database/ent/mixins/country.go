package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/kiwiworks/rodent/slices"

	"github.com/biter777/countries"
)

func toIso3(country countries.CountryCode) string {
	return country.Alpha3()
}

var allCountries = slices.Map(countries.All(), toIso3)

type Country struct {
	mixin.Schema
}

func (Country) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("country_iso3").
			Values(allCountries...).
			Comment("The ISO 3166-1 alpha 3 code of the country"),
	}
}
