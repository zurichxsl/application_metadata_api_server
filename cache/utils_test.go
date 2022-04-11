package cache

import (
	"application_metadata_api_server/server/api"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestGetPaths(t *testing.T) {
	testCases := []struct {
		name     string
		app      api.App
		expected []Path
	}{
		{
			name: "simple 1 direct leaf",
			app: api.App{
				Title: "t1",
			},
			expected: []Path{
				{value: "t1", fields: []string{"title"}},
			},
		},
		{
			name: "simple 1 nested leaf",
			app: api.App{
				Release: api.Release{
					Author: api.Maintainer{
						Name: "a1",
					},
				},
			},
			expected: []Path{
				{value: "a1", fields: []string{"release", "author", "name"}},
			},
		},
		{
			name: "simple 2 direct leaves",
			app: api.App{
				Title:   "t1",
				Company: "some inc",
			},
			expected: []Path{
				{value: "t1", fields: []string{"title"}},
				{value: "some inc", fields: []string{"company"}},
			},
		},
		{
			name: "simple slice",
			app: api.App{
				Maintainers: []api.Maintainer{
					{Name: "first m1"},
					{Name: "sec m2"},
					{Name: "m3"},
				},
			},
			expected: []Path{
				{value: "first m1", fields: []string{"maintainers", "name"}},
				{value: "sec m2", fields: []string{"maintainers", "name"}},
				{value: "m3", fields: []string{"maintainers", "name"}},
			},
		},
		{
			name: "expect adds up from multiple branches",
			app: api.App{
				Title:   "t1",
				Company: "some inc",
				Maintainers: []api.Maintainer{
					{Name: "first m1"},
					{Name: "sec m2"},
				},
				Release: api.Release{
					Name: "r1",
					Author: api.Maintainer{
						Name: "a1",
					},
				},
			},
			expected: []Path{
				{value: "t1", fields: []string{"title"}},
				{value: "some inc", fields: []string{"company"}},
				{value: "first m1", fields: []string{"maintainers", "name"}},
				{value: "sec m2", fields: []string{"maintainers", "name"}},
				{value: "r1", fields: []string{"release", "name"}},
				{value: "a1", fields: []string{"release", "author", "name"}},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			src, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&test.app)
			assert.Nil(t, err)
			rs := GetPaths(src)
			for _, r := range rs {
				e := getPathByValue(r.value, test.expected)
				assert.Equal(t, e.fields, r.fields)
			}
		})
	}

}

func getPathByValue(value string, paths []Path) *Path {
	for i := range paths {
		if value == paths[i].value {
			return &paths[i]
		}
	}
	return nil
}
