package v1alpha1

type FunctionOrEndpointTemplateRef struct {
	// name of action
	Name string `json:"name"`

	// Reference to Function or Endpoint
	Ref ObjectRef `json:"ref"`

	// Parameters for execution
	Params []Param `json:"params,omitempty"`
}

type ObjectRef struct {
	// Type of Object: Has to be Function or EndpointTemplate
	Type string `json:"type"`

	// Name of the Function or EndpointTemplate
	Name string `json:"name"`
}
