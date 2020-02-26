package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LightSwitchSpec defines the desired state of LightSwitch
type LightSwitchSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	serviceName string `json:"name"`
	spotDisable bool   `json:"spotDisable"`
	logzio      bool   `json:"logzio"`
	port        int32  `json:"port"`
	//replicaCount int32  `json:"replicaCount"`
}

// LightSwitchStatus defines the observed state of LightSwitch
type LightSwitchStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LightSwitch is the Schema for the lightswitches API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=lightswitches,scope=Namespaced
type LightSwitch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LightSwitchSpec   `json:"spec,omitempty"`
	Status LightSwitchStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LightSwitchList contains a list of LightSwitch
type LightSwitchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LightSwitch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LightSwitch{}, &LightSwitchList{})
}
