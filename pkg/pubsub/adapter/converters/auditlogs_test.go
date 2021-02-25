/*
Copyright 2019 Google LLC.

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
	"bytes"
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	schemasv1 "github.com/google/knative-gcp/pkg/schemas/v1"
	auditpb "google.golang.org/genproto/googleapis/cloud/audit"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/protobuf/testing/protocmp"
)

const (
	insertID = "test-insert-id"
	logName  = "projects/test-project/pubsub.googleapis.com%2Factivity"
	testTs   = "2006-01-02T15:04:05Z"
)

func TestConvertAuditLog(t *testing.T) {
	auditLog := auditpb.AuditLog{
		ServiceName:  "pubsub.googleapis.com",
		MethodName:   "test-method-name",
		ResourceName: "projects/test-project/topics/test-topic",
	}
	payload, err := ptypes.MarshalAny(&auditLog)
	if err != nil {
		t.Fatalf("Failed to marshal proto payload: %v", err)
	}
	logEntry := logpb.LogEntry{
		InsertId: insertID,
		LogName:  logName,
		Timestamp: &timestamp.Timestamp{
			Seconds: 12345,
		},
		Payload: &logpb.LogEntry_ProtoPayload{
			ProtoPayload: payload,
		},
	}
	testTime, err := time.Parse(time.RFC3339, testTs)
	if err != nil {
		t.Fatalf("Unable to parse test timestamp: %q", err)
	}
	if ts, err := ptypes.TimestampProto(testTime); err != nil {
		t.Fatalf("Invalid test timestamp: %q", err)
	} else {
		logEntry.Timestamp = ts
	}
	var buf bytes.Buffer
	if err := new(jsonpb.Marshaler).Marshal(&buf, &logEntry); err != nil {
		t.Fatalf("Failed to marshal AuditLog pb: %v", err)
	}
	msg := pubsub.Message{
		Data: buf.Bytes(),
	}

	e, err := NewPubSubConverter().Convert(context.Background(), &msg, CloudAuditLogs)

	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}
	if id := schemasv1.CloudAuditLogsEventID(insertID, logName, testTs); e.ID() != id {
		t.Errorf("ID '%s' != '%s'", e.ID(), id)
	}
	if !e.Time().Equal(testTime) {
		t.Errorf("Time '%v' != '%v'", e.Time(), testTime)
	}
	if want := schemasv1.CloudAuditLogsEventSource("projects/test-project", "activity"); e.Source() != want {
		t.Errorf("Source %q != %q", e.Source(), want)
	}
	if e.Type() != "google.cloud.audit.log.v1.written" {
		t.Errorf(`Type %q != "google.cloud.audit.log.v1.written"`, e.Type())
	}
	if want := schemasv1.CloudAuditLogsEventSubject("pubsub.googleapis.com", "projects/test-project/topics/test-topic"); e.Subject() != want {
		t.Errorf("Subject %q != %q", e.Subject(), want)
	}
	if e.DataSchema() != schemasv1.CloudAuditLogsEventDataSchema {
		t.Errorf("DataSchema got=%s, want=%s", e.DataSchema(), schemasv1.CloudAuditLogsEventDataSchema)
	}

	var actualLogEntry logpb.LogEntry
	if err = jsonpb.Unmarshal(bytes.NewReader(e.Data()), &actualLogEntry); err != nil {
		t.Errorf("Unable to unmarshal event data to LogEntry: %q", err)
	} else {
		if diff := cmp.Diff(logEntry, actualLogEntry, protocmp.Transform()); diff != "" {
			t.Errorf("unexpected LogEntry (-want, +got) = %v", diff)
		}
	}

	wantExtensions := map[string]interface{}{
		"servicename":  "pubsub.googleapis.com",
		"methodname":   "test-method-name",
		"resourcename": "projects/test-project/topics/test-topic",
	}
	if diff := cmp.Diff(wantExtensions, e.Extensions()); diff != "" {
		t.Errorf("unexpected (-want, +got) = %v", diff)
	}
}

func TestConvertTextPayload(t *testing.T) {
	logEntry := logpb.LogEntry{
		InsertId: insertID,
		LogName:  logName,
		Timestamp: &timestamp.Timestamp{
			Seconds: 12345,
		},
		Payload: &logpb.LogEntry_TextPayload{
			TextPayload: "test payload",
		},
	}
	testTime, err := time.Parse(time.RFC3339, testTs)
	if err != nil {
		t.Fatalf("Unable to parse test timestamp: %q", err)
	}
	if ts, err := ptypes.TimestampProto(testTime); err != nil {
		t.Fatalf("Invalid test timestamp: %q", err)
	} else {
		logEntry.Timestamp = ts
	}
	var buf bytes.Buffer
	if err := new(jsonpb.Marshaler).Marshal(&buf, &logEntry); err != nil {
		t.Fatalf("Failed to marshal AuditLog pb: %v", err)
	}
	msg := pubsub.Message{
		Data: buf.Bytes(),
	}

	_, err = NewPubSubConverter().Convert(context.Background(), &msg, CloudAuditLogs)

	if err == nil {
		t.Errorf("Expected error when converting non-AuditLog LogEntry.")
	}
}
