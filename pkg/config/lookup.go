package config

import "maps"

type Lookup func(string) (string, bool)

func Map(values map[string]string) Lookup {
	owned := maps.Clone(values)
	return func(key string) (string, bool) { v, ok := owned[key]; return v, ok }
}
