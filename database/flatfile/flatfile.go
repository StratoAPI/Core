package flatfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Vilsol/ResourceAPI/database/filters"
)

var flatfile Flatfile

func Initialize(location string) *Flatfile {
	var data map[string][]map[string]interface{}

	if _, err := os.Stat(location); err == nil {
		bytes, _ := ioutil.ReadFile(location)
		json.Unmarshal(bytes, &data)
	}

	if data == nil {
		data = make(map[string][]map[string]interface{})
	}

	flatfile = Flatfile{
		Location: location,
		Data:     data,
	}

	return &flatfile
}

func (flatfile *Flatfile) Save() {
	if _, err := os.Stat(flatfile.Location); os.IsNotExist(err) {
		f, err := os.Create(flatfile.Location)
		if err != nil {
			fmt.Println("Failed to save: " + err.Error())
		}
		f.Close()
	}

	bytes, _ := json.Marshal(flatfile.Data)
	err := ioutil.WriteFile(flatfile.Location, bytes, 0666)
	if err != nil {
		fmt.Println("Failed to save: " + err.Error())
	}
}

func (flatfile *Flatfile) Wipe() {
	flatfile.Data = make(map[string][]map[string]interface{})
	flatfile.Save()
}

func (flatfile *Flatfile) GetResources(name string, filters []filters.Filter) []map[string]interface{} {
	resources, ok := flatfile.Data[name]

	if !ok {
		return make([]map[string]interface{}, 0)
	}

	// TODO Filters

	return resources
}

func (flatfile *Flatfile) StoreResources(name string, resources []map[string]interface{}) {
	_, ok := flatfile.Data[name]

	if !ok {
		flatfile.Data[name] = resources
	} else {
		flatfile.Data[name] = append(flatfile.Data[name], resources...)
	}

	flatfile.Save()
}

func (flatfile *Flatfile) DeleteResources(name string, filters []filters.Filter) {
	panic("implement me")
}
