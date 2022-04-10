package util

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status map[string]interface{}

func newStatus() Status {
	status := make(Status)
	status["reconciledAt"] = metav1.Now()
	return status
}

func (s *Status) bytes() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Status) merge(status Status, name string) {
	(*s)[name] = status
}
