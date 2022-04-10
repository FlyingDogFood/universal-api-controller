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

// ConfigTemplateSpec defines the desired state of ConfigTemplate
type ConfigTemplateSpec struct {
	// Defines Parameters for ConfigTemplate
	Params []ParamDefinition `json:"params,omitempty"`

	// Defines Ref for reconciliation
	Reconcile []FunctionOrEndpointTemplateRef `json:"reconcile"`

	// Defines Ref for deletion
	Delete []FunctionOrEndpointTemplateRef `json:"delete"`
}

//+kubebuilder:object:root=true

// ConfigTemplate is the Schema for the configtemplates API
type ConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConfigTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigTemplateList contains a list of ConfigTemplate
type ConfigTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigTemplate{}, &ConfigTemplateList{})
}
