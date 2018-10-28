package middleware

import (
	"github.com/StratoAPI/Core/registry"
	"github.com/StratoAPI/Interface/middleware"
	"github.com/StratoAPI/Interface/schema"
)

type CoreProcessor struct {
}

var coreProcessor *CoreProcessor

func InitializeMiddleware() {
	coreProcessor = &CoreProcessor{}
	middleware.SetProcessor(coreProcessor)
}

func (cp CoreProcessor) Request(resource string, headers map[string][]string) *middleware.RequestResponse {
	middlewares := schema.GetProcessor().GetSchema(resource).Meta.Middlewares

	for _, mwMeta := range middlewares {
		mw := registry.GetRegistryInternal().GetMiddleware(mwMeta.Type)
		if mw == nil {
			panic("Middleware " + mwMeta.Type + " is nil")
		}

		request := (*mw).Request(resource, headers, mwMeta.Data)
		if request != nil {
			return request
		}
	}

	return nil
}

func (cp CoreProcessor) Response(resource string, headers map[string][]string, data []map[string]interface{}) ([]map[string]interface{}, *middleware.RequestResponse) {
	middlewares := schema.GetProcessor().GetSchema(resource).Meta.Middlewares

	for _, mwMeta := range middlewares {
		mw := registry.GetRegistryInternal().GetMiddleware(mwMeta.Type)
		if mw == nil {
			panic("Middleware " + mwMeta.Type + " is nil")
		}

		result, response := (*mw).Response(resource, headers, data, mwMeta.Data)
		if response != nil {
			return nil, response
		}

		data = result
	}

	return data, nil
}
