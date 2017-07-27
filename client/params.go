package client

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/cybozu-go/kkok"
)

func printParams(pp kkok.PluginParams, showJSON bool) {
	if showJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		enc.Encode(pp)
		return
	}

	fmt.Println("type:", pp.Type)
	keys := make([]string, 0, len(pp.Params))
	for k := range pp.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s: %v\n", k, pp.Params[k])
	}
}
