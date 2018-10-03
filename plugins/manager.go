package plugins

import (
	"sync"
)

var storageWaitGroup sync.WaitGroup
var facadeWaitGroup sync.WaitGroup
var filtersWaitGroup sync.WaitGroup

func InitializeStores() {
	for _, store := range stores {
		(*store).Initialize()
	}
}

func StartStores() {
	storageWaitGroup.Add(len(stores))

	for _, store := range stores {
		go func() {
			defer storageWaitGroup.Done()
			(*store).Start()
		}()
	}
}

func StopStores() {
	// TODO 30s timeout
	for _, store := range stores {
		(*store).Stop()
	}
}

func InitializeFacades() {
	for _, facade := range facades {
		(*facade).Initialize()
	}
}

func StartFacades() {
	facadeWaitGroup.Add(len(facades))

	for _, facade := range facades {
		go func() {
			defer facadeWaitGroup.Done()
			(*facade).Start()
		}()
	}
}

func StopFacades() {
	// TODO 30s timeout
	for _, facade := range facades {
		(*facade).Stop()
	}
}

func InitializeFilters() {
	for _, filter := range filters {
		(*filter).Initialize()
	}
}

func StartFilters() {
	filtersWaitGroup.Add(len(facades))

	for _, filter := range filters {
		go func() {
			defer filtersWaitGroup.Done()
			(*filter).Start()
		}()
	}
}

func StopFilters() {
	// TODO 30s timeout
	for _, filter := range filters {
		(*filter).Stop()
	}
}

func WaitForGoroutines() {
	facadeWaitGroup.Wait()
	storageWaitGroup.Wait()
	filtersWaitGroup.Wait()
}
