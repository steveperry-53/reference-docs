/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"fmt"
	"strings"

	"errors"
	"github.com/go-openapi/spec"
)

func GetGroupVersionKind() {

}

// GetDefinitionVersionKind returns the api group, version, and kind for the spec.  This is the primary key of a Definition.
func GetDefinitionVersionKind(s spec.Schema) (string, string, string) {
	// Get the reference for complex types
	if IsDefinition(s) {
		s := fmt.Sprintf("%s", s.SchemaProps.Ref.GetPointer())
		s = strings.Replace(s, "/definitions/", "", -1)
		name := strings.Split(s, ".")

		var group, version, kind string
		if name[len(name)-3] == "api" {
			// e.g. "io.k8s.apimachinery.pkg.api.resource.Quantity"
			group = "core"
			version = name[len(name)-2]
			kind = name[len(name)-1]
		} else if name[len(name)-4] == "api" {
			// e.g. "io.k8s.api.core.v1.Pod"
			group = name[len(name)-3]
			version = name[len(name)-2]
			kind = name[len(name)-1]
		} else if name[len(name)-4] == "apis" {
			// e.g. "io.k8s.apimachinery.pkg.apis.meta.v1.Status"
			group = name[len(name)-3]
			version = name[len(name)-2]
			kind = name[len(name)-1]
		} else if name[len(name)-3] == "util" || name[len(name)-3] == "pkg" {
			// e.g. io.k8s.apimachinery.pkg.util.intstr.IntOrString
			// e.g. io.k8s.apimachinery.pkg.runtime.RawExtension
			return "", "", ""
		} else {
			panic(errors.New(fmt.Sprintf("Could not locate group for %s", name)))
		}
		return group, version, kind
	}
	// Recurse if type is array
	if IsArray(s) {
		return GetDefinitionVersionKind(*s.Items.Schema)
	}
	return "", "", ""
}

// GetTypeName returns the display name of a Schema.  This is the api kind for definitions and the type for
// primitive types.  Arrays of objects have "array" appended.
func GetTypeName(s spec.Schema) string {
	// Get the reference for complex types
	if IsDefinition(s) {
		_, _, name := GetDefinitionVersionKind(s)
		return name
	}
	// Recurse if type is array
	if IsArray(s) {
		return fmt.Sprintf("%s array", GetTypeName(*s.Items.Schema))
	}
	// Get the value for primitive types
	if len(s.Type) > 0 {
		return fmt.Sprintf("%s", s.Type[0])
	}
	panic(fmt.Errorf("No type found for object %v", s))
}

// IsArray returns true if the type is an array type.
func IsArray(s spec.Schema) bool {
	//if s == nil {
	//	return false
	//}
	return len(s.Type) > 0 && s.Type[0] == "array"
}

// IsDefinition returns true if Schema is a complex type that should have a Definition.
func IsDefinition(s spec.Schema) bool {
	return len(s.SchemaProps.Ref.GetPointer().String()) > 0
}

func PrintApiGroups(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("ApiGroups")

	for _, g := range config.ApiGroups {
		fmt.Println("   ", g)
	}
}

func PrintOperationCategories(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("OperationCategories")
	
	for _, oc := range config.OperationCategories {
		fmt.Println("   ", oc.Name)
	}
}

func PrintResourceCategories(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("ResourceCategories")

	for _, rc := range config.ResourceCategories {
		fmt.Println("   ", rc.Name)
	}
}

func PrintGroupMap(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("GroupMap")

	fmt.Println("   GroupMap length: ", len(config.GroupMap))

	for k, v := range config.GroupMap {
		fmt.Println("   ", k, "  ", v)
	}
}

func PrintResourcesInResourceCategory(rc ResourceCategory) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("Resources in resource category")
	fmt.Println("   ", rc.Name)

	for _, r := range rc.Resources {
		fmt.Println("      ", r.Group, " ", r.Version, " ", r.Name)
	}
}

func PrintDefinition(config *Config, gvk string) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("Definition")
	definition := config.Definitions.ByGroupVersionKind[gvk]
	fmt.Println("   Name: ", definition.Name)
	fmt.Println("   Group: ", definition.Group)
	fmt.Println("   ShowGroup: ", definition.ShowGroup)
	fmt.Println("   Version: ", definition.Version)
	fmt.Println("   Kind: ", definition.Kind)
	fmt.Println("   InToc: ", definition.InToc)
	fmt.Println("   IsInlined: ", definition.IsInlined)
	fmt.Println("   IsOldVersion: ", definition.IsOldVersion)
	fmt.Println("   FoundInField: ", definition.FoundInField)
	fmt.Println("   FoundInOperation: ", definition.FoundInOperation)
}

func PrintDefinitionByVersionKindKeys(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("Definition ByVersionKind keys")

	for k, _ := range config.Definitions.ByGroupVersionKind {
		fmt.Println(k)
	}
}

func PrintDefinitionByKindKeys(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("Definition ByKind keys")

	for k, _ := range config.Definitions.ByKind {
		fmt.Println(k)
	}
}

func PrintDefinitionSchema(d *Definition) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("Definition schema")
	fmt.Println("   Name: ", d.Name)

	for p := range d.schema.Properties {
		fmt.Println("   ", p)
	}
}

func PrintAllDefinitionVersions(config *Config) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("All definition versions")

	count := 0

	for k, slc := range config.Definitions.ByKind {
		count = count + 1
		fmt.Println()
		fmt.Println("   ", k, count)

		for _, d := range slc {
			fmt.Println("        ", d.Key())
		}
	}
}

func PrintDefinitionVersions(config *Config, kind string) {
	fmt.Println()
	fmt.Println("----------------------------")
	fmt.Println("Definition versions")

	definitions := config.Definitions.ByKind[kind]
	for _, d := range definitions {
		fmt.Println("   ", d.Key())
	}
}

func PrintResource(r *Resource) {
	fmt.Println()
	fmt.Println("   -----------------------------")
	fmt.Println("   Resource")
	fmt.Println("      Name: ", r.Name)
	fmt.Println("      Version: ", r.Version)
	fmt.Println("      Group: ", r.Group)

	if r.Definition != nil {
		fmt.Println("      Definition key: ", r.Definition.Key())
	}

	for _, idef := range r.InlineDefinition {
		fmt.Println("         InlineDefinition: ", idef) 
	}
}

func PrintResourceCategory(rc *ResourceCategory) {
	fmt.Println()
	fmt.Println("-----------------------------")
	fmt.Println("Resource category")
	fmt.Println("   Category name: ", rc.Name)

	for _, r := range rc.Resources {
		PrintResource(r)
	}
}

func PrintOperations(config *Config) {
	fmt.Println()
	fmt.Println("PrintOperations")

	ops := config.Operations

	var yesFriendlyName []string
	var noFriendlyName []string

	for _, op := range ops {
		if op.FriendlyName == "" {
			noFriendlyName = append(noFriendlyName, "\"" + op.ID + "\": \"TODO\"")
		} else {
			yesFriendlyName = append(yesFriendlyName, "\"" + op.ID + "\": " + "\"" + op.FriendlyName + "\"")
		}
	}

	for _, op := range yesFriendlyName {
		fmt.Println(op)
	}

	for _, op := range noFriendlyName {
		fmt.Println(op)
	}

	fmt.Println(len(ops))
	fmt.Println(len(yesFriendlyName))
	fmt.Println(len(noFriendlyName))
}
