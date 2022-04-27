/*
Copyright 2022 Amazon.com, Inc. or its affiliates. All rights reserved.

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
	"context"
	"fmt"
	"net/url"

	"k8s.io/apimachinery/pkg/util/sets"

	"knative.dev/pkg/apis"
)

var knownActions = sets.NewString(
	Actions.CordonAndDrain,
	Actions.Cordon,
	Actions.NoAction,
)

func (t *Terminator) Validate(_ context.Context) (errs *apis.FieldError) {
	return errs.Also(
		apis.ValidateObjectMetadata(t).ViaField("metadata"),
		t.Spec.validate().ViaField("spec"),
	)
}

func (t *TerminatorSpec) validate() (errs *apis.FieldError) {
	return errs.Also(
		t.validateMatchLabels().ViaField("matchLabels"),
		t.SQS.validate().ViaField("sqs"),
		t.Events.validate().ViaField("events"),
	)
}

func (t *TerminatorSpec) validateMatchLabels() (errs *apis.FieldError) {
	for name, value := range t.MatchLabels {
		if value == "" {
			errs = errs.Also(apis.ErrInvalidValue(value, name, "label value cannot be empty"))
		}
	}
	return errs
}

func (s *SQSSpec) validate() (errs *apis.FieldError) {
	if _, err := url.Parse(s.QueueURL); err != nil {
		errs = errs.Also(apis.ErrInvalidValue(s.QueueURL, "queueURL", "must be a valid URL"))
	}
	return errs
}

func (e *EventsSpec) validate() (errs *apis.FieldError) {
	errMsg := fmt.Sprintf("must be one of %s", knownActions.List())
	for name, value := range map[string]string{
		"autoScalingTermination":  e.AutoScalingTermination,
		"rebalanceRecommendation": e.RebalanceRecommendation,
		"scheduledChange":         e.ScheduledChange,
		"spotInterruption":        e.SpotInterruption,
		"stateChange":             e.StateChange,
	} {
		if !knownActions.Has(value) {
			errs = errs.Also(apis.ErrInvalidValue(value, name, errMsg))
		}
	}
	return errs
}
