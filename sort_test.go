package main

import (
	"testing"
)

func TestSortResultsEachIsDifferent(t *testing.T) {
	routes := []ShippingRoute{
		*(NewShippingRoute(*(NewPoint("13.397634,52.529407")), 606, 1884.9)),
		*(NewShippingRoute(*(NewPoint("13.428555,52.523219")), 691.6, 4128.3)),
		*(NewShippingRoute(*(NewPoint("13.397633,52.529599")), 615.6, 1900)),
		*(NewShippingRoute(*(NewPoint("13.397634,52.529598")), 615.6, 1896.8)),
		*(NewShippingRoute(*(NewPoint("13.397632,52.529488")), 606.2, 1885.5)),
	}

	sorted := SortResults(routes)

	for i, r := range []struct {
		Duration float64
		Distance float64
	}{
		{606, 1884.9},
		{606.2, 1885.5},
		{615.6, 1896.8},
		{615.6, 1900},
		{691.6, 4128.3},
	} {
		if sorted[i].Duration != r.Duration {
			t.Errorf("sorted results element %d should has %f duration instead of %f.", i, r.Duration, sorted[i].Duration)
		}

		if sorted[i].Distance != r.Distance {
			t.Errorf("sorted results element %d should has %f distance instead of %f.", i, r.Distance, sorted[i].Distance)
		}
	}
}
