package main

import (
	"io"
	"regexp"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, reader io.Reader, output io.Writer) (*model.FileInfo, string) {
	policy := p.getConfiguration().AttachmentPolicy

	users, err := p.GetChannelUsers(info)
	if err != nil {
		p.API.LogError(
			"Failed to query GetChannelMembers in FileWillBeUploaded",
			"error", err,
			"infopath", info.Path,
		)
		return nil, "Something went wrong with the upload. Please report it to the administrator."
	}
	attachmentInfo := new(AttachmentInfo)
	attachmentInfo.ChannelUsers = users
	attachmentInfo.FileInfo = info
	attachmentInfo.Context = c
	result, err := apply(policy, attachmentInfo)
	if err != nil {
		json, _ := infoToString(attachmentInfo)
		p.API.LogError(
			"Failed to query JsonLogic.apply in FileWillBeUploaded",
			"error", err,
			"attachmentinfo", json,
		)
		return nil, "Something went wrong with the attachment filtering process. Please report it to the administrator."
	}
	if result {
		return info, ""
	} else {
		json, _ := infoToString(attachmentInfo)
		p.API.LogInfo(
			"File uploads were not allowed in FileWillBeUploaded",
			"attachmentinfo", json,
		)
		return nil, "File uploads were not allowed due to attachment policy."
	}
}

func (p *Plugin) GetChannelUsers(info *model.FileInfo) ([]*model.User, error) {
	path := info.Path
	r := regexp.MustCompile("^.+/channels/([^/]+)/.+$")
	fss := r.FindStringSubmatch(path)
	channelID := fss[1]

	ms, err := p.API.GetChannelMembers(channelID, 0, 1000)
	if err != nil {
		p.API.LogError(
			"Failed to query GetChannelMembers in FileWillBeUploaded",
			"error", err.Error(),
		)
		return nil, err
	}

	members := ([]model.ChannelMember)(*ms)
	users := make([]*model.User, len(members))
	for i, member := range members {
		user, err2 := p.API.GetUser(member.UserId)
		if err2 != nil {
			p.API.LogError(
				"Failed to query GetUser in FileWillBeUploaded",
				"error", err2.Error(),
			)
			return nil, err2
		}
		users[i] = user
	}

	return users, nil
}
