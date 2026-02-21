package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// WebPage is the Schema for the webpages API
// +kubebuilder:object:root=true
type WebPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WebPageSpec `json:"spec"`
}

// WebPageList contains a list of WebPage
// +kubebuilder:object:root=true
type WebPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebPage `json:"items"`
}

// WebPageSpec defines the desired state of WebPage
type WebPageSpec struct {
	Content string `json:"content"`
	Image   string `json:"image"`
}

func init() {
	SchemeBuilder.Register(&WebPage{}, &WebPageList{})
}
