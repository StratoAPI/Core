package nodes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Vilsol/GoLib"
	"github.com/Vilsol/ResourceAPI/database"
	"github.com/Vilsol/ResourceAPI/database/filters"
	"github.com/Vilsol/ResourceAPI/schema"
	"github.com/gorilla/mux"
)

func RegisterResourceRoutes(router GoLib.RegisterRoute) {
	router("GET", "/resource/{resource}", getResource)
	router("POST", "/resource/{resource}", storeResource)
	router("DELETE", "/resource/{resource}", deleteResource)
}

func getResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resource := mux.Vars(r)["resource"]

	if !schema.ResourceExists(resource) {
		return nil, &ErrorResourceDoesNotExist
	}

	// TODO Filters

	return database.Get().GetResources(resource, []filters.Filter{}), nil
}

func storeResource(r *http.Request) (interface{}, *GoLib.ErrorResponse) {
	resource := mux.Vars(r)["resource"]

	if !schema.ResourceExists(resource) {
		return nil, &ErrorResourceDoesNotExist
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, &ErrorCouldNotReadBody
	}

	valid, err := schema.ResourceValid(resource, string(body))

	if !valid {
		resp := ErrorResourceInvalid
		resp.Message += err.Error()
		return nil, &resp
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	database.Get().StoreResources(resource, []map[string]interface{}{data})

	return nil, nil
}

func deleteResource(_ *http.Request) (interface{}, *GoLib.ErrorResponse) {
	return nil, nil
}
