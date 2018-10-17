package resource

import (
	"github.com/StratoAPI/Core/registry"
	"github.com/StratoAPI/Interface/plugins"
	"github.com/StratoAPI/Interface/resource"
	"github.com/StratoAPI/Interface/schema"
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
