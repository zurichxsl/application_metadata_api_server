package cache

import (
	"application_metadata_api_server/server/api"
	"fmt"
	"reflect"
	"strings"
)

// TreeNode represents a data node in the search space
type TreeNode struct {
	// key is the schema field name
	key string
	// data contains all the possible values for the key
	data InvertedIndex
	// children is the immediate children of current tree node
	children map[string]*TreeNode
}

func newTreeNode(key string) *TreeNode {
	return &TreeNode{
		key:      key,
		data:     make(InvertedIndex),
		children: make(map[string]*TreeNode),
	}
}

func (p *TreeNode) addNode(appId api.Id, unstructured map[string]interface{}) {
	for k, obj := range unstructured {
		_, ok := p.children[k]
		if !ok {
			// k does not exist in the tree
			p.children[k] = newTreeNode(k)
		}
		if obj == nil {
			continue
		}
		value := reflect.ValueOf(obj)
		kind := value.Type().Kind()
		switch kind {
		case reflect.String:
			p.children[k].data.Add(appId, value.String())
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				m := value.Index(i).Interface().(map[string]interface{})
				p.children[k].addNode(appId, m)
			}
		case reflect.Map:
			p.children[k].addNode(appId, value.Interface().(map[string]interface{}))
		default:
			fmt.Printf("%s is not handled\n", kind)
		}
	}
}

// InvertedIndex represents an index data structure storing a mapping from content
// lowercase words to its Id in a document or a set of documents
type InvertedIndex map[string][]api.Id

// Add adds a str and this str's tokenized words(split by space) to its search space
// e.g: appId 1, and value "this is a Cat" is stored as: "this":[1], "is":[1], "a": [1], "cat": [1], "this is a cat": [1]
// and if "this a" from appId 2 gets added, it would be: "this":[1, 2], "is":[1], "a": [1,2], "cat": [1], "this is": [2]
func (x *InvertedIndex) Add(appId api.Id, value string) {
	value = getLowercase(value)
	words := strings.Fields(value)
	if len(words) > 1 {
		// store the full value
		valueRef, valueRefok := (*x)[value]
		if !valueRefok {
			valueRef = nil
		}
		if !exist(appId, valueRef) {
			(*x)[value] = append(valueRef, appId)
		}
	}
	// store the value token
	for _, word := range words {
		word = getLowercase(word)
		ref, ok := (*x)[word]
		if !ok {
			ref = nil
		}
		if !exist(appId, ref) {
			(*x)[word] = append(ref, appId)
		}
	}
}

// Search is the plain text search, and query is a plain text. (The nested query is handled in the tree data structure, not here)
// e.g. if you stored "this is acat", a query string of "this", "is", "acat" returns positive match,
// however a query string of "this is", "this acat" will result in not found.
func (x *InvertedIndex) Search(query string) []api.Id {
	q := getLowercase(query)
	ref, ok := (*x)[q]
	if ok {
		return ref
	}
	return []api.Id{}
}

func getLowercase(query string) string {
	return strings.ToLower(query)
}
