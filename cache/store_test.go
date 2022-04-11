package cache

import (
	"application_metadata_api_server/server/api"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestStoreImpl_Search(t *testing.T) {
	testCases := []struct {
		name        string
		apps        []api.App
		queryFields []string
		expected    map[string][]api.Id
	}{
		{name: "search on immediate leaf node on field title",
			apps: []api.App{
				{
					Title: "t1",
					Id:    "1",
				},
				{
					Title: "t2",
					Id:    "2",
				},
				{
					Title: "T2",
					Id:    "3",
				},
				{
					Title: "t2 abc",
					Id:    "4",
				},
				{
					Title: "t4 abc",
					Id:    "5",
				},
				{
					Title: "t4 t2 abcd",
					Id:    "6",
				},
				{
					Title: "t4 t2 abcd",
					Id:    "7",
				},
			},
			queryFields: []string{"title"},
			expected: map[string][]api.Id{
				"t1":  {"1"},
				"t2":  {"2", "3", "4", "6", "7"},
				"abc": {"4", "5"},
				"t4":  {"5", "6", "7"},
				"t":   {},
				"zyc": {},
			},
		},
		{name: "search on nested leaf node on field users.name",
			apps: []api.App{
				{
					Maintainers: []api.Maintainer{
						{Name: "abx",
							Email: "a@b.com"},
						{Name: "xyz",
							Email: "c@d.com"},
						{Name: "bob david",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "1",
				},
				{
					Maintainers: []api.Maintainer{
						{Name: "a",
							Email: "a@b.com"},
						{Name: "bob samuel",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "2",
				},
			},
			queryFields: []string{"maintainers", "name"},
			expected: map[string][]api.Id{
				"a":          {"2"},
				"abx":        {"1"},
				"david":      {"1"},
				"bob samuel": {"2"},
				"bob":        {"1", "2"},
				"mary":       {"1", "2"},
				"zzz":        {},
				"bobdavid":   {},
			},
		},
		{name: "search on nested leaf node on field release.author.email",
			apps: []api.App{
				{
					Release: api.Release{
						Name: "ab",
						Author: api.Maintainer{
							Email: "a@b.com",
						},
					},
					Id: "1",
				},
				{
					Release: api.Release{
						Name: "ab",
						Author: api.Maintainer{
							Email: "a@b.com",
						},
					},
					Id: "2",
				},
				{
					Release: api.Release{
						Name: "cd",
						Author: api.Maintainer{
							Email: "c@d.com",
						},
					},
					Id: "3",
				},
			},
			queryFields: []string{"release", "author", "email"},
			expected: map[string][]api.Id{
				"c@d.com":    {"3"},
				"c@d.comddd": {},
				"a@b.com":    {"1", "2"},
			},
		},
		{name: "search on nested array on field users.name",
			apps: []api.App{
				{
					Maintainers: []api.Maintainer{
						{Name: "abx",
							Email: "a@b.com"},
						{Name: "xyz",
							Email: "c@d.com"},
						{Name: "bob david",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "1",
				},
				{
					Maintainers: []api.Maintainer{
						{Name: "a",
							Email: "a@b.com"},
						{Name: "bob samuel",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "2",
				},
			},
			queryFields: []string{"maintainers", "name"},
			expected: map[string][]api.Id{
				"a":          {"2"},
				"abx":        {"1"},
				"david":      {"1"},
				"bob":        {"1", "2"},
				"mary":       {"1", "2"},
				"zzz":        {},
				"bob david":  {"1"},
				"bob david2": {},
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			tree := InitStore()
			// add to the search space
			for _, app := range test.apps {
				data, err := yaml.Marshal(&app)
				assert.Nil(t, err)
				_, err = tree.Add(&app, data)
				assert.Nil(t, err)
			}
			// search
			for query, expectedRS := range test.expected {
				t.Log(fmt.Sprintf("querying %s, expecting rs len %d", query, len(expectedRS)))
				assert.Equal(t, expectedRS, tree.Search(query, test.queryFields...))
			}
		})
	}
}

func TestStoreImpl_SearchStruct(t *testing.T) {
	testCases := []struct {
		name     string
		apps     []api.App
		query    api.App
		expected []api.Id
	}{
		{
			name: "search on 1 simple field",
			apps: []api.App{
				{
					Title: "t1",
					Id:    "1",
				},
				{
					Title: "t2",
					Id:    "2",
				},
			},
			query: api.App{
				Title: "t1",
			},
			expected: []api.Id{
				"1",
			},
		},
		{
			name: "substring search on 1 simple field",
			apps: []api.App{
				{
					Title: "t1 abc efg",
					Id:    "1",
				},
				{
					Title: "t2",
					Id:    "2",
				},
			},
			query: api.App{
				Title: "abc",
			},
			expected: []api.Id{
				"1",
			},
		},
		{
			name: "substring search on 1 simple field with uppercase",
			apps: []api.App{
				{
					Title: "t1 AbC efg",
					Id:    "1",
				},
				{
					Title: "t2",
					Id:    "2",
				},
			},
			query: api.App{
				Title: "abc",
			},
			expected: []api.Id{
				"1",
			},
		},
		{
			name: "substring search on 1 simple field with 1 common substring",
			apps: []api.App{
				{
					Title: "t1 abc efg",
					Id:    "1",
				},
				{
					Title: "t2 abc",
					Id:    "2",
				},
			},
			query: api.App{
				Title: "abc",
			},
			expected: []api.Id{
				"1", "2",
			},
		},
		{
			name: "substring search on multiple nested fields, positive find",
			apps: []api.App{
				{
					Release: api.Release{
						Name: "ab",
						Author: api.Maintainer{
							Email: "c@d.com",
						},
					},
					Id: "1",
				},
				{
					Release: api.Release{
						Name: "cd",
						Author: api.Maintainer{
							Email: "a@d.com",
						},
					},
					Id: "2",
				},
				{
					Release: api.Release{
						Name: "cd",
						Author: api.Maintainer{
							Email: "c@d.com",
						},
					},
					Id: "3",
				},
			},
			query: api.App{
				Release: api.Release{
					Name: "cd",
					Author: api.Maintainer{
						Email: "c@d.com",
					},
				},
			},
			expected: []api.Id{
				"3",
			},
		},
		{
			name: "substring search on 1 nested field, positive find",
			apps: []api.App{
				{
					Release: api.Release{
						Name: "ab",
						Author: api.Maintainer{
							Email: "c@d.com",
						},
					},
					Id: "1",
				},
				{
					Release: api.Release{
						Name: "cd",
						Author: api.Maintainer{
							Email: "a@d.com",
						},
					},
					Id: "2",
				},
				{
					Release: api.Release{
						Name: "cd",
						Author: api.Maintainer{
							Email: "c@d.com",
						},
					},
					Id: "3",
				},
			},
			query: api.App{
				Release: api.Release{
					Name: "cd",
				},
			},
			expected: []api.Id{
				"2", "3",
			},
		},
		{
			name: "substring search on nested field, negative find",
			apps: []api.App{
				{
					Release: api.Release{
						Name: "ab",
						Author: api.Maintainer{
							Email: "c@d.com",
						},
					},
					Id: "1",
				},
				{
					Release: api.Release{
						Name: "cd",
						Author: api.Maintainer{
							Email: "a@d.com",
						},
					},
					Id: "2",
				},
			},
			query: api.App{
				Release: api.Release{
					Name: "cd",
					Author: api.Maintainer{
						Email: "c@d.com",
					},
				},
			},
			expected: []api.Id{},
		},
		{
			name: "substring search on slice, positive find",
			apps: []api.App{
				{
					Maintainers: []api.Maintainer{
						{Name: "abx",
							Email: "a@b.com"},
						{Name: "xyz",
							Email: "c@d.com"},
						{Name: "bob david",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "1",
				},
				{
					Maintainers: []api.Maintainer{
						{Name: "a",
							Email: "a@b.com"},
						{Name: "bob samuel",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "2",
				},
				{
					Maintainers: []api.Maintainer{},
					Id:          "3",
				},
			},
			query: api.App{
				Maintainers: []api.Maintainer{
					{Name: "bob",
						Email: "a@b.com"},
				},
			},
			expected: []api.Id{
				"1", "2",
			},
		},
		{
			name: "substring search on slice, negative find",
			apps: []api.App{
				{
					Maintainers: []api.Maintainer{
						{Name: "abx",
							Email: "a@b.com"},
						{Name: "xyz",
							Email: "c@d.com"},
						{Name: "bob david",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "1",
				},
				{
					Maintainers: []api.Maintainer{
						{Name: "a",
							Email: "a@b.com"},
						{Name: "bob samuel",
							Email: "a@b.com"},
						{Name: "mary",
							Email: "c@d.com"},
					},
					Id: "2",
				},
				{
					Maintainers: []api.Maintainer{},
					Id:          "3",
				},
			},
			query: api.App{
				Maintainers: []api.Maintainer{
					{Name: "sam",
						Email: "a@b.com"},
				},
			},
			expected: []api.Id{},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			tree := InitStore()
			// add to the search space
			for _, app := range test.apps {
				data, err := yaml.Marshal(&app)
				assert.Nil(t, err)
				_, err = tree.Add(&app, data)
				assert.Nil(t, err)
			}
			// search struct
			t.Log(fmt.Sprintf("querying app struct %+v, expecting rs len %d", test.query, len(test.expected)))
			rs, err := tree.SearchStruct(&test.query)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, rs)
		})
	}
}
