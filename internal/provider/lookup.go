package provider

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func lookupMapValue[K comparable](lookup map[K]string, value string) (K, error) {
	for k, v := range lookup {
		if v == value {
			return k, nil
		}
	}

	var zero K
	return zero, fmt.Errorf("value %s not found in lookup", value)
}

func mapValues[K comparable](lookup map[K]string) []string {
	return slices.Sorted(maps.Values(lookup))
}

func formatValues(values []string) string {
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

func listMapValues[K comparable](lookup map[K]string) string {
	return formatValues(mapValues(lookup))
}
