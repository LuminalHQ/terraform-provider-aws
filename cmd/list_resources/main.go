package main

import (
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

type TypesList struct {
	Types []string
}

func main() {
	pi := provider.Provider()

	types := TypesList{}
	for k := range pi.ResourcesMap {
		types.Types = append(types.Types, k)
	}

	data, err := json.MarshalIndent(types, "", "  ")
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(data)
}
