package main

import (
	"fmt"
	"github.com/kubernetes-incubator/reference-docs/gen-apidocs/generators/api"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	//"reflect"
)

func main() {
	fmt.Println("Study")

	var def api.Definitions
	def = api.Definitions {
		ByGroupVersionKind: map[string]*api.Definition{},
		ByKind:             map[string]api.SortDefinitionsByVersion{},
	}
	fmt.Println(def)

	var doc *loads.Document
	var err error
	doc, err = loads.JSONSpec("test-spec.json")
	if err == nil {
		fmt.Println("Document loaded")

		var swag *spec.Swagger
		swag = doc.Spec()

		if swag != nil {
			fmt.Println("Got Swagger object")

			var defs spec.Definitions  // map[string]Schema
			defs = swag.Definitions

			fmt.Println("Length of defs:", len(defs))

			var k string
			var sch spec.Schema
			for k, sch = range defs {
				fmt.Println("Key", k)

				var k2 string
				var sch2 spec.Schema
				for k2, sch2 = range sch.Properties {
					fmt.Println("   ", k2, sch2.Type)
				}

				var exts spec.Extensions  // map[string]interface{}
				exts = sch.Extensions

				var g, v, k string
				var ok bool
				g, v, k, ok = GetGroupVersionKind(exts)
				fmt.Println("   ", g, v, k, ok)
			}
		}
	}
}

func GetGroupVersionKind(exts spec.Extensions) (string, string, string, bool) {

	if gvk, ok := exts["x-kubernetes-group-version-kind"]; ok {

		if gvkmap, ok := gvk.(map[string]interface{}); ok {
			g := gvkmap["group"].(string)
			v := gvkmap["version"].(string)
			k := gvkmap["kind"].(string)
			return g, v, k, true
		}
	}

	return "", "", "", false
}
