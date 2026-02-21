package v1beta1

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"
)

// TODO hilft das überhaupt? Es muss ja zeigen das die CRD zum YAML passt?
func Test_example_resource_should_match_crd_struct(t *testing.T) {
	file, err := os.ReadFile("../../config/examples/webpage.yaml")
	require.NoError(t, err, "Failed to read example file")

	webpage := &WebPage{}
	err = yaml.Unmarshal(file, &webpage)
	require.NoError(t, err, "Failed to unmarshal example file into WebPage struct")

	assert.Equal(t, "this is a test", webpage.Spec.Content)
}
