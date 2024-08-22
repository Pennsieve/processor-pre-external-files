package preprocessor

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromEnv(t *testing.T) {

	expectedIntegrationID := uuid.NewString()
	t.Setenv(IntegrationIDKey, expectedIntegrationID)
	expectedInputDirectory := fmt.Sprintf("input/%s", uuid.NewString())
	t.Setenv(InputDirectoryKey, expectedInputDirectory)
	expectedOutputDirectory := fmt.Sprintf("output/%s", uuid.NewString())
	t.Setenv(OutputDirectoryKey, expectedOutputDirectory)

	processor, err := FromEnv()
	require.NoError(t, err)

	assert.Equal(t, expectedIntegrationID, processor.IntegrationID)
	assert.Equal(t, expectedInputDirectory, processor.InputDirectory)
	assert.Equal(t, expectedOutputDirectory, processor.OutputDirectory)

}
