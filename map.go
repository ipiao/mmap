package mmap

import (
	"encoding/json"
	"strconv"
	"strings"
)

type AutoModel struct {
	Value string
}

type Map map[string]interface{}

func (m Map) Get(key string) interface{} {
	if len(key) == 0 {
		return nil
	}
	rs := strings.Split(key, ".")

	ret := m[rs[0]]
	for _, r := range rs[1:] {
		rm, ok := ret.(map[string]interface{})
		if ok {
			ret = rm[r]
		} else {
			return nil
		}
	}
	return ret
}

func (m Map) GetString(key string) string {
	i := m.Get(key)
	switch i.(type) {
	case string:
		return i.(string)
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', 6, 64)
	case json.Number:
		return i.(json.Number).String()
	default:
		return ""
	}
}

func (m Map) GetInt(key string) int {
	i := m.Get(key)
	switch i.(type) {
	case float64:
		return int(i.(float64))
	case json.Number:
		n, _ := i.(json.Number).Int64()
		return int(n)
	case string:
		n, _ := strconv.Atoi(i.(string))
		return n
	default:
		return 0
	}
}

func (m Map) GetInt64(key string) int64 {
	i := m.Get(key)
	switch i.(type) {
	case float64:
		return int64(i.(float64))
	case json.Number:
		n, _ := i.(json.Number).Int64()
		return n
	case string:
		n, _ := strconv.ParseInt(i.(string), 10, 64)
		return n
	default:
		return 0
	}
}

func (m Map) GetFloat64(key string) float64 {
	i := m.Get(key)
	switch i.(type) {
	case float64:
		return i.(float64)
	case json.Number:
		n, _ := i.(json.Number).Float64()
		return n
	case string:
		n, _ := strconv.ParseFloat(i.(string), 64)
		return n
	default:
		return 0
	}
}
