package resource

import (
	"github.com/ResourceAPI/Core/plugins"
	"github.com/ResourceAPI/Core/schema"
)

func GetStore(resource string) *plugins.Storage {
	return plugins.GetStore(schema.GetSchema(resource).Meta.Store)
}

func GetResources() []string {
	// TODO
	return []string{}
}