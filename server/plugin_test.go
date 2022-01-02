package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileWillBeUpload(t *testing.T) {
	t.Run("png", func(t *testing.T) {
		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			return api
		}

		api := setupAPI()
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		cm := make([]model.ChannelMember, 1)
		cm[0] = model.ChannelMember{
			UserId: "123",
		}
		cms := (*model.ChannelMembers)(&cm)
		api.On("GetChannelMembers", mock.Anything, mock.Anything, mock.Anything).Return(cms, nil)
		api.On("GetUser", mock.Anything).Return(&model.User{Username: "Red"}, nil)

		defer api.AssertExpectations(t)
		p := &Plugin{}
		p.API = api

		bts, err := ioutil.ReadFile("attachment_policy_test.json")
		if err != nil {
			t.Fatal(err)
		}
		p.configuration = &configuration{AttachmentPolicy: string(bts)}
		// "{ \"match\": [{ \"var\": \"FileInfo.name\" }, \"\\\\.png$\"] }"}

		data, err := ioutil.ReadFile("../assets/test.png")
		assert.Nil(t, err)

		c := &plugin.Context{
			IpAddress: "192.168.0.100",
		}

		fi := &model.FileInfo{
			Name: "test.png",
			Path: "foo/bar/channels/dummychannels/baz",
		}

		r := bytes.NewReader(data)

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		fi, reason := p.FileWillBeUploaded(c, fi, r, w)
		assert.Equal(t, reason, "")
		assert.NotNil(t, fi)
	})
}
