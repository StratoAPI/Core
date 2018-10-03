package Core

import (
	"github.com/ResourceAPI/Core/config"
	"github.com/ResourceAPI/Core/database"
	"github.com/ResourceAPI/Core/plugins"
	"github.com/ResourceAPI/Core/schema"
)

func Run() {
	// Initialize Core
	config.InitializeConfig()
	schema.InitializeSchemas()
	database.InitializeDatabase()

	// Initialize Plugins
	plugins.InitializePlugins()
	plugins.InitializeStores()
	plugins.InitializeFilters()
	plugins.InitializeFacades()

	// Start up stores
	plugins.StartStores()

	// Start up filters
	plugins.StartFilters()

	// Start up facades
	plugins.StartFacades()

	// Wait for goroutines
	plugins.WaitForGoroutines()

	// TODO Graceful shutdown
}
