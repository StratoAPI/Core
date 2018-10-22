package filter

import (
	"errors"
	"github.com/StratoAPI/Core/registry"
	"github.com/StratoAPI/Interface/filter"
)

type CoreProcessor struct {
}

var coreProcessor *CoreProcessor

func InitializeFilters() {
	coreProcessor = &CoreProcessor{}
	filter.SetProcessor(coreProcessor)
	registerSimpleFilter()
}

func (cp CoreProcessor) FilterExists(filter string) bool {
	return registry.GetRegistryInternal().GetFilter(filter) != nil
}

func (cp CoreProcessor) CreateFilter(filter string) interface{} {
	f := registry.GetRegistryInternal().GetFilter(filter)
	if f == nil {
		return nil
	}

	created, err := (*f).CreateFilter(filter)

	if err != nil {
		panic(err)
	}

	return created
}

func (cp CoreProcessor) ValidateFilter(filter filter.ProcessedFilter) (bool, error) {
	f := registry.GetRegistryInternal().GetFilter(filter.Type)
	if f == nil {
		return false, errors.New("filter not found")
	}
	return (*f).ValidateFilter(filter)
}
