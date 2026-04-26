package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/api"
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/data"
	"github.com/Francisco-Guillermo-Hernandez/reverse-geocode-service/internal/search"
)

var g *api.GeocodeHandler

func init() {
	locationsPath := "./data/geocode.json"
	countriesPath := "./data/countries.csv"
	locations, err := data.PrepareData(locationsPath, countriesPath)
	if err != nil {
		panic("Failed to load data in init: " + err.Error())
	}

	searcher := search.NewKDTreeGeocoder(locations)
	g = api.NewGeocodeHandler(searcher)
}

func LambdaHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	req, err := http.NewRequestWithContext(ctx, request.RequestContext.HTTP.Method, request.RawPath, nil)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Invalid request"}, err
	}

	req.Header = make(http.Header)
	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}
	q := req.URL.Query()
	for k, v := range request.QueryStringParameters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	switch request.RawPath {
	case "/geocode":
		g.HandleGeocode(w, req)
	case "/search-city":
		g.HandleSearchCity(w, req)
	default:
		return events.APIGatewayV2HTTPResponse{StatusCode: 404, Body: "Not Found"}, nil
	}

	bodyStr := w.Body.String()
	headerMap := make(map[string]string)
	for k, v := range w.Header() {
		headerMap[k] = v[0] 
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: w.Code,
		Headers:    headerMap,
		Body:       bodyStr,
	}, nil
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(LambdaHandler)
	}
}
