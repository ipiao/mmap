package mmap

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

type Map map[string]interface{}

func CreateMap(i interface{}) Map {
	return Map(Struct2Map(i))
}

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

func (m Map) GetBool(key string) bool {
	i := m.Get(key)
	switch i.(type) {
	case float64:
		return int64(i.(float64)) == 1
	case json.Number:
		n, _ := i.(json.Number).Int64()
		return n == 1
	case string:
		return i.(string) == "true"
	case bool:
		return i.(bool)
	default:
		return false
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

func (m Map) GetUint64(key string) uint64 {
	i := m.Get(key)
	switch i.(type) {
	case float64:
		return uint64(i.(float64))
	case json.Number:
		n, _ := i.(json.Number).Int64()
		return uint64(n)
	case string:
		n, _ := strconv.ParseUint(i.(string), 10, 64)
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

func (m Map) ToStruct(obj interface{}) error {
	return m.ToStructWithTagName(obj, "json")
}

func (m Map) ToStructWithTagName(obj interface{}, tagName string) error {
	v := reflect.ValueOf(obj)
	return m.mapToStruct(v, tagName)
}

func (m Map) mapToStruct(v reflect.Value, tagName string) error {
	if m != nil {
		k := v.Kind()
		if k == reflect.Ptr || k == reflect.Interface {
			v = v.Elem()
		}
		t := v.Type()
		if v.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				if v.Field(i).CanSet() {
					tag := t.Field(i).Tag.Get(tagName)
					if tag == "-" {
						continue
					} else if tag == "" {
						tag = SnakeName(t.Field(i).Name)
					}
					switch v.Field(i).Kind() {
					case reflect.Bool:
						v.Field(i).SetBool(m.GetBool(tag))
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						v.Field(i).SetInt(m.GetInt64(tag))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						v.Field(i).SetUint(m.GetUint64(tag))
					case reflect.String:
						v.Field(i).SetString(m.GetString(tag))
					case reflect.Float32, reflect.Float64:
						v.Field(i).SetFloat(m.GetFloat64(tag))
					case reflect.Ptr:
						// m.Get()
					}
					// // 这是没有支持基础类型指针
					// if (v.Field(i).Kind() == reflect.Struct || v.Field(i).Kind() == reflect.Ptr) && v.Field(i).Type().String() != "time.Time" {
					// 	fm, _ := data[tag].(map[string]interface{})
					// 	map2struct(fm, v.Field(i), tagName)
					// } else {

					// }
				}

			}
		}
	}
	return nil
}
