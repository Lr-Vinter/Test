package index

import "errors"

type Index struct {
	Indexmap map[string]int
}

func NewIndexMap() *Index {
	Indexmap := make(map[string]int)
	return &Index{Indexmap}
}

func (ix *Index) GetIndex(key string) (int, error) {
	if value, ok := ix.Indexmap[key]; ok {
		return value, nil
	}
	err := errors.New("Index : No such key")
	return 0, err
}

func (ix *Index) PushPair(key string, value int) {
	ix.Indexmap[key] = value
}
