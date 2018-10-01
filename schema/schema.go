package schema

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var schemas = make(map[string]gojsonschema.JSONLoader)

func InitializeSchemas() {
	files, err := ioutil.ReadDir("./resources")

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			bytes, _ := ioutil.ReadFile("./resources/" + f.Name())
			schema := gojsonschema.NewBytesLoader(bytes)
			schemas[f.Name()[0:strings.LastIndex(f.Name(), ".")]] = schema
		}
	}

	fmt.Printf("Loaded %d schema(s)\n", len(schemas))
}

func ResourceExists(resource string) bool {
	_, ok := schemas[resource]
	return ok
}

func ResourceValid(resource string, data string) (bool, error) {
	result, err := gojsonschema.Validate(schemas[resource], gojsonschema.NewStringLoader(data))

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
