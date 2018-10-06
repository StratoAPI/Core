package plugins

import (
	"fmt"
	"io/ioutil"
	"plugin"

	"github.com/ResourceAPI/Core/config"
)

var facades map[string]*Facade
var stores map[string]*Storage
var filters map[string]*Filter
var associates map[string][]string

func InitializePlugins() {
	facades = make(map[string]*Facade)
	stores = make(map[string]*Storage)
	filters = make(map[string]*Filter)
	associates = make(map[string][]string)

	files, err := ioutil.ReadDir(config.Get().PluginDirectory)

	if err != nil {
		panic(err)
	}

	loadedPlugins := make([]Plugin, 0)
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

		var pl Plugin
		pl, ok := entrypoint.(Plugin)
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

	for k := range facades {
		facadeNames = append(facadeNames, k)
	}

	for k := range stores {
		storageNames = append(storageNames, k)
	}

	for k := range filters {
		filterNames = append(filterNames, k)
	}

	fmt.Printf("Loaded %d plugin(s): %+v\n", len(loadedPlugins), pluginNames)
	fmt.Printf("Loaded %d facade(s): %+v\n", len(facades), facadeNames)
	fmt.Printf("Loaded %d storage(s): %+v\n", len(stores), storageNames)
	fmt.Printf("Loaded %d filter(s): %+v\n", len(filters), filterNames)
	fmt.Printf("Loaded %d filter association(s)\n", len(associates))
}

func RegisterFacade(name string, facade Facade) error {
	if _, ok := facades[name]; ok {
		panic("Facade with name " + name + " is already registered!")
	}

	facades[name] = &facade

	return nil
}

func RegisterStorage(name string, storage Storage) error {
	if _, ok := stores[name]; ok {
		panic("Storage with name " + name + " is already registered!")
	}

	stores[name] = &storage

	return nil
}

func RegisterFilter(name string, filter Filter) error {
	if _, ok := filters[name]; ok {
		panic("Filter with name " + name + " is already registered!")
	}

	filters[name] = &filter

	return nil
}

func AssociateFilter(filter string, storage string) error {
	if _, ok := associates[filter]; !ok {
		associates[filter] = make([]string, 0)
	}

	supportedStores := associates[filter]

	for _, store := range supportedStores {
		if store == storage {
			panic("Filter " + filter + " is already associated with storage " + storage + "!")
		}
	}

	supportedStores = append(supportedStores, storage)
	associates[filter] = supportedStores

	return nil
}

func GetStore(store string) *Storage {
	return stores[store]
}
