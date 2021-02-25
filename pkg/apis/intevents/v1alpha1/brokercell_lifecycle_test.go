/*
Copyright 2020 Google LLC

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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

var (
	replicaUnavailableDeployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-deployment",
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{
					Type:    appsv1.DeploymentAvailable,
					Status:  corev1.ConditionFalse,
					Reason:  "MinimumReplicasUnavailable",
					Message: "False Status",
				},
			},
		},
	}
)

var (
	brokerCellConditionReady = apis.Condition{
		Type:   BrokerCellConditionReady,
		Status: corev1.ConditionTrue,
	}

	brokerCellConditionIngress = apis.Condition{
		Type:   BrokerCellConditionIngress,
		Status: corev1.ConditionTrue,
	}

	brokerCellConditionIngressFalse = apis.Condition{
		Type:   BrokerCellConditionIngress,
		Status: corev1.ConditionFalse,
	}

	brokerCellConditionFanout = apis.Condition{
		Type:   BrokerCellConditionFanout,
		Status: corev1.ConditionTrue,
	}

	brokerCellConditionRetry = apis.Condition{
		Type:   BrokerCellConditionRetry,
		Status: corev1.ConditionTrue,
	}
)

func TestBrokerCellGetCondition(t *testing.T) {
	tests := []struct {
		name      string
		ts        *BrokerCellStatus
		condQuery apis.ConditionType
		want      *apis.Condition
	}{{
		name: "single condition",
		ts: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					brokerCellConditionReady,
				},
			},
		},
		condQuery: apis.ConditionReady,
		want:      &brokerCellConditionReady,
	}, {
		name: "multiple conditions",
		ts: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					brokerCellConditionIngress,
					brokerCellConditionFanout,
				},
			},
		},
		condQuery: BrokerCellConditionIngress,
		want:      &brokerCellConditionIngress,
	}, {
		name: "multiple conditions, condition false",
		ts: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					brokerCellConditionIngressFalse,
					brokerCellConditionFanout,
				},
			},
		},
		condQuery: BrokerCellConditionIngress,
		want:      &brokerCellConditionIngressFalse,
	}, {
		name: "unknown condition",
		ts: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{
					brokerCellConditionIngress,
				},
			},
		},
		condQuery: apis.ConditionType("foo"),
		want:      nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.ts.GetCondition(test.condQuery)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("unexpected condition (-want, +got) = %v", diff)
			}
		})
	}
}

func TestBrokerCellInitializeConditions(t *testing.T) {
	tests := []struct {
		name string
		ts   *BrokerCellStatus
		want *BrokerCellStatus
	}{{
		name: "empty",
		ts:   &BrokerCellStatus{},
		want: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   BrokerCellConditionFanout,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionIngress,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionReady,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionRetry,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionTargetsConfig,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	}, {
		name: "one false",
		ts: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   BrokerCellConditionIngress,
					Status: corev1.ConditionFalse,
				}},
			},
		},
		want: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   BrokerCellConditionFanout,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionIngress,
					Status: corev1.ConditionFalse,
				}, {
					Type:   BrokerCellConditionReady,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionRetry,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionTargetsConfig,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	}, {
		name: "one true",
		ts: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   BrokerCellConditionIngress,
					Status: corev1.ConditionTrue,
				}},
			},
		},
		want: &BrokerCellStatus{
			Status: duckv1.Status{
				Conditions: []apis.Condition{{
					Type:   BrokerCellConditionFanout,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionIngress,
					Status: corev1.ConditionTrue,
				}, {
					Type:   BrokerCellConditionReady,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionRetry,
					Status: corev1.ConditionUnknown,
				}, {
					Type:   BrokerCellConditionTargetsConfig,
					Status: corev1.ConditionUnknown,
				}},
			},
		},
	}}

	ignoreAllButTypeAndStatus := cmpopts.IgnoreFields(
		apis.Condition{},
		"LastTransitionTime", "Message", "Reason", "Severity")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.ts.InitializeConditions()
			if diff := cmp.Diff(test.want, test.ts, ignoreAllButTypeAndStatus); diff != "" {
				t.Errorf("unexpected conditions (-want, +got) = %v", diff)
			}
		})
	}
}

func TestBrokerCellConditionStatus(t *testing.T) {
	tests := []struct {
		name                string
		fanoutStatus        *appsv1.Deployment
		ingressStatus       *corev1.Endpoints
		retryStatus         *appsv1.Deployment
		targetsStatus       bool
		wantConditionStatus corev1.ConditionStatus
	}{{
		name:                "all happy",
		fanoutStatus:        TestHelper.AvailableDeployment(),
		ingressStatus:       TestHelper.AvailableEndpoints(),
		retryStatus:         TestHelper.AvailableDeployment(),
		targetsStatus:       true,
		wantConditionStatus: corev1.ConditionTrue,
	}, {
		name:                "fanout sad",
		fanoutStatus:        TestHelper.UnavailableDeployment(),
		ingressStatus:       TestHelper.AvailableEndpoints(),
		retryStatus:         TestHelper.AvailableDeployment(),
		targetsStatus:       true,
		wantConditionStatus: corev1.ConditionFalse,
	}, {
		name:                "fanout unknown",
		fanoutStatus:        TestHelper.UnknownDeployment(),
		ingressStatus:       TestHelper.AvailableEndpoints(),
		retryStatus:         TestHelper.AvailableDeployment(),
		targetsStatus:       true,
		wantConditionStatus: corev1.ConditionUnknown,
	}, {
		name:                "ingress sad",
		fanoutStatus:        TestHelper.AvailableDeployment(),
		ingressStatus:       TestHelper.UnavailableEndpoints(),
		retryStatus:         TestHelper.AvailableDeployment(),
		targetsStatus:       true,
		wantConditionStatus: corev1.ConditionFalse,
	}, {
		name:                "retry sad",
		fanoutStatus:        TestHelper.AvailableDeployment(),
		ingressStatus:       TestHelper.AvailableEndpoints(),
		retryStatus:         TestHelper.UnavailableDeployment(),
		targetsStatus:       true,
		wantConditionStatus: corev1.ConditionFalse,
	}, {
		name:                "retry unknown",
		fanoutStatus:        TestHelper.AvailableDeployment(),
		ingressStatus:       TestHelper.AvailableEndpoints(),
		retryStatus:         TestHelper.UnknownDeployment(),
		targetsStatus:       true,
		wantConditionStatus: corev1.ConditionUnknown,
	}, {
		name:                "targets sad",
		fanoutStatus:        TestHelper.AvailableDeployment(),
		ingressStatus:       TestHelper.AvailableEndpoints(),
		retryStatus:         TestHelper.AvailableDeployment(),
		targetsStatus:       false,
		wantConditionStatus: corev1.ConditionFalse,
	}, {
		name:                "all sad",
		fanoutStatus:        TestHelper.UnavailableDeployment(),
		ingressStatus:       TestHelper.UnavailableEndpoints(),
		retryStatus:         TestHelper.UnavailableDeployment(),
		targetsStatus:       false,
		wantConditionStatus: corev1.ConditionFalse,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bs := &BrokerCellStatus{}
			if test.fanoutStatus != nil {
				bs.PropagateFanoutAvailability(test.fanoutStatus)
			} else {
				bs.PropagateFanoutAvailability(&appsv1.Deployment{})
			}
			if test.ingressStatus != nil {
				bs.PropagateIngressAvailability(test.ingressStatus, &appsv1.Deployment{})
			} else {
				bs.PropagateIngressAvailability(&corev1.Endpoints{}, &appsv1.Deployment{})
			}
			if test.retryStatus != nil {
				bs.PropagateRetryAvailability(test.retryStatus)
			} else {
				bs.PropagateRetryAvailability(&appsv1.Deployment{})
			}
			if test.targetsStatus {
				bs.MarkTargetsConfigReady()
			} else {
				bs.MarkTargetsConfigFailed("Unable to sync targets config", "induced failure")
			}
			got := bs.GetTopLevelCondition().Status
			if test.wantConditionStatus != got {
				t.Errorf("unexpected readiness: want %v, got %v", test.wantConditionStatus, got)
			}
			happy := bs.IsReady()
			switch test.wantConditionStatus {
			case corev1.ConditionTrue:
				if !happy {
					t.Error("expected happy true, got false")
				}
			case corev1.ConditionFalse, corev1.ConditionUnknown:
				if happy {
					t.Error("expected happy false, got true")
				}
			}
		})
	}
}

func TestMarkBrokerCellStatus(t *testing.T) {
	tests := []struct {
		name          string
		s             *BrokerCellStatus
		wantType      apis.ConditionType
		wantCondition corev1.ConditionStatus
	}{{
		name: "mark IngressReady unknown",
		s: func() *BrokerCellStatus {
			s := &BrokerCell{}
			s.Status.InitializeConditions()
			s.Status.MarkIngressUnknown("test", "the status of ingressReady is unknown")
			return &s.Status
		}(),
		wantType:      BrokerCellConditionIngress,
		wantCondition: corev1.ConditionUnknown,
	}, {}, {
		name: "mark IngressReady false",
		s: func() *BrokerCellStatus {
			s := &BrokerCell{}
			s.Status.InitializeConditions()
			s.Status.MarkIngressFailed("test", "the status of ingressReady is false")
			return &s.Status
		}(),
		wantType:      BrokerCellConditionIngress,
		wantCondition: corev1.ConditionFalse,
	}, {
		name: "mark FanoutReady unknown",
		s: func() *BrokerCellStatus {
			s := &BrokerCell{}
			s.Status.InitializeConditions()
			s.Status.MarkFanoutUnknown("test", "the status of fanoutReady is unknown")
			return &s.Status
		}(),
		wantType:      BrokerCellConditionFanout,
		wantCondition: corev1.ConditionUnknown,
	}, {
		name: "mark FanoutReady false",
		s: func() *BrokerCellStatus {
			s := &BrokerCell{}
			s.Status.InitializeConditions()
			s.Status.MarkFanoutFailed("test", "the status of fanoutReady is false")
			return &s.Status
		}(),
		wantType:      BrokerCellConditionFanout,
		wantCondition: corev1.ConditionFalse,
	}, {
		name: "mark RetryReady unknown",
		s: func() *BrokerCellStatus {
			s := &BrokerCell{}
			s.Status.InitializeConditions()
			s.Status.MarkRetryUnknown("test", "the status of retryReady is unknown")
			return &s.Status
		}(),
		wantType:      BrokerCellConditionRetry,
		wantCondition: corev1.ConditionUnknown,
	}, {
		name: "mark RetryReady false",
		s: func() *BrokerCellStatus {
			s := &BrokerCell{}
			s.Status.InitializeConditions()
			s.Status.MarkRetryFailed("test", "the status of retryReady is false")
			return &s.Status
		}(),
		wantType:      BrokerCellConditionRetry,
		wantCondition: corev1.ConditionFalse,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantType != "" {
				gotConditionStatus := test.s.Status
				for _, cd := range gotConditionStatus.Conditions {
					if cd.Type == test.wantType {
						if cd.Status == test.wantCondition {
							return
						}
						t.Errorf("unexpected condition status for %v: want %v, got %v", test.wantType, test.wantCondition, cd.Status)
					}
				}
				t.Error("didn't see the expected condition: ", test.wantType)
			}
		})
	}
}

func TestPropagateDeploymentAvailability(t *testing.T) {
	t.Run("propagate ingress availability", func(t *testing.T) {
		s := &BrokerCellStatus{}
		got := s.PropagateIngressAvailability(&corev1.Endpoints{}, replicaUnavailableDeployment)
		want := false
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error("unexpected condition (-want, +got) =", diff)
		}
	})

	t.Run("propagate fanout availability", func(t *testing.T) {
		s := &BrokerCellStatus{}
		got := s.PropagateFanoutAvailability(replicaUnavailableDeployment)
		want := false
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error("unexpected condition (-want, +got) =", diff)
		}
	})

	t.Run("propagate retry availability", func(t *testing.T) {
		s := &BrokerCellStatus{}
		got := s.PropagateRetryAvailability(replicaUnavailableDeployment)
		want := false
		if diff := cmp.Diff(want, got); diff != "" {
			t.Error("unexpected condition (-want, +got) =", diff)
		}
	})
}
