package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type Schema struct {
	Data   gojsonschema.JSONLoader
	Source map[string]interface{}
}

var schemas = make(map[string]Schema)

func InitializeSchemas() {
	files, err := ioutil.ReadDir("./resources")

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile("./resources/" + f.Name())
			schema := gojsonschema.NewBytesLoader(bytes)
			var source interface{}
			json.Unmarshal(bytes, &source)

			s := Schema{
				Data:   schema,
				Source: source.(map[string]interface{}),
			}

			resource, err := s.GetResource()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			schemas[resource] = s
		}
	}

	fmt.Printf("Loaded %d schema(s)\n", len(schemas))
}

func ResourceExists(resource string) bool {
	_, ok := schemas[resource]
	return ok
}

func ResourceValid(resource string, data string) (bool, error) {
	result, err := gojsonschema.Validate(schemas[resource].Data, gojsonschema.NewStringLoader(data))

	if err != nil {
		return false, err
	}

	if result.Valid() {
		return true, nil
	}

	errs := ""

	for i, err := range result.Errors() {
		if i > 0 {
			errs += ", "
		}

		errs += err.String()
	}

	return false, errors.New(errs)
}

func GetSchema(resource string) *Schema {
	if _, ok := schemas[resource]; !ok {
		return nil
	}

	schema := schemas[resource]
	return &schema
}

func (schema Schema) GetRaw(key string) (interface{}, error) {
	if _, ok := schema.Source[key]; !ok {
		return nil, errors.New("key '" + key + "' does not exist in schema")
	}

	return schema.Source[key], nil
}

func (schema Schema) GetRawString(key string) (string, error) {
	raw, err := schema.GetRaw(key)

	if err != nil {
		return "", err
	}

	casted, ok := raw.(string)

	if !ok {
		return "", errors.New("key '" + key + "' is not of type string in schema")
	}

	return casted, nil
}

func (schema Schema) GetStore() (string, error) {
	return schema.GetRawString("store")
}

func (schema Schema) GetResource() (string, error) {
	return schema.GetRawString("resource")
}
