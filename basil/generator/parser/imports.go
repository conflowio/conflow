package parser

import "sort"

func SortedImportKeys(imports map[string]string) []string {
	var keys []string
	for k := range imports {
		keys = append(keys, k)
	}

	is := importSorter{
		Keys:    keys,
		Imports: imports,
	}
	sort.Sort(is)

	return is.Keys
}

type importSorter struct {
	Keys    []string
	Imports map[string]string
}

func (i importSorter) Len() int {
	return len(i.Imports)
}

func (i importSorter) Less(a, b int) bool {
	return i.Imports[i.Keys[a]] < i.Imports[i.Keys[b]]
}

func (i importSorter) Swap(a, b int) {
	i.Keys[a], i.Keys[b] = i.Keys[b], i.Keys[a]
}
