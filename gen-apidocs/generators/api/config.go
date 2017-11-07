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
//	"strings"
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
		fmt.Println()
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
		fmt.Println("Got operations from API spec.")
	} else {
		fmt.Println("Failed to get operations from API spec.")
	}

	config.Definitions = getDefinitionsFromApiSpec(specs)

	if config.Definitions.ByGroupVersionKind != nil && config.Definitions.ByKind != nil {
		fmt.Println("Got definitions from API spec.")
	} else {
		fmt.Println("Failed to get definitions from API spec.")
	}

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
