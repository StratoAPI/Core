package registry

import (
	"sync"
)

var storageWaitGroup sync.WaitGroup
var facadeWaitGroup sync.WaitGroup
var filtersWaitGroup sync.WaitGroup

func InitializeStores() {
	for _, store := range coreRegistry.stores {
		err := (*store).Initialize()

		if err != nil {
			panic(err)
		}
	}
}

func StartStores() {
	storageWaitGroup.Add(len(coreRegistry.stores))

	for _, store := range coreRegistry.stores {
		go func() {
			defer storageWaitGroup.Done()
			err := (*store).Start()

			if err != nil {
				panic(err)
			}
		}()
	}
}

func StopStores() {
	// TODO 30s timeout
	for _, store := range coreRegistry.stores {
		err := (*store).Stop()

		if err != nil {
			panic(err)
		}
	}
}

func InitializeFacades() {
	for _, facade := range coreRegistry.facades {
		err := (*facade).Initialize()

		if err != nil {
			panic(err)
		}
	}
}

func StartFacades() {
	facadeWaitGroup.Add(len(coreRegistry.facades))

	for _, facade := range coreRegistry.facades {
		go func() {
			defer facadeWaitGroup.Done()
			err := (*facade).Start()

			if err != nil {
				panic(err)
			}
		}()
	}
}

func StopFacades() {
	// TODO 30s timeout
	for _, facade := range coreRegistry.facades {
		err := (*facade).Stop()

		if err != nil {
			panic(err)
		}
	}
}

func InitializeFilters() {
	for _, filter := range coreRegistry.filters {
		err := (*filter).Initialize()

		if err != nil {
			panic(err)
		}
	}
}

func StartFilters() {
	filtersWaitGroup.Add(len(coreRegistry.facades))

	for _, filter := range coreRegistry.filters {
		go func() {
			defer filtersWaitGroup.Done()
			err := (*filter).Start()

			if err != nil {
				panic(err)
			}
		}()
	}
}

func StopFilters() {
	// TODO 30s timeout
	for _, filter := range coreRegistry.filters {
		err := (*filter).Stop()

		if err != nil {
			panic(err)
		}
	}
}

func WaitForGoroutines() {
	facadeWaitGroup.Wait()
	storageWaitGroup.Wait()
	filtersWaitGroup.Wait()
}
