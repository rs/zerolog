package zerolog

import (
	"fmt"
	"net"
)

func appendProp(props map[string][]interface{}, key string, val interface{}) map[string][]interface{} {
	if props == nil {
		props = make(map[string][]interface{})
	}
	if _, ok := props[key]; !ok {
		props[key] = make([]interface{}, 0)
	}
	props[key] = append(props[key], val)
	return props
}

func appendStringerProp(props map[string][]interface{}, key string, val fmt.Stringer) map[string][]interface{} {
	if val == nil {
		return appendProp(props, key, nil)
	}
	return appendProp(props, key, val.String())
}

func appendStringersProp(props map[string][]interface{}, key string, vals []fmt.Stringer) map[string][]interface{} {
	if vals == nil {
		return appendProp(props, key, nil)
	}
	for _, value := range vals {
		props = appendStringerProp(props, key, value)
	}
	return props
}

func appendIPProp(props map[string][]interface{}, key string, ip net.IP) map[string][]interface{} {
	return appendProp(props, key, ip.String())
}

func appendIPPrefixProp(props map[string][]interface{}, key string, pfx net.IPNet) map[string][]interface{} {
	return appendProp(props, key, pfx.String())
}

func appendMACAddrProp(props map[string][]interface{}, key string, ma net.HardwareAddr) map[string][]interface{} {
	return appendProp(props, key, ma.String())
}
