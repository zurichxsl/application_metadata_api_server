package cache

import (
	"application_metadata_api_server/server/api"
	"fmt"
	"reflect"
)

type Path struct {
	value  string
	fields []string
}

// GetPaths get a list of paths from searchRoot to the leaves
func GetPaths(src map[string]interface{}) []Path {
	rs := make([]Path, 0)
	for k, obj := range src {
		if obj == nil {
			continue
		}
		value := reflect.ValueOf(obj)
		kind := value.Type().Kind()
		switch kind {
		case reflect.String:
			if len(value.String()) > 0 {
				rs = append(rs, Path{
					value:  value.String(),
					fields: []string{k},
				})
			}
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				subPathRS := GetPaths(value.Index(i).Interface().(map[string]interface{}))
				for _, subPath := range subPathRS {
					rs = append(rs, Path{
						value:  subPath.value,
						fields: append([]string{k}, subPath.fields...),
					})
				}
			}
		case reflect.Map:
			subPathRS := GetPaths(value.Interface().(map[string]interface{}))
			for _, subPath := range subPathRS {
				rs = append(rs, Path{
					value:  subPath.value,
					fields: append([]string{k}, subPath.fields...),
				})
			}
		default:
			fmt.Printf("%s is not handled\n", kind)
		}
	}
	return rs
}

// intersect take the intersection of two slices
func intersect(slice1 []api.Id, slice2 []api.Id) []api.Id {
	m := make(map[api.Id]int)
	n := make([]api.Id, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}
	return n
}

func exist(target api.Id, srcList []api.Id) bool {
	for _, src := range srcList {
		if src == target {
			return true
		}
	}
	return false
}
