package utils

import (
	"bytes"
	"encoding/json"
)

func ToPrettyJson(obj interface{}) string {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "\t")
	_ = encoder.Encode(obj)
	return PrettifyJson(b.Bytes())
}
func PrettifyJson(data []byte) string {
	bb := new(bytes.Buffer)
	err := json.Indent(bb, data, "", "\t")
	if err != nil {
		return ""
	}
	return bb.String()
}
