package resource

import (
	"testing"
	"time"

	"github.com/digitalocean/artifactory-resource/internal"
	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestAddModifiedTime(t *testing.T) {
	tests := []struct {
		description string
		aql         AQL
		version     Version
		expected    string
	}{
		{
			description: "empty aql",
			aql:         AQL{},
			version:     Version{},
			expected:    "",
		},
		{
			description: "simple",
			aql:         AQL{Raw: `{"repo": "artifacts-local", "path": {"$match" : "changeset/*"}, "name": "artifact"}`},
			version:     Version{Modified: internal.GetTimePointer(time.Date(2020, time.May, 26, 0, 0, 0, 0, time.UTC))},
			expected:    `{"repo": "artifacts-local", "path": {"$match" : "changeset/*"}, "name": "artifact", "modified": {"$gt": "2020-05-26T00:00:00Z"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			tc.aql.SetModifiedTime(tc.version)
			out := tc.aql.Raw
			Expect(t, out).To(Equal(tc.expected))
		})
	}
}

func TestCheckRequestUnmarshal(t *testing.T) {
	tests := []struct {
		description   string
		input         []byte
		expected      CheckRequest
		errorExpected bool
	}{
		{
			description:   "empty string",
			input:         []byte("{}"),
			expected:      CheckRequest{},
			errorExpected: false,
		},
		{
			description: "simple source w/empty version",
			input: []byte(`
			{
				"source": {
					"endpoint": "https://artifactory.example.com",
					"user": "me",
					"password": "xxxx"
				},
				"version": {}
			}
			`),
			expected: CheckRequest{
				Source: Source{
					Endpoint: "https://artifactory.example.com",
					User:     "me",
					Password: "xxxx",
				},
			},
			errorExpected: false,
		},
		{
			description: "simple source w/empty version",
			input: []byte(`
			{
				"source": {
					"endpoint": "https://artifactory.example.com",
					"user": "me",
					"password": "xxxx",
					"aql": {
						"raw": "{\"repo\": \"artifacts-local\", \"path\": {\"$match\": \"project/*\"}, \"name\": \"artifact\"}"
					}
				},
				"version": {}
			}
			`),
			expected: CheckRequest{
				Source: Source{
					Endpoint: "https://artifactory.example.com",
					User:     "me",
					Password: "xxxx",
					AQL: AQL{
						Raw: "{\"repo\": \"artifacts-local\", \"path\": {\"$match\": \"project/*\"}, \"name\": \"artifact\"}",
					},
				},
			},
			errorExpected: false,
		},
		{
			description: "simple source w/empty version",
			input: []byte(`
			{
				"source": {
					"endpoint": "https://artifactory.example.com",
					"user": "me",
					"password": "xxxx",
					"aql": {
						"repo": "artifacts-local",
						"path": "project/*",
						"name": "artifact"
					}
				},
				"version": {}
			}
			`),
			expected: CheckRequest{
				Source: Source{
					Endpoint: "https://artifactory.example.com",
					User:     "me",
					Password: "xxxx",
					AQL: AQL{
						Raw:  "{\"repo\": \"artifacts-local\", \"path\": {\"$match\": \"project/*\"}, \"name\": {\"$match\": \"artifact\"}}",
						Repo: "artifacts-local",
						Path: "project/*",
						Name: "artifact",
					},
				},
			},
			errorExpected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			r := CheckRequest{}
			err := r.Read(tc.input)
			if tc.errorExpected {
				Expect(t, err).To(Not(BeNil()))
			}

			Expect(t, err).To(BeNil())
			Expect(t, r).To(Equal(tc.expected))
		})
	}
}
