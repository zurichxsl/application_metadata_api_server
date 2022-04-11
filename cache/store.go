package cache

import (
	"application_metadata_api_server/server/api"
	"fmt"
	"strconv"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
)

// Store is the interface of the in-memory data store.
// stores Struct as a tree data structure, tree node = field name of the struct, tree node value = field values
type Store interface {
	// Add inserts an App to the in-memory data store
	Add(app *api.App, raw []byte) (api.Id, error)
	// Get gets an App based on its Id
	Get(id api.Id) ([]byte, error)
	// Search takes a value str and its field or its nested field
	Search(value string, fields ...string) []api.Id
	// SearchStruct takes an App struct, and traverse along the struct with store's tree structure
	SearchStruct(app *api.App) ([]api.Id, error)
}

type storeImpl struct {
	rwLock sync.RWMutex
	// searchRoot is the App search space
	searchRoot *TreeNode
	// rawData contains direct Id to App mapping
	rawData map[api.Id][]byte
	// cnt is used to generate auto increment ids : TODO: could use a gen-id service
	cnt int
}

func InitStore() Store {
	return &storeImpl{
		searchRoot: newTreeNode(""),
		rawData:    make(map[api.Id][]byte),
		cnt:        0}
}

func (t *storeImpl) Add(app *api.App, rawContent []byte) (api.Id, error) {
	t.rwLock.Lock()
	defer t.rwLock.Unlock()

	id := t.cnt + 1
	app.Id = api.Id(strconv.Itoa(id))

	// 1. add to raw, overwrite if exists
	t.rawData[app.Id] = rawContent
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(app)
	if err != nil {
		return "", err
	}
	// 2. add to search space
	t.searchRoot.addNode(app.Id, unstructuredObj)
	// only increase if add success
	t.cnt++
	return app.Id, nil
}

func (t *storeImpl) Get(id api.Id) ([]byte, error) {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	rawApp, ok := t.rawData[id]
	if ok {
		return rawApp, nil
	}
	return nil, fmt.Errorf("%v not found", id)
}

func (t *storeImpl) Search(value string, fields ...string) []api.Id {
	t.rwLock.RLock()
	defer t.rwLock.RUnlock()

	rs := make([]api.Id, 0)
	p := t.searchRoot
	for _, field := range fields {
		found := false
		for i := range p.children {
			if p.children[i].key == field {
				found = true
				p = p.children[i]
				break
			}
		}
		if !found {
			return rs
		}
	}
	return p.data.Search(value)
}

func (t *storeImpl) SearchStruct(app *api.App) ([]api.Id, error) {
	rs := make([]api.Id, 0)
	unstructuredSrc, err := runtime.DefaultUnstructuredConverter.ToUnstructured(app)
	if err != nil {
		return rs, err
	}
	paths := GetPaths(unstructuredSrc)
	if len(paths) == 0 {
		return rs, nil
	}
	// get the intersection of result from paths
	p1 := paths[0]
	r1 := t.Search(p1.value, p1.fields...)
	for i := 1; i < len(paths); i++ {
		p2 := paths[i]
		r2 := t.Search(p2.value, p2.fields...)
		r1 = intersect(r1, r2)
	}
	return r1, nil
}
