package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type kv struct {
	key   string
	value int
}

func Top10(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	fields := strings.Fields(text)
	cache := make(map[string]int)

	for _, field := range fields {
		value, ok := cache[field]
		if ok {
			cache[field] = value + 1
		} else {
			cache[field] = 1
		}
	}

	kvs := make([]kv, 0, len(cache))

	for k, v := range cache {
		kvs = append(kvs, kv{k, v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].value == kvs[j].value {
			return kvs[i].key < kvs[j].key
		}
		return kvs[i].value > kvs[j].value
	})

	results := make([]string, 0)

	for _, kv := range kvs {
		results = append(results, kv.key)
	}

	if len(results) < 10 {
		return results
	}
	return results[:10]
}
