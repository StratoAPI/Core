package registry

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
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
	for _, filter := range coreRegistry.filters {
		err := (*filter).Stop()

		if err != nil {
			panic(err)
		}
	}
}

func WaitForGoroutines() {
	done := make(chan bool, 1)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		facadeWaitGroup.Wait()
		storageWaitGroup.Wait()
		filtersWaitGroup.Wait()
		done <- true
	}()

	<-stop

	fmt.Println("Shutting down the server...")

	go func() {
		StopFacades()
		StopStores()
		StopFilters()
	}()

	select {
	// TODO Make time configurable
	case <-time.After(30 * time.Second):
		panic("Graceful shutdown failed!")
	case <-done:
		fmt.Println("Server shut down")
	}
}
