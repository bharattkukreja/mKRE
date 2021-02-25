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

package converters

import (
	"context"
	"errors"

	"github.com/google/knative-gcp/pkg/apis/events/v1beta1"
	schemasv1 "github.com/google/knative-gcp/pkg/schemas/v1"

	"cloud.google.com/go/pubsub"
	cev2 "github.com/cloudevents/sdk-go/v2"
)

func convertCloudScheduler(ctx context.Context, msg *pubsub.Message) (*cev2.Event, error) {
	event := cev2.NewEvent(cev2.VersionV1)
	event.SetID(msg.ID)
	event.SetTime(msg.PublishTime)
	event.SetType(schemasv1.CloudSchedulerJobExecutedEventType)
	event.SetDataSchema(schemasv1.CloudSchedulerEventDataSchema)

	jobName, ok := msg.Attributes[v1beta1.CloudSchedulerSourceJobName]
	if !ok {
		return nil, errors.New("received event did not have jobName")
	}
	event.SetSource(schemasv1.CloudSchedulerEventSource(jobName))

	if err := event.SetData(cev2.ApplicationJSON, &schemasv1.SchedulerJobData{CustomData: msg.Data}); err != nil {
		return nil, err
	}
	return &event, nil
}
