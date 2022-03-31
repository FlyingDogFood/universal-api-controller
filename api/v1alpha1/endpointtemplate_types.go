/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EndpointTemplateSpec defines the desired state of EndpointTemplate
type EndpointTemplateSpec struct {
	//Defines the parameters for the template
	Params []ParamDefinition `json:"params,omitempty"`

	// Defines the http method for this endpoint
	Method string `json:"method"`

	// Defines the url to call
	URL string `json:"url"`

	// Defines the http headers for the Request
	Headers map[string]string `json:"headers,omitempty"`

	// Defines the http body
	Body string `json:"body,omitempty"`
}

//+kubebuilder:object:root=true

// EndpointTemplate is the Schema for the endpointtemplates API
type EndpointTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EndpointTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// EndpointTemplateList contains a list of EndpointTemplate
type EndpointTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EndpointTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EndpointTemplate{}, &EndpointTemplateList{})
}
