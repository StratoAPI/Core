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
	Meta   ResourceMeta
}

type ResourceMeta struct {
	Resource string `json:"resource"`
	Store    string `json:"store"`
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
			meta := new(ResourceMeta)
			err = json.Unmarshal(temp, &meta)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			s := Schema{
				Data:   schema,
				Source: source,
				Meta:   *meta,
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load resource %s:\n", f.Name())
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			schemas[meta.Resource] = s
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
