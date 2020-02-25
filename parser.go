package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

// APIURL constant API url getting routes informations
const APIURL = "http://router.project-osrm.org/route/v1/driving/%s;%s?overview=false"

// Point structure used for keeping single coordinates variable
type Point struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// String function returns a string represtentation of this structure
func (p *Point) String() string {
	return fmt.Sprintf("%s,%s", p.Lat, p.Lng)
}

// NewPoint - returns pointer to newly created Point struct
func NewPoint(s string) *Point {
	coordStr := strings.Split(s, ",")
	if len(coordStr) != 2 {
		return nil
	}

	p := Point{
		Lat: coordStr[0],
		Lng: coordStr[1],
	}

	return &p
}

// ShippingRoute structure used for keeping calculated routes from source to given destination point
type ShippingRoute struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}

// NewShippingRoute - returns pointer to newly created ShippingRoute struct
func NewShippingRoute(dest Point, duration float64, distance float64) *ShippingRoute {

	sr := ShippingRoute{
		Destination: dest.String(),
		Duration:    duration,
		Distance:    distance,
	}

	return &sr
}

// RawResponse structure used to deserialize raw JSON response data into object
type RawResponse struct {
	Routes []struct {
		Duration float64
		Distance float64
	}
	Code    string
	Message string
}

func callAPI(wg *sync.WaitGroup, source string, destination *Point, queue chan ShippingRoute) {
	defer wg.Done()

	url := fmt.Sprintf(APIURL, source, destination.String())

	// Request call
	res, reqErr := http.Get(url)
	if reqErr != nil {
		log.Printf("error: %v", reqErr)
	}

	// Read body
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Printf("error: %v", readErr)
	}

	defer res.Body.Close()

	var r RawResponse

	// Deserialize response body
	jsonErr := json.Unmarshal(body, &r)
	if jsonErr != nil {
		log.Printf("error: %v", jsonErr)
	}

	// Check if API returned correct status
	if r.Code == "Ok" {

		// Create driving details object and add to queue channel
		queue <- *(NewShippingRoute(*(destination), r.Routes[0].Duration, r.Routes[0].Distance))
	}
}

func getDrivingDetails(source *Point, destinations []*Point) *[]ShippingRoute {

	routes := []ShippingRoute{}

	queue := make(chan ShippingRoute)
	done := make(chan bool)

	// Dispatch listener to provide non-blocking write operations
	go func() {
		for {
			select {
			case r := <-queue:
				routes = append(routes, r)
			case <-done:
				return
			}
		}
	}()

	var wg sync.WaitGroup

	// Dispatch API calls in goroutines for each route
	for _, destination := range destinations {
		wg.Add(1)
		go callAPI(&wg, source.String(), destination, queue)
	}

	wg.Wait()

	done <- true

	return &routes
}
