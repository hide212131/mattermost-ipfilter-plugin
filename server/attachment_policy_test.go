package main

import (
	"io/ioutil"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/stretchr/testify/assert"
)

func TestAttachmentPolicy(t *testing.T) {
	initialize()

	bts, err := ioutil.ReadFile("attachment_policy_test.json")
	if err != nil {
		t.Fatal(err)
	}
	policy := string(bts)

	cm := make([]*model.User, 1)
	cm[0] = &model.User{
		Username: "Red",
	}

	info := AttachmentInfo{
		ChannelUsers: cm,
		Context: &plugin.Context{
			IpAddress: "192.168.0.100",
		},
		FileInfo: &model.FileInfo{
			Name: "test.png",
		},
	}

	result, errmsg := apply(policy, &info)
	assert.Equal(t, "", errmsg)
	assert.True(t, result)
}
