package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

// BaseController structure
type BaseController struct{}

// IBaseController for exposed methods
type IBaseController interface {
	JSONResponse(w http.ResponseWriter, v interface{}, code int)
	ErrorJSONResponse(err error, w http.ResponseWriter, status ...int)
}

// APIController structure
type APIController struct {
	BaseController
}

// NewBaseController - returns pointer to newly created BaseController struct
func NewBaseController() *BaseController {
	return &BaseController{}
}

// NewAPIController - returns pointer to newly created APIController struct
func NewAPIController() *APIController {
	return &APIController{
		BaseController: BaseController{},
	}
}

// JSONResponse - function returns JSON response of any object
func (c *BaseController) JSONResponse(w http.ResponseWriter, v interface{}, code int) {

	// Marshal any object to JSON format
	content, marsharErr := json.Marshal(v)
	if marsharErr != nil {
		http.Error(w, marsharErr.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(content)
}

// ErrorJSONResponse - function returns error JSON message
func (c *BaseController) ErrorJSONResponse(err error, w http.ResponseWriter, status ...int) {

	// Prepare message
	msg := map[string]string{
		"status":  "fail",
		"message": err.Error(),
	}

	// Set proper return status or let 500 as default
	returnStatus := http.StatusInternalServerError
	if len(status) > 0 {
		returnStatus = status[0]
	}

	// Send back JSON response
	c.JSONResponse(w, &msg, returnStatus)
}

// HealthCheck - endpoint returns status and hostname
func (c *APIController) HealthCheck(w http.ResponseWriter, r *http.Request) {

	res := map[string]interface{}{
		"status": http.StatusOK,
	}

	c.JSONResponse(w, &res, http.StatusOK)
}

// Routes - endpoint takes the source and a list of destinations
// and returns a list of routes between source and each destination
func (c *APIController) Routes(w http.ResponseWriter, r *http.Request) {

	var routes *[]ShippingRoute

	// Extract data from query params
	source, destinations, err := extractQueryParams(r.URL.Query())
	if err != nil {
		c.ErrorJSONResponse(err, w, 422)
		return
	}

	// Call external API to retrieve shipping routes details
	routes = getDrivingDetails(source, destinations)

	// Sort results
	results := SortResults(*(routes))

	res := map[string]interface{}{
		"source": source.String(),
		"routes": results,
	}

	c.JSONResponse(w, &res, http.StatusOK)
}

func extractQueryParams(v url.Values) (source *Point, destinations []*Point, err error) {

	r, _ := regexp.Compile("[-+]?([0-9]*\\.[0-9]+|[0-9]+),[-+]?([0-9]*\\.[0-9]+|[0-9]+)")

	// Validate source coordinates param
	srcParams, hasSrcParam := v["src"]
	if !hasSrcParam || len(srcParams[0]) < 1 {
		return nil, nil, errors.New("missing src query param")
	}

	if r.MatchString(srcParams[0]) == false {
		return nil, nil, errors.New("src query param is incorrect format")
	}

	if len(srcParams) > 1 {
		return nil, nil, errors.New("only one src query param allowed")
	}

	srcPoint := NewPoint(srcParams[0])
	if srcPoint == nil {
		return nil, nil, errors.New("missing src query param")
	}

	// Validate destinations coordinates params
	var destsPoints []*Point

	destParams, hasDestParam := v["dst"]
	if !hasDestParam || len(destParams) < 1 {
		return nil, nil, errors.New("missing dst query param. It should be at least one")
	}

	// Validate multiple
	for i, d := range destParams {
		if len(d) < 1 {
			return nil, nil, errors.New(fmt.Sprintf("%d dst query param is empty", i))
		}

		if r.MatchString(d) == false {
			return nil, nil, errors.New(fmt.Sprintf("%d dst query param is incorrect format", i))
		}

		newPoint := NewPoint(d)
		if newPoint != nil {
			destsPoints = append(destsPoints, newPoint)
		}
	}

	return srcPoint, destsPoints, nil
}
