package resource

import (
	"github.com/ResourceAPI/Core/registry"
	"github.com/ResourceAPI/Interface/plugins"
	"github.com/ResourceAPI/Interface/resource"
	"github.com/ResourceAPI/Interface/schema"
)

type CoreProcessor struct {
}

var coreProcessor *CoreProcessor

func InitializeResources() {
	coreProcessor = &CoreProcessor{}
	resource.SetProcessor(coreProcessor)
}

func (cp CoreProcessor) GetStore(resource string) *plugins.Storage {
	return registry.GetRegistryInternal().GetStore(schema.GetProcessor().GetSchema(resource).Meta.Store)
}

func (cp CoreProcessor) GetResources() []string {
	// TODO
	return []string{}
}
