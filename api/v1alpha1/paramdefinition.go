package v1alpha1

type ParamDefinition struct {
	Name string `json:"name"`

	Description string `json:"description,omitempty"`
}