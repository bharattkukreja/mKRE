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

package context

import (
	"context"

	"github.com/google/knative-gcp/pkg/broker/config"
)

// The key used to store/retrieve broker in the context.
type brokerKey struct{}

// WithBrokerKey sets a broker key in the context.
func WithBrokerKey(ctx context.Context, key *config.CellTenantKey) context.Context {
	return context.WithValue(ctx, brokerKey{}, key)
}

// GetBrokerKey gets the broker key from the context.
func GetBrokerKey(ctx context.Context) (*config.CellTenantKey, error) {
	untyped := ctx.Value(brokerKey{})
	if untyped == nil {
		return nil, ErrBrokerKeyNotPresent
	}
	return untyped.(*config.CellTenantKey), nil
}
