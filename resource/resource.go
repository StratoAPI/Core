package resource

import (
	"errors"
	"github.com/StratoAPI/Core/registry"
	"github.com/StratoAPI/Interface/filter"
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

func (cp CoreProcessor) GetResourceList() []string {
	// TODO
	return []string{}
}

func (cp CoreProcessor) GetResources(resource string, filters []filter.ProcessedFilter) ([]map[string]interface{}, error) {
	resourceStore := schema.GetProcessor().GetSchema(resource).Meta.Store

	for _, f := range filters {
		associates := registry.GetRegistryInternal().GetAssociates(f.Type)

		if len(associates) == 0 {
			return nil, errors.New("the store does not support a provided filter")
		}

		found := false
		for _, store := range associates {
			if store == resourceStore {
				found = true
				break
			}
		}

		if found {
			break
		}

		return nil, errors.New("the store does not support a provided filter")
	}

	return (*registry.GetRegistryInternal().GetStore(schema.GetProcessor().GetSchema(resource).Meta.Store)).GetResources(resource, filters)
}
