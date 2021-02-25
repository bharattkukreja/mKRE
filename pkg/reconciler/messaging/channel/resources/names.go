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

package resources

import (
	"github.com/google/knative-gcp/pkg/apis/messaging/v1beta1"
	"github.com/google/knative-gcp/pkg/utils/naming"
	"k8s.io/apimachinery/pkg/types"
)

// For reference, the minimum number of characters available for a name
// is 146. However, any name longer than 146 will be truncated and suffixed
// with a 32-char hash, making its max length 114 chars.
//
// pubsub resource name max length: 255 chars
// Namespace max length: 63 chars
// broker name max length: 253 chars
// trigger name max length: 253 chars
// uid length: 36 chars
// prefix + separators: 10 chars
// 255 - 10 - 63 - 36 = 146

// GenerateDecouplingTopicName generates a deterministic Topic name for a
// Channel. If the Topic name would be longer than allowed by PubSub, the
// Channel name is truncated to fit.
func GenerateDecouplingTopicName(c *v1beta1.Channel) string {
	return naming.TruncatedPubsubResourceName("cre-ch", c.Namespace, c.Name, c.UID)
}

// GenerateDecouplingSubscriptionName generates a deterministic Subscription
// name for a Channel. If the Subscription name would be longer than allowed by
// PubSub, the Channel name is truncated to fit.
func GenerateDecouplingSubscriptionName(c *v1beta1.Channel) string {
	return naming.TruncatedPubsubResourceName("cre-ch", c.Namespace, c.Name, c.UID)
}

// GenerateSubscriberRetryTopicName generates a deterministic Topic name for a Knative
// Subscription's retry topic. If the Topic name would be longer than allowed by PubSub, the name is
// truncated to fit.
func GenerateSubscriberRetryTopicName(c *v1beta1.Channel, subscriberUID types.UID) string {
	return naming.TruncatedPubsubResourceName("cre-sub", c.Namespace, c.Name, subscriberUID)
}

// GenerateSubscriberRetrySubscriptionName generates a deterministic Pub/Sub Subscription name for a
// Knative Subscription's retry Topic. If the subscription name would be longer than allowed by
// PubSub, the Subscription name is truncated to fit.
func GenerateSubscriberRetrySubscriptionName(c *v1beta1.Channel, subsciberUID types.UID) string {
	return naming.TruncatedPubsubResourceName("cre-sub", c.Namespace, c.Name, subsciberUID)
}
