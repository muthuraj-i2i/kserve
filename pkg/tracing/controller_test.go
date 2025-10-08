/*
Copyright 2023 The KServe Authors.

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

package tracing

import "testing"

func TestIsTruthy(t *testing.T) {
	cases := map[string]bool{
		"true":  true,
		"TRUE":  true,
		"1":     true,
		"yes":   true,
		"on":    true,
		"false": false,
		"0":     false,
		"":      false,
		"no":    false,
	}

	for input, expected := range cases {
		if actual := isTruthy(input); actual != expected {
			t.Fatalf("isTruthy(%q) = %t, expected %t", input, actual, expected)
		}
	}
}

func TestIsTruthyWithDefault(t *testing.T) {
	if !isTruthyWithDefault("", true) {
		t.Fatalf("expected default true when value empty")
	}
	if isTruthyWithDefault("", false) {
		t.Fatalf("expected default false when value empty")
	}
	if !isTruthyWithDefault("true", false) {
		t.Fatalf("expected override true when value explicit")
	}
	if isTruthyWithDefault("no", true) {
		t.Fatalf("expected override false when value explicit")
	}
}
