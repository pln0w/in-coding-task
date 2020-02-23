package main

import (
	"sort"
)

// ByDurationOrDistance type implements sort.Interface
type ByDurationOrDistance []ShippingRoute

func (a ByDurationOrDistance) Len() int {
	return len(a)
}

func (a ByDurationOrDistance) Less(i, j int) bool {
	if a[i].Duration == a[j].Duration {
		return a[i].Distance < a[j].Distance
	}

	return a[i].Duration < a[j].Duration
}

func (a ByDurationOrDistance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// SortResults - a function soriting slice of ShippingRoute structs by time or distance
func SortResults(results []ShippingRoute) []ShippingRoute {

	sort.Sort(ByDurationOrDistance(results))

	return results
}
