package mock

import (
	"encoding/json"
	"io"
	"strings"
)

func RequestBody(v map[string]any) (io.Reader, error) {
	data, err := json.Marshal(v)
	return strings.NewReader(string(data)), err
}
