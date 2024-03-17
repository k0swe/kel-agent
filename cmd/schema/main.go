package main

import (
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/k0swe/kel-agent/internal/config"
)

func main() {
	sch := jsonschema.Reflect(config.Config{})
	schema, _ := sch.MarshalJSON()
	fmt.Println(string(schema))
}
