package resource

import (
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestBuildInfo(t *testing.T) {
	tests := []struct {
		description string
		in          utils.ResultItem
		expected    Version
		expectError bool
	}{
		{
			description: "empty result item",
			in:          utils.ResultItem{},
			expected:    Version{},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, err := processItem(tc.in)

			if tc.expectError {
				Expect(t, err).To(Not(BeNil()))
				return
			}

			Expect(t, err).To(BeNil())
			Expect(t, out).To(Equal(tc.expected))
		})
	}
}
