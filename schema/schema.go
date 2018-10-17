package schema

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Data   gojsonschema.JSONLoader
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

			s := CoreSchema{
				Data: schema,
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

func (cp CoreProcessor) ResourceExists(resource string) bool {
	_, ok := cp.schemas[resource]
	return ok
}

func (cp CoreProcessor) ResourceValid(resource string, data string) (bool, error) {
	result, err := gojsonschema.Validate(cp.schemas[resource].Data, gojsonschema.NewStringLoader(data))

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

func (cp CoreProcessor) GetSchema(resource string) *schemaInt.Schema {
	if _, ok := cp.schemas[resource]; !ok {
		return nil
	}

	schema := cp.schemas[resource]
	return &(schema.Parent)
}
