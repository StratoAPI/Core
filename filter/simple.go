package filter

import (
	"errors"
	"github.com/StratoAPI/Core/registry"
	"github.com/StratoAPI/Interface/filter"
	"reflect"
)

func registerSimpleFilter() {
	registry.GetRegistryInternal().RegisterFilter("simple", &SimpleFilter{})
}

type SimpleFilter struct {
}

func (sf SimpleFilter) Initialize() error {
	return nil
}

func (sf SimpleFilter) Start() error {
	return nil
}

func (sf SimpleFilter) Stop() error {
	return nil
}

func (sf SimpleFilter) ValidateFilter(processed filter.ProcessedFilter) (bool, error) {
	simple, ok := processed.Data.(*filter.Simple)

	if !ok {
		return false, errors.New("data does not match filter type")
	}

	if !simple.Operation.Valid() {
		return false, errors.New("filter operation not valid")
	}

	if simple.Operation != filter.OpEQ && simple.Operation != filter.OpNEQ {
		k := reflect.ValueOf(simple.Value).Type().Kind()
		if k == reflect.Invalid ||
			k == reflect.Bool ||
			k == reflect.Array ||
			k == reflect.Chan ||
			k == reflect.Func ||
			k == reflect.Interface ||
			k == reflect.Map ||
			k == reflect.Ptr ||
			k == reflect.Slice ||
			k == reflect.String ||
			k == reflect.Struct ||
			k == reflect.UnsafePointer {
			return false, errors.New("filter operation can only be applied to numbers")
		}
	}

	return true, nil
}

func (sf SimpleFilter) CreateFilter(_ string) (interface{}, error) {
	return &filter.Simple{}, nil
}
