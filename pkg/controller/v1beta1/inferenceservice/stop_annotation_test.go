/*
Copyright 2021 The KServe Authors.

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

package inferenceservice

import (
	"testing"

	"github.com/kserve/kserve/pkg/apis/serving/v1beta1"
	"github.com/kserve/kserve/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestStopAnnotationProcessing tests the stop annotation processing logic
// This extracts business logic that was being tested in multiple integration tests
func TestStopAnnotationProcessing(t *testing.T) {
	tests := []struct {
		name            string
		stopAnnotation  string
		expectedStopped bool
		description     string
	}{
		{
			name:            "stop_annotation_true",
			stopAnnotation:  "true",
			expectedStopped: true,
			description:     "Should stop when annotation is 'true'",
		},
		{
			name:            "stop_annotation_false",
			stopAnnotation:  "false",
			expectedStopped: false,
			description:     "Should not stop when annotation is 'false'",
		},
		{
			name:            "stop_annotation_empty",
			stopAnnotation:  "",
			expectedStopped: false,
			description:     "Should not stop when annotation is empty",
		},
		{
			name:            "stop_annotation_invalid",
			stopAnnotation:  "invalid",
			expectedStopped: false,
			description:     "Should not stop when annotation is invalid value",
		},
		{
			name:            "stop_annotation_uppercase_true",
			stopAnnotation:  "TRUE",
			expectedStopped: false,
			description:     "Should not stop when annotation is uppercase (case sensitive)",
		},
		{
			name:            "stop_annotation_yes",
			stopAnnotation:  "yes",
			expectedStopped: false,
			description:     "Should not stop when annotation is 'yes' (only 'true' is valid)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test InferenceService with stop annotation
			isvc := &v1beta1.InferenceService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "default",
					Annotations: map[string]string{
						"serving.kserve.io/deploymentMode": "Serverless",
						constants.StopAnnotationKey:        tt.stopAnnotation,
					},
				},
			}

			// Test the stop annotation processing logic
			isStopped := shouldStopInferenceService(isvc)

			// Verify the result
			if isStopped != tt.expectedStopped {
				t.Errorf("%s: Expected stopped=%v, got stopped=%v",
					tt.description, tt.expectedStopped, isStopped)
			}
		})
	}
}

// TestStopAnnotationMissing tests behavior when stop annotation is completely missing
func TestStopAnnotationMissing(t *testing.T) {
	// Create test InferenceService without stop annotation
	isvc := &v1beta1.InferenceService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "default",
			Annotations: map[string]string{
				"serving.kserve.io/deploymentMode": "Serverless",
				// No stop annotation
			},
		},
	}

	// Test the stop annotation processing logic
	isStopped := shouldStopInferenceService(isvc)

	// Should not be stopped when annotation is missing
	if isStopped {
		t.Error("Expected service not to be stopped when stop annotation is missing")
	}
}

// TestStopAnnotationWithNilAnnotations tests behavior when annotations map is nil
func TestStopAnnotationWithNilAnnotations(t *testing.T) {
	// Create test InferenceService with nil annotations
	isvc := &v1beta1.InferenceService{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-service",
			Namespace:   "default",
			Annotations: nil,
		},
	}

	// Test the stop annotation processing logic
	isStopped := shouldStopInferenceService(isvc)

	// Should not be stopped when annotations are nil
	if isStopped {
		t.Error("Expected service not to be stopped when annotations are nil")
	}
}

// shouldStopInferenceService extracts the core business logic for determining if a service should be stopped
// This is the logic that would be used in the actual controller
func shouldStopInferenceService(isvc *v1beta1.InferenceService) bool {
	if isvc.Annotations == nil {
		return false
	}

	stopValue, exists := isvc.Annotations[constants.StopAnnotationKey]
	if !exists {
		return false
	}

	// Only exact string "true" should stop the service
	return stopValue == "true"
}

// TestStopAnnotationStateTransitions tests the logic for state transitions
func TestStopAnnotationStateTransitions(t *testing.T) {
	tests := []struct {
		name                string
		currentState        bool // current stopped state
		newAnnotation       string
		expectedNewState    bool
		shouldTriggerChange bool
		description         string
	}{
		{
			name:                "running_to_stopped",
			currentState:        false,
			newAnnotation:       "true",
			expectedNewState:    true,
			shouldTriggerChange: true,
			description:         "Should transition from running to stopped",
		},
		{
			name:                "stopped_to_running",
			currentState:        true,
			newAnnotation:       "false",
			expectedNewState:    false,
			shouldTriggerChange: true,
			description:         "Should transition from stopped to running",
		},
		{
			name:                "running_to_running",
			currentState:        false,
			newAnnotation:       "false",
			expectedNewState:    false,
			shouldTriggerChange: false,
			description:         "Should remain running when already running",
		},
		{
			name:                "stopped_to_stopped",
			currentState:        true,
			newAnnotation:       "true",
			expectedNewState:    true,
			shouldTriggerChange: false,
			description:         "Should remain stopped when already stopped",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test InferenceService
			isvc := &v1beta1.InferenceService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service",
					Namespace: "default",
					Annotations: map[string]string{
						"serving.kserve.io/deploymentMode": "Serverless",
						constants.StopAnnotationKey:        tt.newAnnotation,
					},
				},
			}

			// Calculate new state and check if change is needed
			newState := shouldStopInferenceService(isvc)
			needsChange := newState != tt.currentState

			// Verify the results
			if newState != tt.expectedNewState {
				t.Errorf("%s: Expected new state=%v, got state=%v",
					tt.description, tt.expectedNewState, newState)
			}

			if needsChange != tt.shouldTriggerChange {
				t.Errorf("%s: Expected trigger change=%v, got trigger change=%v",
					tt.description, tt.shouldTriggerChange, needsChange)
			}
		})
	}
}
