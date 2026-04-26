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

func LambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, request.HTTPMethod, request.Path, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Invalid request"}, err
	}

	// Copiar headers de AWS a http.Request
	req.Header = make(http.Header)
	for key, value := range request.Headers {
		req.Header.Add(key, value)
	}
	// Map query parameters
	q := req.URL.Query()
	for k, v := range request.QueryStringParameters {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()

	// Route based on path
	switch request.Path {
	case "/geocode", "/api/v1/geocode":
		g.HandleGeocode(w, req)
	case "/search-city", "/api/v1/search-city":
		g.HandleSearchCity(w, req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Not Found"}, nil
	}

	// 6. Convertir la respuesta de Go a respuesta de AWS Lambda
	bodyStr := w.Body.String()
	headerMap := make(map[string]string)
	for k, v := range w.Header() {
		headerMap[k] = v[0] // Usamos el primer valor si hay múltiples
	}

	return events.APIGatewayProxyResponse{
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
