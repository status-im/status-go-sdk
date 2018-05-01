package sdk

import (
	"encoding/json"
)

func unmarshalJSON(j string) interface{} {
	var v interface{}
	// TODO(adriacidre) Handle this error
	_ = json.Unmarshal([]byte(j), &v)
	return v
}
