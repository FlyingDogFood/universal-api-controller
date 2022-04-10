package util

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status map[string]interface{}

func NewStatus() Status {
	status := make(Status)
	status["reconciledAt"] = metav1.Now()
	return status
}

func (s *Status) Bytes() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Status) Merge(status Status, name string) {
	(*s)[name] = status
}
