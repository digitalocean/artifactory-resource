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
			expected:    `[{"name":"ArtifactoryPath","value":""},{"name":"LocalPath","value":""},{"name":"Created","value":""},{"name":"Modified","value":""},{"name":"Name","value":""},{"name":"Repo","value":""},{"name":"Size","value":"0"},{"name":"Type","value":""},{"name":"Properties","value":"null"}]`,
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
			expected: `[{"name":"ArtifactoryPath","value":""},{"name":"LocalPath","value":""},{"name":"Created","value":""},{"name":"Modified","value":""},{"name":"Name","value":""},{"name":"Repo","value":""},{"name":"Size","value":"0"},{"name":"Type","value":""},{"name":"Properties","value":"[{\"Key\":\"vcs.revision\",\"Value\":\"xxxxx\"}]"}]`,
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
