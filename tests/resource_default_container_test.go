package tests

import (
	"testing"

	"github.com/UCF/terraform-ecr-default-container/provider"
	"github.com/stretchr/testify/assert"
)

func TestPushDefaultContainer(t *testing.T) {
	repoName := "test-repo"
	err := provider.PushDefaultContainer(repoName)
	assert.NoError(t, err, "Expected no error when pushing default container")

}
