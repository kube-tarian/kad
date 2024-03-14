package source

import (
	"fmt"
	"sort"
	"strings"
)

type BinData struct {
	// FileNames array should be in ascending order of sequence numbers provided
	// in FileMap keys of SQL data
	FileNames []string
	FilesMap  map[string][]byte
}

func NewBinData(data map[string][]byte) *BinData {
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return &BinData{
		FileNames: keys,
		FilesMap:  data,
	}
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func (b *BinData) Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := b.FilesMap[cannonicalName]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}
