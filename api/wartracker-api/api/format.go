package api

import (
	"encoding/json"
	"io"
)

func jsonOutput(w io.Writer, v any, indent bool) error {
	enc := json.NewEncoder(w)
	if indent {
		enc.SetIndent("", "    ")
	}
	return enc.Encode(v)
}
