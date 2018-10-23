package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StratoAPI/Core/registry"
	"io/ioutil"
	"os"
	"strings"

	schemaInt "github.com/StratoAPI/Interface/schema"

	"github.com/xeipuuv/gojsonschema"
)

type CoreProcessor struct {
	schemas map[string]CoreSchema
}

var coreProcessor *CoreProcessor

type CoreSchema struct {
	Parent schemaInt.Schema
	Data   *gojsonschema.Schema
}

func InitializeSchemas() {
	files, err := ioutil.ReadDir("./resources")

	if err != nil {
		panic(err)
	}

	coreProcessor = &CoreProcessor{
		schemas: make(map[string]CoreSchema),
	}

	schemaInt.SetProcessor(coreProcessor)

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile("./resources/" + f.Name())
			schema := gojsonschema.NewBytesLoader(bytes)
			var source map[string]interface{}
			err := json.Unmarshal(bytes, &source)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			if _, ok := source["meta"]; !ok {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, "resource does not contain meta object")
				continue
			}

			temp, _ := json.Marshal(source["meta"])
			meta := new(schemaInt.ResourceMeta)
			err = json.Unmarshal(temp, &meta)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			newSchema, err := gojsonschema.NewSchema(schema)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load compile resource schema %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			s := CoreSchema{
				Data: newSchema,
				Parent: schemaInt.Schema{
					Source: source,
					Meta:   *meta,
				},
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			coreProcessor.schemas[meta.Resource] = s
		}
	}

	fmt.Printf("Loaded %d schema(s)\n", len(coreProcessor.schemas))
}

func ValidateSchemas() {
	for _, schema := range coreProcessor.schemas {
		if registry.GetRegistryInternal().GetStore(schema.Parent.Meta.Store) == nil {
			panic("resource " + schema.Parent.Meta.Resource + " uses an unsupported store: " + schema.Parent.Meta.Store)
		}
	}
}

func (cp CoreProcessor) ResourceExists(resource string) bool {
	_, ok := cp.schemas[resource]
	return ok
}

func (cp CoreProcessor) ResourceValid(resource string, data string, required bool) (bool, error) {
	return cp.validateSchema(resource, gojsonschema.NewStringLoader(data), required)
}

func (cp CoreProcessor) ResourceValidGo(resource string, data interface{}, required bool) (bool, error) {
	return cp.validateSchema(resource, gojsonschema.NewGoLoader(data), required)
}

func (cp CoreProcessor) validateSchema(resource string, loader gojsonschema.JSONLoader, required bool) (bool, error) {
	result, err := cp.schemas[resource].Data.Validate(loader)

	if err != nil {
		return false, err
	}

	if result.Valid() {
		return true, nil
	}

	errs := ""
	validErrors := false

	for i, err := range result.Errors() {
		if !required && err.Type() == "required" {
			continue
		}

		validErrors = true

		err.Type()
		if i > 0 {
			errs += ", "
		}

		errs += err.String()
	}

	if !validErrors {
		return true, nil
	}

	return false, errors.New(errs)
}

func (cp CoreProcessor) GetSchema(resource string) *schemaInt.Schema {
	if _, ok := cp.schemas[resource]; !ok {
		return nil
	}

	schema := cp.schemas[resource]
	return &(schema.Parent)
}
