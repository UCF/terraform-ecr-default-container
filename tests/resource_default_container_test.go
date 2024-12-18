package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushDefaultContainer(t *testing.T) {
	repoName := "test-repo"
	err := PushDefaultContainer(repoName)
	assert.NoError(t, err, "Expected no error when pushing default container")

}
