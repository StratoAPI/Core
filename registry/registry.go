package registry

import (
	"fmt"
	"io/ioutil"
	"plugin"

	"github.com/StratoAPI/Core/config"
	"github.com/StratoAPI/Interface/plugins"
)

type CoreRegistry struct {
	facades     map[string]*plugins.Facade
	stores      map[string]*plugins.Storage
	filters     map[string]*plugins.Filter
	middlewares map[string]*plugins.Middleware
	associates  map[string][]string
}

var coreRegistry *CoreRegistry

func GetRegistryInternal() *CoreRegistry {
	return coreRegistry
}

func InitializeRegistry() {
	coreRegistry = &CoreRegistry{
		facades:     make(map[string]*plugins.Facade),
		stores:      make(map[string]*plugins.Storage),
		filters:     make(map[string]*plugins.Filter),
		middlewares: make(map[string]*plugins.Middleware),
		associates:  make(map[string][]string),
	}
}

func InitializePlugins() {
	files, err := ioutil.ReadDir(config.Get().PluginDirectory)

	if err != nil {
		panic(err)
	}

	plugins.SetRegistry(coreRegistry)

	loadedPlugins := make([]plugins.Plugin, 0)
	for _, f := range files {
		plug, err := plugin.Open(config.Get().PluginDirectory + "/" + f.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		entrypoint, err := plug.Lookup("CorePlugin")
		if err != nil {
			fmt.Println(err)
			continue
		}

		var pl plugins.Plugin
		pl, ok := entrypoint.(plugins.Plugin)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			continue
		}

		pl.Entrypoint()
		loadedPlugins = append(loadedPlugins, pl)
	}

	pluginNames := make([]string, 0)
	facadeNames := make([]string, 0)
	storageNames := make([]string, 0)
	filterNames := make([]string, 0)

	for _, v := range loadedPlugins {
		pluginNames = append(pluginNames, v.Name())
	}

	for k := range coreRegistry.facades {
		facadeNames = append(facadeNames, k)
	}

	for k := range coreRegistry.stores {
		storageNames = append(storageNames, k)
	}

	for k := range coreRegistry.filters {
		filterNames = append(filterNames, k)
	}

	fmt.Printf("Loaded %d plugin(s): %+v\n", len(loadedPlugins), pluginNames)
	fmt.Printf("Loaded %d facade(s): %+v\n", len(coreRegistry.facades), facadeNames)
	fmt.Printf("Loaded %d storage(s): %+v\n", len(coreRegistry.stores), storageNames)
	fmt.Printf("Loaded %d filter(s): %+v\n", len(coreRegistry.filters), filterNames)
	fmt.Printf("Loaded %d filter association(s)\n", len(coreRegistry.associates))
}

func (cr CoreRegistry) RegisterFacade(name string, facade plugins.Facade) error {
	if _, ok := cr.facades[name]; ok {
		panic("Facade with name " + name + " is already registered!")
	}

	cr.facades[name] = &facade

	return nil
}

func (cr CoreRegistry) RegisterStorage(name string, storage plugins.Storage) error {
	if _, ok := cr.stores[name]; ok {
		panic("Storage with name " + name + " is already registered!")
	}

	cr.stores[name] = &storage

	return nil
}

func (cr CoreRegistry) RegisterFilter(name string, filter plugins.Filter) error {
	if _, ok := cr.filters[name]; ok {
		panic("Filter with name " + name + " is already registered!")
	}

	cr.filters[name] = &filter

	return nil
}

func (cr CoreRegistry) RegisterMiddleware(name string, middleware plugins.Middleware) error {
	if _, ok := cr.middlewares[name]; ok {
		panic("Middleware with name " + name + " is already registered!")
	}

	cr.middlewares[name] = &middleware

	return nil
}

func (cr CoreRegistry) AssociateFilter(filter string, storage string) error {
	if _, ok := cr.associates[filter]; !ok {
		cr.associates[filter] = make([]string, 0)
	}

	supportedStores := cr.associates[filter]

	for _, store := range supportedStores {
		if store == storage {
			panic("Filter " + filter + " is already associated with storage " + storage + "!")
		}
	}

	supportedStores = append(supportedStores, storage)
	cr.associates[filter] = supportedStores

	return nil
}

func (cr CoreRegistry) GetStore(store string) *plugins.Storage {
	return cr.stores[store]
}

func (cr CoreRegistry) GetFilter(filter string) *plugins.Filter {
	return cr.filters[filter]
}

func (cr CoreRegistry) GetAssociates(filter string) []string {
	return cr.associates[filter]
}

func (cr CoreRegistry) GetMiddleware(middleware string) *plugins.Middleware {
	return cr.middlewares[middleware]
}
