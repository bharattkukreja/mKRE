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

package filter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/extensions"
	"github.com/google/go-cmp/cmp"
	"go.opencensus.io/trace"

	"github.com/google/knative-gcp/pkg/broker/config"
	"github.com/google/knative-gcp/pkg/broker/config/memory"
	handlerctx "github.com/google/knative-gcp/pkg/broker/handler/context"
	"github.com/google/knative-gcp/pkg/broker/handler/processors"
)

const (
	traceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanID  = "00f067aa0ba902b7"
)

func TestInvalidContext(t *testing.T) {
	p := &Processor{}
	e := event.New()
	err := p.Process(context.Background(), &e)
	if err != handlerctx.ErrTargetKeyNotPresent {
		t.Errorf("Process error got=%v, want=%v", err, handlerctx.ErrTargetKeyNotPresent)
	}
}

type VerifyTraceID struct {
	processors.BaseProcessor
	wantTraceID string
}

func (p *VerifyTraceID) Process(ctx context.Context, event *event.Event) error {
	if got := trace.FromContext(ctx).SpanContext().TraceID.String(); got != p.wantTraceID {
		return fmt.Errorf("unexpected trace ID: got %s, want %s", got, p.wantTraceID)
	}
	return nil
}

func TestTriggerSpanCreatedFromTraceparent(t *testing.T) {
	e := event.New()
	e.SetID("id")
	e.SetSubject("foo")
	e.SetType("bar")
	extensions.DistributedTracingExtension{
		TraceParent: fmt.Sprintf("00-%s-%s-01", traceID, spanID),
	}.AddTracingAttributes(&e)

	ctx, testTargets := newTestTargets(nil)

	p := &Processor{Targets: testTargets}
	p.WithNext(&VerifyTraceID{wantTraceID: traceID})

	if err := p.Process(ctx, &e); err != nil {
		t.Error(err)
	}
}

