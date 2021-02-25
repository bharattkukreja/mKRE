/*
Copyright 2019 Google LLC

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

package v1beta1

import (
	"testing"

	gcpauthtesthelper "github.com/google/knative-gcp/pkg/apis/configs/gcpauth/testhelper"

	"github.com/google/go-cmp/cmp"
	duckv1beta1 "github.com/google/knative-gcp/pkg/apis/duck/v1beta1"
	schemasv1 "github.com/google/knative-gcp/pkg/schemas/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestCloudStorageSource_SetDefaults(t *testing.T) {
	testCases := map[string]struct {
		orig     *CloudStorageSourceSpec
		expected *CloudStorageSourceSpec
	}{
		"missing defaults": {
			orig: &CloudStorageSourceSpec{},
			expected: &CloudStorageSourceSpec{
				EventTypes: allEventTypes,
				PubSubSpec: duckv1beta1.PubSubSpec{
					Secret: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "google-cloud-key",
						},
						Key: "key.json",
					},
				},
			},
		},
		"defaults present": {
			orig: &CloudStorageSourceSpec{
				EventTypes: []string{schemasv1.CloudStorageObjectFinalizedEventType, schemasv1.CloudStorageObjectDeletedEventType},
				PubSubSpec: duckv1beta1.PubSubSpec{
					Secret: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "secret-name",
						},
						Key: "secret-key.json",
					},
				},
			},
			expected: &CloudStorageSourceSpec{
				EventTypes: []string{schemasv1.CloudStorageObjectFinalizedEventType, schemasv1.CloudStorageObjectDeletedEventType},
				PubSubSpec: duckv1beta1.PubSubSpec{
					Secret: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "secret-name",
						},
						Key: "secret-key.json",
					},
				},
			},
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			tc.orig.SetDefaults(gcpauthtesthelper.ContextWithDefaults())
			if diff := cmp.Diff(tc.expected, tc.orig); diff != "" {
				t.Errorf("Unexpected differences (-want +got): %v", diff)
			}
		})
	}
}
