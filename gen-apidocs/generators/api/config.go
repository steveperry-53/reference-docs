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
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
//	"log"
	"os"
	"path/filepath"
//	"regexp"
//	"sort"
	"strings"
//	"unicode"

	"github.com/go-openapi/loads"
)

var AllowErrors = flag.Bool("allow-errors", false, "If true, don't fail on errors.")
var ConfigDir = flag.String("config-dir", "", "Directory contain api files.")
var MungeGroups = flag.Bool("munge-groups", true, "If true, munge the group names for the operations to match.")

func NewConfig() *Config {

	var specs []*loads.Document = nil

	specs = LoadOpenApiSpec()

	if specs != nil {
		fmt.Println("Loaded API spec.")
	} else {
		fmt.Println("Failed to load API spec")
		os.Exit(1)
	}

	var friendlyOperationNames FriendlyOperationNames = nil

	friendlyOperationNames = loadFriendlyOperationNames()

	if friendlyOperationNames != nil {
		fmt.Println("Loaded friendly operation names", len(friendlyOperationNames))
	} else {
		fmt.Println("Failed to load friendly operation names.")
		os.Exit(1)
	}

	var friendlyNamesVerified bool = false

	friendlyNamesVerified = checkFriendlyOperationNames(specs, friendlyOperationNames)
	
	if friendlyNamesVerified {
		fmt.Println("Checked friendly operation names.")
	} else {
		fmt.Println("Some operations do not have friendly names.")
		fmt.Println("TODO: Give a list of those operations.")
		os.Exit(1)
	}

	var config *Config

	config = loadYamlConfig()

	if config != nil {
		fmt.Println("Loaded config.yaml.")
	} else {
		fmt.Println("Failed to load config.yaml.")
		os.Exit(1)
	}

	config.Operations = nil

	config.Operations = getOperationsFromApiSpec(specs)

	if config.Operations != nil {
		fmt.Println("Loaded operations from API spec.")
	} else {
		fmt.Println("Failed to load operations from API spec.")
	}

	config.Definitions = getDefinitionsFromApiSpec(specs)

	if config.Definitions.ByGroupVersionKind != nil && config.Definitions.ByKind != nil {
		fmt.Println("Loaded  definitions from API spec.")
	} else {
		fmt.Println("Failed to load definitions from API spec.")
	}

	config.initOperationCategories(friendlyOperationNames)
	fmt.Println("Initialized operation categories")
	fmt.Println("   Temporary verification:", config.OperationCategories[1])

	// What's next?
	// Finishing initializing operations.
	// We have operation categories and operatio friendly names.
	// What else do we need?
	//   Each OperationCategory has an Operations field that points to a []*Operation.
	//   We need to create those slices and hook them up.
	//   For each Operation in our map, get the key, which is the operation ID.
	//   Use the operation ID to get the operation friendly name from the map of friendlyOperationNames.
	//   Say I look an an Operation and figure out that its friendly name is Create.
	//   Now I need to find its operation category.
	//      Go through the slice of OperationCategory.
	//      For each category, go through the OperationFriendlyNames.
	//      If you find a match, you know the category.
	//         For that OperationCategory append a *Operation to its Operations field. 

	return config
}

func loadYamlConfig() *Config {

	f := filepath.Join(*ConfigDir, "config.yaml")

	config := &Config{}

	contents, err := ioutil.ReadFile(f)

	if err != nil {
		fmt.Printf("Failed to read yaml file %s: %v", f, err)
		os.Exit(2)
	}  else {
		err = yaml.Unmarshal(contents, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return config
}

func loadFriendlyOperationNames()FriendlyOperationNames {

	f := filepath.Join(*ConfigDir, "config-op-names.yaml")

	if contents, err := ioutil.ReadFile(f); err == nil {

		friendlyOperationNames := FriendlyOperationNames{}

		if err = yaml.Unmarshal(contents, &friendlyOperationNames); err == nil {

			return friendlyOperationNames

		} else {
			fmt.Println("Failed to unmarshal YAML", err)
		}
	} else {
		fmt.Println("Failed to read config-op-names.yaml", err)
	}

	return nil
}

func checkFriendlyOperationNames(specs []*loads.Document, friendlyNames FriendlyOperationNames) bool {

	VisitOperationsInApiSpec(specs, func(operation Operation) {
		if _, ok := friendlyNames[operation.ID]; !ok {
			fmt.Println("No friendly name found for", operation.ID)
		}
	})

	return true
}

func getOperationsFromApiSpec(specs []*loads.Document) Operations {

	o := Operations{}
	VisitOperationsInApiSpec(specs, func(operation Operation) {
		o[operation.ID] = &operation
	})

	return o
}

func getDefinitionsFromApiSpec(specs []*loads.Document) Definitions {

	d := Definitions{
		ByGroupVersionKind: map[string]*Definition{},
		ByKind: map[string]SortDefinitionsByVersion{},
	}

	VisitDefinitionsInApiSpec(specs, func(definition *Definition) {
		d.Put(definition)
	})

	return d
}

func (config *Config) initOperationCategories(frOpNames FriendlyOperationNames) {

	writeCategory := OperationCategory{
		Name: "Write Operations",
		FriendlyNames: []FriendlyOperationName{
			FriendlyOperationName{
				Name: "Create",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Patch",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Replace",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Delete",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Delete Collection",
				Operations: []*Operation{},
			},
		},
	}

	readCategory := OperationCategory{
		Name: "Read Operations:",
		FriendlyNames: []FriendlyOperationName{
			FriendlyOperationName{
				Name: "Read",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "List",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "List All Namespaces",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Watch",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Watch List",
				Operations: []*Operation{},
			},
			FriendlyOperationName{
				Name: "Watch List All Namespaces",
				Operations: []*Operation{},
			},
		},
	}

	config.OperationCategories = append(config.OperationCategories, writeCategory)
	config.OperationCategories = append(config.OperationCategories, readCategory)

	var opID string
	var opPointer *Operation

	for opID, opPointer = range config.Operations {

		friendlyName, found := frOpNames[opID]

		if found {

			// We have the friendly name of the operation and a pointer to the Operation.
			fmt.Println("Found", friendlyName)
			if opPointer != nil {

				var opCat OperationCategory

				for _, opCat = range config.OperationCategories {

					var opFrName FriendlyOperationName

					for _, opFrName = range opCat.FriendlyNames {

						if strings.Compare(opFrName.Name, friendlyName) == 0 {
							fmt.Println("opFrName.Name matches friendlyName.")
							if opPointer != nil {
								fmt.Println(".")
							}
						}
					}
				}
			}

		} else {
			fmt.Println("Friendly name not found for", opID)
		}
	}
}
