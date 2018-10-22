package Core

import (
	"github.com/StratoAPI/Core/config"
	"github.com/StratoAPI/Core/filter"
	"github.com/StratoAPI/Core/registry"
	"github.com/StratoAPI/Core/resource"
	"github.com/StratoAPI/Core/schema"
)

func Run() {
	// Initialize Core
	config.InitializeConfig()
	registry.InitializeRegistry()
	schema.InitializeSchemas()
	resource.InitializeResources()
	filter.InitializeFilters()

	// Initialize Plugins
	registry.InitializePlugins()
	registry.InitializeStores()
	registry.InitializeFilters()
	registry.InitializeFacades()

	// Start up stores
	registry.StartStores()

	// Start up filters
	registry.StartFilters()

	// Start up facades
	registry.StartFacades()

	// Wait for goroutines
	registry.WaitForGoroutines()

	// TODO Graceful shutdown
}
