package resource

import (
	"testing"

	"github.com/digitalocean/concourse-resource-library/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestMetadata(t *testing.T) {
	tests := []struct {
		description string
		input       artifactory.Artifact
		expected    string
	}{
		{
			description: "empty artifact",
			input:       artifactory.Artifact{File: utils.FileInfo{}, Item: utils.ResultItem{}},
			expected:    `[{"name":"artifactory-path","value":""},{"name":"local-path","value":""},{"name":"created","value":""},{"name":"modified","value":""},{"name":"name","value":""},{"name":"repo","value":""},{"name":"size","value":"0"},{"name":"type","value":""},{"name":"properties","value":"null"}]`,
		},
		{
			description: "filled properties",
			input: artifactory.Artifact{
				File: utils.FileInfo{},
				Item: utils.ResultItem{
					Properties: []utils.Property{
						{Key: "vcs.revision", Value: "xxxxx"},
					},
				},
			},
			expected: `[{"name":"artifactory-path","value":""},{"name":"local-path","value":""},{"name":"created","value":""},{"name":"modified","value":""},{"name":"name","value":""},{"name":"repo","value":""},{"name":"size","value":"0"},{"name":"type","value":""},{"name":"properties","value":"[{\"Key\":\"vcs.revision\",\"Value\":\"xxxxx\"}]"}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out := metadata(tc.input)

			data, err := out.JSON()
			Expect(t, err).To(BeNil())
			Expect(t, string(data)).To(Equal(tc.expected))
		})
	}
}
