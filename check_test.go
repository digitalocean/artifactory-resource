package resource

import (
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestProcessItem(t *testing.T) {
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
		{
			description: "bad date time",
			in: utils.ResultItem{
				Repo:       "artifact-local",
				Path:       "some/path",
				Name:       "artifact",
				Modified:   "2020-05-2620:00:00.000Z",
				Properties: []utils.Property{},
				Type:       "file",
			},
			expected: Version{
				Repo:     "artifact-local",
				Path:     "some/path",
				Name:     "artifact",
				Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
			},
			expectError: true,
		},
		{
			description: "valid data",
			in: utils.ResultItem{
				Repo:       "artifact-local",
				Path:       "some/path",
				Name:       "artifact",
				Modified:   "2020-05-26T20:00:00.000Z",
				Properties: []utils.Property{},
				Type:       "file",
			},
			expected: Version{
				Repo:     "artifact-local",
				Path:     "some/path",
				Name:     "artifact",
				Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
			},
			expectError: false,
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

func TestSelectVersions(t *testing.T) {
	tests := []struct {
		description string
		input       Version
		found       CheckResponse
		expected    CheckResponse
	}{
		{
			description: "empty input & empty found",
			input:       Version{},
			found:       CheckResponse{},
			expected:    CheckResponse{},
		},
		{
			description: "no new versions found",
			input: Version{
				Repo:     "artifact-local",
				Path:     "some/path",
				Name:     "artifact",
				Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
			},
			found: CheckResponse{},
			expected: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			description: "new version found, no input",
			input:       Version{},
			found: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
				},
			},
			expected: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			description: "new versions found, no input",
			input:       Version{},
			found: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
				},
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 27, 20, 0, 0, 0, time.UTC),
				},
			},
			expected: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 27, 20, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			description: "new versions found with input",
			input: Version{
				Repo:     "artifact-local",
				Path:     "some/path",
				Name:     "artifact",
				Modified: time.Date(2020, time.May, 25, 20, 0, 0, 0, time.UTC),
			},
			found: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
				},
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 27, 20, 0, 0, 0, time.UTC),
				},
			},
			expected: []Version{
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 26, 20, 0, 0, 0, time.UTC),
				},
				{
					Repo:     "artifact-local",
					Path:     "some/path",
					Name:     "artifact",
					Modified: time.Date(2020, time.May, 27, 20, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := selectVersions(tc.input, tc.found)
			Expect(t, out).To(Equal(tc.expected))
		})
	}
}
