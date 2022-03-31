package v1alpha1

// Parameter Definition e.g.
// name: url
// value: https://example.com
type Param struct {
	// Name of parameter
	Name string `json:"name"`

	// Value of parameter
	Value string `json:"value"`
}
