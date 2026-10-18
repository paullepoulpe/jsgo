// Command jsgo-worker accepts a stream of JSON requests and emits one JSON
// response per request. Arguments and generated files use the ordinary jsgo CLI
// contract; stdout/stderr/stdin are base64 so every supported charset survives.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/paullepoulpe/jsgo/internal/cliworker"
	"github.com/paullepoulpe/jsgo/internal/parsing"
)

type request struct {
	ID    json.RawMessage `json:"id,omitempty"`
	Args  []string        `json:"args"`
	Stdin []byte          `json:"stdin_base64,omitempty"`
}
type response struct {
	ID json.RawMessage `json:"id,omitempty"`
	cliworker.Result
}

func main() {
	parsing.ParserRunner_EnableBuiltinExternCache(true)
	decoder, encoder := json.NewDecoder(os.Stdin), json.NewEncoder(os.Stdout)
	for {
		var next request
		if err := decoder.Decode(&next); err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintln(os.Stderr, "jsgo-worker: invalid request:", err)
			os.Exit(2)
		}
		result := response{ID: next.ID, Result: cliworker.Run(next.Args, next.Stdin)}
		if err := encoder.Encode(result); err != nil {
			fmt.Fprintln(os.Stderr, "jsgo-worker: write response:", err)
			os.Exit(2)
		}
	}
}
