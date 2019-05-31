package docker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTag(t *testing.T) {
	execStringCommand = func (command string) (string, error) {
		assert.Equal(t, "docker tag sourceimage:sourcetag targetimage:targettag", command)
		return "", nil
	}

	err := Tag("sourceImage", "sourceTag", "targetImage", "targetTag")
	if err != nil {
		t.Error(err)
	}
}

func TestTag_emptySourceTag(t *testing.T) {
	execStringCommand = func (command string) (string, error) {
		assert.Equal(t, "docker tag sourceimage targetimage:targettag", command)
		return "", nil
	}

	err := Tag("sourceImage", "", "targetImage", "targetTag")
	if err != nil {
		t.Error(err)
	}
}

func  TestTag_emptyTargetTag(t *testing.T) {
	execStringCommand = func (command string) (string, error) {
		assert.Equal(t, "docker tag sourceimage:sourcetag targetimage", command)
		return "", nil
	}

	err := Tag("sourceImage", "sourceTag", "targetImage", "")
	if err != nil {
		t.Error(err)
	}
}

func TestPush(t *testing.T) {
	execCommandStreamingOut = func (command string) error {
		assert.Equal(t, "docker push registryurl/imagename:tag", command)
		return nil
	}

	err := Push("registryUrl", "imageName", "tag")
	if err != nil {
		t.Error(err)
	}
}

func TestPush_emptyTag(t *testing.T) {
	execCommandStreamingOut = func (command string) error {
		assert.Equal(t, "docker push registryurl/imagename", command)
		return nil
	}

	err := Push("registryUrl", "imageName", "")
	if err != nil {
		t.Error(err)
	}
}

func TestBuild(t *testing.T) {
	execCommandStreamingOut = func (command string) error {
		assert.Equal(t, "docker build --pull --shm-size 256m --memory=3g --memory-swap=-1 -t imagename:commit .", command)
		return nil
	}

	err := Build("imageName", "commit")
	if err != nil {
		t.Error(err)
	}
}

func TestBuild_emptyComitTag(t *testing.T) {
	execCommandStreamingOut = func (command string) error {
		assert.Equal(t, "docker build --pull --shm-size 256m --memory=3g --memory-swap=-1 -t imagename .", command)
		return nil
	}

	err := Build("imageName", "")
	if err != nil {
		t.Error(err)
	}
}

func TestBuild_noTag(t *testing.T) {
	execCommandStreamingOut = func (command string) error {
		assert.Equal(t, "docker build --pull --shm-size 256m --memory=3g --memory-swap=-1 .", command)
		return nil
	}

	err := Build("", "")
	if err != nil {
		t.Error(err)
	}
}

func TestLogin(t *testing.T) {
	execStringCommand = func (command string) (string, error) {
		assert.Equal(t, "docker login -u Us3r -p PaSSword234! registryUrl", command)
		return "", nil
	}

	err := Login("registryUrl", "Us3r", "PaSSword234!")
	if err != nil {
		t.Error(err)
	}
}