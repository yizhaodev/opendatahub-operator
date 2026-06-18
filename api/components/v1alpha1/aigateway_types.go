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

package v1alpha1

import (
	"github.com/opendatahub-io/opendatahub-operator/v2/api/common"
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AIGatewayComponentName = "aigateway"
	// value should match whats set in the XValidation below
	AIGatewayInstanceName = "default-" + AIGatewayComponentName
	AIGatewayKind         = "AIGateway"
)

// Check that the component implements common.PlatformObject.
var _ common.PlatformObject = (*AIGateway)(nil)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'default-aigateway'",message="AIGateway name must be default-aigateway"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="Ready"
// +kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].reason`,description="Reason"

// AIGateway is the Schema for the aigateways API
type AIGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AIGatewaySpec   `json:"spec,omitempty"`
	Status AIGatewayStatus `json:"status,omitempty"`
}

// AIGatewaySpec defines the desired state of AIGateway
type AIGatewaySpec struct {
	AIGatewayCommonSpec `json:",inline"`
}

type AIGatewayCommonSpec struct {
	// BatchGateway sub-component configuration.
	BatchGateway AIGatewayBatchGatewaySpec `json:"batchGateway,omitempty"`
}

// AIGatewayBatchGatewaySpec defines the configuration for the BatchGateway sub-component.
type AIGatewayBatchGatewaySpec struct {
	// Set to one of the following values:
	// - "Managed" : the operator is actively managing the component and trying to keep it active.
	//               It will only upgrade the component if it is safe to do so.
	// - "Removed" : the operator is actively managing the component and will not install it,
	//               or if it is installed, the operator will try to remove it.
	//
	// +kubebuilder:validation:Enum=Managed;Removed
	ManagementState operatorv1.ManagementState `json:"managementState,omitempty"`
}

// AIGatewayCommonStatus defines the shared observed state of AIGateway
type AIGatewayCommonStatus struct {
	common.ComponentReleaseStatus `json:",inline"`
}

// AIGatewayStatus defines the observed state of AIGateway
type AIGatewayStatus struct {
	common.Status         `json:",inline"`
	AIGatewayCommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// AIGatewayList contains a list of AIGateway
type AIGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AIGateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AIGateway{}, &AIGatewayList{})
}

func (c *AIGateway) GetStatus() *common.Status {
	return &c.Status.Status
}

func (c *AIGateway) GetConditions() []common.Condition {
	return c.Status.GetConditions()
}

func (c *AIGateway) SetConditions(conditions []common.Condition) {
	c.Status.SetConditions(conditions)
}

func (c *AIGateway) GetReleaseStatus() *[]common.ComponentRelease { return &c.Status.Releases }

func (c *AIGateway) SetReleaseStatus(releases []common.ComponentRelease) {
	c.Status.Releases = releases
}

// DSCAIGateway contains all the configuration exposed in DSC instance for AIGateway component
type DSCAIGateway struct {
	common.ManagementSpec `json:",inline"`
	// configuration fields common across components
	AIGatewayCommonSpec `json:",inline"`
}

// DSCAIGatewayStatus struct holds the status for the AIGateway component exposed in the DSC
type DSCAIGatewayStatus struct {
	common.ManagementSpec  `json:",inline"`
	*AIGatewayCommonStatus `json:",inline"`
}
