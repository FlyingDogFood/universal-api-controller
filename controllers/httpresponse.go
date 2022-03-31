package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

type HttpResponse struct {
	Status     string `json:"Status"`
	StatusCode int    `json:"StatusCode"`
	Proto      string `json:"Proto"`
	ProtoMajor int    `json:"ProtoMajor"`
	ProtoMinor int    `json:"ProtoMinor"`

	Header http.Header `json:"Header"`

	Body interface{} `json:"Body"`

	ContentLength int64 `json:"ContentLength"`
}

func fromHttpResponse(response http.Response) (HttpResponse, error) {
	httpResponse := HttpResponse{
		Status:        response.Status,
		StatusCode:    response.StatusCode,
		Proto:         response.Proto,
		ProtoMajor:    response.ProtoMajor,
		ProtoMinor:    response.ProtoMinor,
		Header:        response.Header.Clone(),
		ContentLength: response.ContentLength,
	}
	defer response.Body.Close()
	for _, contentType := range response.Header.Values("Content-Type") {
		if contentType == "appication/json" {
			err := json.NewDecoder(response.Body).Decode(&httpResponse.Body)
			return httpResponse, err
		}
	}
	// As Content-Type is none of supported types body is interpreted as a string
	body, err := io.ReadAll(response.Body)
	httpResponse.Body = string(body)
	return httpResponse, err
}
