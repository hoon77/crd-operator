/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebAppSpec defines the desired state of WebApp.
type WebAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of WebApp. Edit webapp_types.go to remove/update
	Image      string            `json:"image"`
	Replicas   *int32            `json:"replicas"`
	ConfigData map[string]string `json:"configData,omitempty"`
	Ingress    *IngressSpec      `json:"ingress,omitempty"`
}

// WebAppStatus defines the observed state of WebApp.
type WebAppStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WebApp is the Schema for the webapps API.
type WebApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebAppSpec   `json:"spec,omitempty"`
	Status WebAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebAppList contains a list of WebApp.
type WebAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebApp `json:"items"`
}

type IngressSpec struct {
	Enabled       bool   `json:"enabled"`
	ClassName     string `json:"className,omitempty"`
	Host          string `json:"host,omitempty"`
	Path          string `json:"path,omitempty"`
	Port          int32  `json:"port,omitempty"`
	RewriteTarget string `json:"rewriteTarget,omitempty"`
	TLS           bool   `json:"tls,omitempty"`
}

func init() {
	SchemeBuilder.Register(&WebApp{}, &WebAppList{})
}
