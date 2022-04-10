package util

import (
	"testing"

	universalapicontrolleriov1alpha1 "github.com/flyingdogfood/universal-api-controller/api/v1alpha1"
	"github.com/google/go-cmp/cmp"
)

func TestNewParameters(t *testing.T) {
	parameters := NewParameters()
	parameters2 := Parameters{
		Parameters: make(map[string]string),
		Responses:  make(map[string]interface{}),
	}
	if !cmp.Equal(parameters, parameters2) {
		t.Error("newParameters Function failed")
	}
}

func TestGenerateParameters(t *testing.T) {
	parameters := Parameters{
		Parameters: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
		Responses: make(map[string]interface{}),
	}
	params := []universalapicontrolleriov1alpha1.Param{
		{
			Name:  "param1",
			Value: "{{ .Parameters.a }}",
		},
		{
			Name:  "param2",
			Value: "{{ .Parameters.c }}{{ .Parameters.b }}",
		},
	}
	expectedResult := Parameters{
		Parameters: map[string]string{
			"param1": "a",
			"param2": "cb",
		},
		Responses: make(map[string]interface{}),
	}
	generateResult, _ := parameters.GenerateParameters(params)
	if !cmp.Equal(expectedResult, generateResult) {
		t.Errorf("Expected %v, got %v", expectedResult, generateResult)
	}
}

func TestMerge(t *testing.T) {
	parameters := Parameters{
		Parameters: make(map[string]string),
		Responses: map[string]interface{}{
			"resp1": HttpResponse{
				StatusCode: 200,
			},
			"resp2": HttpResponse{
				StatusCode: 404,
			},
		},
	}
	mergeParameters := Parameters{
		Parameters: make(map[string]string),
		Responses: map[string]interface{}{
			"resp3": HttpResponse{
				StatusCode: 200,
			},
			"resp4": HttpResponse{
				StatusCode: 404,
			},
		},
	}
	expectedParameters := Parameters{
		Parameters: make(map[string]string),
		Responses: map[string]interface{}{
			"resp1": HttpResponse{
				StatusCode: 200,
			},
			"resp2": HttpResponse{
				StatusCode: 404,
			},
			"merge": map[string]interface{}{
				"resp3": HttpResponse{
					StatusCode: 200,
				},
				"resp4": HttpResponse{
					StatusCode: 404,
				},
			},
		},
	}
	parameters.Merge(mergeParameters, "merge")
	if !cmp.Equal(expectedParameters, parameters) {
		t.Errorf("Expected %v, got %v", expectedParameters, parameters)
	}
}
