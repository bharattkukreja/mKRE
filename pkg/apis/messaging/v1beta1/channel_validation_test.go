/*
Copyright 2019 The Knative Authors

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
	"context"
	"testing"

	pkgduckv1 "knative.dev/pkg/apis/duck/v1"

	"github.com/google/go-cmp/cmp"
	eventingduck "knative.dev/eventing/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/webhook/resourcesemantics"
)

var (
	backoffPolicy = eventingduck.BackoffPolicyExponential
	backoffDelay  = "PT1S"
)

func TestChannelValidation(t *testing.T) {
	tests := []struct {
		name string
		cr   resourcesemantics.GenericCRD
		want *apis.FieldError
	}{{
		name: "empty",
		cr: &Channel{
			Spec: ChannelSpec{},
		},
		want: nil,
	}, {
		name: "valid subscribers array",
		cr: &Channel{
			Spec: ChannelSpec{
				SubscribableSpec: &eventingduck.SubscribableSpec{
					Subscribers: []eventingduck.SubscriberSpec{{
						SubscriberURI: apis.HTTP("subscriberendpoint"),
						ReplyURI:      apis.HTTP("resultendpoint"),
					}},
				}},
		},
		want: nil,
	}, {
		name: "empty subscriber at index 1",
		cr: &Channel{
			Spec: ChannelSpec{
				SubscribableSpec: &eventingduck.SubscribableSpec{
					Subscribers: []eventingduck.SubscriberSpec{{
						SubscriberURI: apis.HTTP("subscriberendpoint"),
						ReplyURI:      apis.HTTP("replyendpoint"),
					}, {}},
				}},
		},
		want: func() *apis.FieldError {
			fe := apis.ErrMissingField("spec.subscribable.subscriber[1].replyURI", "spec.subscribable.subscriber[1].subscriberURI")
			fe.Details = "expected at least one of, got none"
			return fe
		}(),
	}, {
		name: "2 empty subscribers",
		cr: &Channel{
			Spec: ChannelSpec{
				SubscribableSpec: &eventingduck.SubscribableSpec{
					Subscribers: []eventingduck.SubscriberSpec{{}, {}},
				},
			},
		},
		want: func() *apis.FieldError {
			var errs *apis.FieldError
			fe := apis.ErrMissingField("spec.subscribable.subscriber[0].replyURI", "spec.subscribable.subscriber[0].subscriberURI")
			fe.Details = "expected at least one of, got none"
			errs = errs.Also(fe)
			fe = apis.ErrMissingField("spec.subscribable.subscriber[1].replyURI", "spec.subscribable.subscriber[1].subscriberURI")
			fe.Details = "expected at least one of, got none"
			errs = errs.Also(fe)
			return errs
		}(),
	}, {
		name: "invalid Delivery DeadLetterSink",
		cr: &Channel{
			Spec: ChannelSpec{
				SubscribableSpec: &eventingduck.SubscribableSpec{
					Subscribers: []eventingduck.SubscriberSpec{
						{
							ReplyURI: apis.HTTP("subscriberendpoint"),
							Delivery: &eventingduck.DeliverySpec{
								BackoffDelay:  &backoffDelay,
								BackoffPolicy: &backoffPolicy,
								DeadLetterSink: &pkgduckv1.Destination{
									URI: apis.HTTP("example.com"),
								},
							},
						},
					},
				},
			},
		},
		want: func() *apis.FieldError {
			var errs *apis.FieldError
			fe := apis.ErrInvalidValue("Dead letter sink URI scheme should be pubsub", "uri")
			errs = errs.Also(fe.ViaField("spec.delivery.subscriber[0].deadLetterSink"))
			return errs
		}(),
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.cr.Validate(context.TODO())
			if diff := cmp.Diff(test.want.Error(), got.Error()); diff != "" {
				t.Errorf("%s: validate (-want, +got) = %v", test.name, diff)
			}
		})
	}
}

func TestCheckImmutableFields(t *testing.T) {
	testCases := map[string]struct {
		orig    interface{}
		updated ChannelSpec
		allowed bool
	}{
		"nil orig": {
			updated: ChannelSpec{},
			allowed: true,
		},
	}

	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {
			var orig *Channel

			if tc.orig != nil {
				if spec, ok := tc.orig.(*ChannelSpec); ok {
					orig = &Channel{
						Spec: *spec,
					}
				}
			}
			updated := &Channel{
				Spec: tc.updated,
			}
			err := updated.CheckImmutableFields(context.TODO(), orig)
			if tc.allowed != (err == nil) {
				t.Fatalf("Unexpected immutable field check. Expected %v. Actual %v", tc.allowed, err)
			}
		})
	}
}
