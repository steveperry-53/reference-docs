
package main

import (
	"fmt"
//	"os"
	"github.com/go-openapi/loads"
)

func main() {
	fmt.Println("Hello: test-spec")

	_, err := loads.JSONSpec("openapi-spec/swagger-test.json")
	
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Success")
	}
}
