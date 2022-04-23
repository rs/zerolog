package zerolog

import (
	"sort"
	"unsafe"
)

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

func appendFields(dst []byte, fields interface{}) []byte {
	switch fields := fields.(type) {
	case []interface{}:
		if n := len(fields); n&0x1 == 1 { // odd number
			fields = fields[:n-1]
		}
		dst = appendFieldList(dst, fields)
	case map[string]interface{}:
		keys := make([]string, 0, len(fields))
		for key := range fields {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		kv := make([]interface{}, 2)
		for _, key := range keys {
			kv[0], kv[1] = key, fields[key]
			dst = appendFieldList(dst, kv)
		}
	}
	return dst
}
