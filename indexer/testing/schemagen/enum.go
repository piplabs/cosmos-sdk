package schemagen

import (
	indexerbase "cosmossdk.io/schema"
	"pgregory.net/rapid"
)

var enumValuesGen = rapid.SliceOfNDistinct(Name, 1, 10, func(x string) string { return x })

var EnumDefinition = rapid.Custom(func(t *rapid.T) schema.EnumDefinition {
	enum := schema.EnumDefinition{
		Name:   Name.Draw(t, "name"),
		Values: enumValuesGen.Draw(t, "values"),
	}

	return enum
})