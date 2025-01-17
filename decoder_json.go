package zerolog

import (
	"bytes"
	"encoding/json"
	"io"
)

func getValues(buf []byte, key ...string) (map[string][]interface{}, error) {
	if len(key) == 0 {
		return nil, nil
	}
	decoder := json.NewDecoder(bytes.NewReader(buf))
	decoder.UseNumber()
	kv := make(map[string][]interface{})
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return kv, nil
		}
		t, ok := token.(string)
		if !ok {
			continue
		}
		for _, k := range key {
			if t != k {
				continue
			}
			var val interface{}
			if err := decoder.Decode(&val); err != nil {
				return nil, err
			}
			if _, ok := kv[k]; !ok {
				kv[k] = make([]interface{}, 0)
			}
			kv[k] = append(kv[k], val)
		}
	}
}
