/*
Copyright 2021 James.Liu.

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

// RedisSentinelSpec defines the desired state of RedisSentinel
type RedisSentinelSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of RedisSentinel. Edit redissentinel_types.go to remove/update
	Redis    RedisService `json:"redis"`
	Sentinel Sentinel     `json:"sentinel,omitempty"`
}

type Common struct {
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
	Port     int32  `json:"port"`
}

type RedisService struct {
	Common `json:",inline"`
	PassWd string `json:"passwd"`
}

type Sentinel struct {
	Common `json:",inline"`
}

// RedisSentinelStatus defines the observed state of RedisSentinel
type RedisSentinelStatus struct {
	Redis_HEALTH    string `json:"redis_health,omitempty"`    // 集群健康状态
	Sentinel_HEALTH string `json:"sentinel_health,omitempty"` // 集群健康状态
	Master          string `json:"master,omitempty"`          // 那个是主
	Message         string `json:"message,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RedisSentinel is the Schema for the redissentinels API
type RedisSentinel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSentinelSpec   `json:"spec,omitempty"`
	Status RedisSentinelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisSentinelList contains a list of RedisSentinel
type RedisSentinelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisSentinel `json:"items"`
}

// redisSentinelCommonLabel 返回节点 label 选择器
// func (redisSentinel *RedisSentinel) CommonLabelSelector() labels.Selector {
// 	return labels.SelectorFromSet(map[string]string{
// 		"app": "redisSentinel",
// 	})
// }

func init() {
	SchemeBuilder.Register(&RedisSentinel{}, &RedisSentinelList{})
}

func (r *RedisSentinel) GetRedisSentinelName(typeName string) string {
	return r.Name + "-" + typeName
}

// func (s *RedisSentinelSpec) Apply(svc corev1.Service) *corev1.Service {
// 	svc = corev1.Service{}
// 	return &svc
// }

// func (r *RedisSentinel) RuntimeClass() *v1.RuntimeClass {
// 	// s := r.Spec
// 	return &v1.RuntimeClass{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "",
// 			Namespace: "",
// 		},
// 		Handler: "runc",
// 	}
// }