func TestFilterProcessor(t *testing.T) {
	tn := time.Now()
	cases := []struct {
		name       string
		e          event.Event
		filter     map[string]string
		shouldPass bool
	}{{
		name: "no filter pass",
		e: func() event.Event {
			e := event.New()
			e.SetID("id")
			e.SetSubject("foo")
			e.SetType("bar")
			return e
		}(),
		shouldPass: true,
	}, {
		name: "match spec version pass",
		e: func() event.Event {
			e := event.New(event.CloudEventsVersionV1)
			return e
		}(),
		filter: map[string]string{
			"specversion": event.CloudEventsVersionV1,
		},
		shouldPass: true,
	}, {
		name: "match spec version not pass",
		e: func() event.Event {
			e := event.New(event.CloudEventsVersionV1)
			return e
		}(),
		filter: map[string]string{
			"specversion": event.CloudEventsVersionV03,
		},
		shouldPass: false,
	}, {
		name: "match type pass",
		e: func() event.Event {
			e := event.New()
			e.SetType("foo")
			return e
		}(),
		filter: map[string]string{
			"type": "foo",
		},
		shouldPass: true,
	}, {
		name: "match type not pass",
		e: func() event.Event {
			e := event.New()
			e.SetType("foo")
			return e
		}(),
		filter: map[string]string{
			"type": "bar",
		},
		shouldPass: false,
	}, {
		name: "match source pass",
		e: func() event.Event {
			e := event.New()
			e.SetSource("foo")
			return e
		}(),
		filter: map[string]string{
			"source": "foo",
		},
		shouldPass: true,
	}, {
		name: "match source not pass",
		e: func() event.Event {
			e := event.New()
			e.SetSource("foo")
			return e
		}(),
		filter: map[string]string{
			"source": "bar",
		},
		shouldPass: false,
	}, {
		name: "match subject pass",
		e: func() event.Event {
			e := event.New()
			e.SetSubject("foo")
			return e
		}(),
		filter: map[string]string{
			"subject": "foo",
		},
		shouldPass: true,
	}, {
		name: "match subject not pass",
		e: func() event.Event {
			e := event.New()
			e.SetSubject("foo")
			return e
		}(),
		filter: map[string]string{
			"subject": "bar",
		},
		shouldPass: false,
	}, {
		name: "match time pass",
		e: func() event.Event {
			e := event.New()
			e.SetTime(tn)
			return e
		}(),
		filter: map[string]string{
			"time": tn.String(),
		},
		shouldPass: true,
	}, {
		name: "match time not pass",
		e: func() event.Event {
			e := event.New()
			e.SetTime(tn)
			return e
		}(),
		filter: map[string]string{
			"time": time.Now().Add(time.Hour).String(),
		},
		shouldPass: false,
	}, {
		name: "match schemaurl pass",
		e: func() event.Event {
			e := event.New()
			e.SetDataSchema("foo")
			return e
		}(),
		filter: map[string]string{
			"schemaurl": "foo",
		},
		shouldPass: true,
	}, {
		name: "match schemaurl not pass",
		e: func() event.Event {
			e := event.New()
			e.SetDataSchema("foo")
			return e
		}(),
		filter: map[string]string{
			"schemaurl": "bar",
		},
		shouldPass: false,
	}, {
		name: "match datacontenttype pass",
		e: func() event.Event {
			e := event.New()
			e.SetDataContentType("foo")
			return e
		}(),
		filter: map[string]string{
			"datacontenttype": "foo",
		},
		shouldPass: true,
	}, {
		name: "match datacontenttype not pass",
		e: func() event.Event {
			e := event.New()
			e.SetDataContentType("foo")
			return e
		}(),
		filter: map[string]string{
			"datacontenttype": "bar",
		},
		shouldPass: false,
	}, {
		name: "mixed fileter pass",
		e: func() event.Event {
			e := event.New()
			e.SetID("id")
			e.SetSource("foo")
			e.SetType("bar")
			e.SetSubject("subject")
			return e
		}(),
		filter: map[string]string{
			"id":      "id",
			"subject": "subject",
		},
		shouldPass: true,
	}, {
		name: "mixed fileter not pass",
		e: func() event.Event {
			e := event.New()
			e.SetID("id")
			e.SetSource("foo")
			e.SetType("bar")
			e.SetSubject("subject")
			return e
		}(),
		filter: map[string]string{
			"id":      "id",
			"subject": "subject",
			"source":  "unknown",
		},
		shouldPass: false,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, testTargets := newTestTargets(tc.filter)
			next := &processors.FakeProcessor{}
			p := &Processor{Targets: testTargets}
			p.WithNext(next)
			ch := make(chan *event.Event, 1)
			next.PrevEventsCh = ch
			defer func() {
				gotEvent := <-ch
				if tc.shouldPass {
					if diff := cmp.Diff(&tc.e, gotEvent); diff != "" {
						t.Errorf("processed event (-want,+got): %v", diff)
					}
				} else {
					if gotEvent != nil {
						t.Errorf("unexpected event %v passed filter %v", gotEvent, tc.filter)
					}
				}
			}()

			if err := p.Process(ctx, &tc.e); err != nil {
				t.Errorf("unexpected error from processing: %v", err)
			}
			// In case the event doesn't pass the filter,
			// we need to close the channel to make sure defer func returns.
			close(ch)
		})
	}
}

func newTestTargets(filter map[string]string) (context.Context, config.Targets) {
	testTarget := &config.Target{
		Name:             "target",
		CellTenantType:   config.CellTenantType_BROKER,
		CellTenantName:   "broker",
		Namespace:        "ns",
		FilterAttributes: filter,
	}
	testTargets := memory.NewEmptyTargets()
	testTargets.MutateCellTenant(testTarget.Key().ParentKey(), func(bm config.CellTenantMutation) {
		bm.UpsertTargets(testTarget)
	})
	ctx := handlerctx.WithTargetKey(context.Background(), testTarget.Key())
	return ctx, testTargets
}
