package resource

import (
	"errors"
	"github.com/StratoAPI/Core/config"
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

func (cp CoreProcessor) GetStoreName(resource string) string {
	resourceStore := schema.GetProcessor().GetSchema(resource).Meta.Store

	if resourceStore == "" {
		resourceStore = config.Get().DefaultStore
	}

	return resourceStore
}

func (cp CoreProcessor) GetStore(resource string) *plugins.Storage {
	return registry.GetRegistryInternal().GetStore(cp.GetStoreName(resource))
}

func (cp CoreProcessor) GetResourceList() []string {
	// TODO
	return []string{}
}

func (cp CoreProcessor) GetResources(resource string, filters []filter.ProcessedFilter) ([]map[string]interface{}, error) {
	resourceStore := cp.GetStoreName(resource)

	err := checkStoreFilterAssociates(resourceStore, filters)

	if err != nil {
		return nil, err
	}

	return (*registry.GetRegistryInternal().GetStore(resourceStore)).GetResources(resource, filters)
}

func (cp CoreProcessor) CreateResources(resource string, data []map[string]interface{}) error {
	return (*cp.GetStore(resource)).CreateResources(resource, data)
}

func (cp CoreProcessor) UpdateResources(resource string, data map[string]interface{}, filters []filter.ProcessedFilter) error {
	resourceStore := cp.GetStoreName(resource)

	err := checkStoreFilterAssociates(resourceStore, filters)

	if err != nil {
		return err
	}

	return (*registry.GetRegistryInternal().GetStore(resourceStore)).UpdateResources(resource, data, filters)
}

func (cp CoreProcessor) DeleteResources(resource string, filters []filter.ProcessedFilter) error {
	resourceStore := cp.GetStoreName(resource)

	err := checkStoreFilterAssociates(resourceStore, filters)

	if err != nil {
		return err
	}

	return (*registry.GetRegistryInternal().GetStore(resourceStore)).DeleteResources(resource, filters)
}

func checkStoreFilterAssociates(resourceStore string, filters []filter.ProcessedFilter) error {
	for _, f := range filters {
		associates := registry.GetRegistryInternal().GetAssociates(f.Type)

		if len(associates) == 0 {
			return errors.New("the store does not support a provided filter")
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

		return errors.New("the store does not support a provided filter")
	}

	return nil
}
