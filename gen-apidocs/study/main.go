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

				var group, version, kind string
				group, version, kind, _ = GetGroupVersionKind(exts)

				def  := api.Definition{
					Name: kind,
					Group: api.ApiGroup(group),
					Version: api.ApiVersion(version),
					Kind: api.ApiKind(kind),
					ShowGroup: true,
					Resource: "",
				}

				fmt.Println("    def:", def.Name, def.Group, def.Version, def.Kind, def.ShowGroup)
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
