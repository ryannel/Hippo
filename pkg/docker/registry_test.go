package docker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRegistryDomain(t *testing.T) {
	domain := GetRegistryDomain("quay")
	assert.Equal(t, "quay.io", domain)
}

func TestBuildReigistryUrl(t *testing.T) {
	url := BuildReigistryUrl("quay","namespace", "repo")
	assert.Equal(t, "https://quay.io/namespace/repo", url)
}