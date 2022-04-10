package util

import (
	universalapicontrolleriov1alpha1 "github.com/flyingdogfood/universal-api-controller/api/v1alpha1"
)

type Parameters struct {
	Parameters map[string]string      `json:"Parameters"`
	Responses  map[string]interface{} `json:"Responses"`
}

func newParameters() Parameters {
	return Parameters{
		Parameters: make(map[string]string),
		Responses:  make(map[string]interface{}),
	}
}

func (p *Parameters) generateParameters(parameters []universalapicontrolleriov1alpha1.Param) (Parameters, error) {
	params := newParameters()
	var err error
	for _, parameter := range parameters {
		params.Parameters[parameter.Name], err = templateString(parameter.Value, *p)
		if err != nil {
			return params, err
		}
	}
	return params, nil
}

func (p *Parameters) merge(parameters Parameters, name string) {
	(*p).Responses[name] = make(map[string]interface{})
	(*p).Responses[name] = parameters.Responses
}
