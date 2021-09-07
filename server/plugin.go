package main

import (
	"net/http"
	"strings"
	"sync"

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

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	addr := c.IpAddress
	AllowedIps := p.configuration.AllowedIps
	if strings.TrimSpace(AllowedIps) != "" {
		ips := strings.Split(AllowedIps, ",")
		for _, s := range ips {
			text := strings.TrimRight(strings.TrimSpace(s), "*")
			if strings.HasPrefix(addr, text) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				p.API.LogInfo("Allowed from " + addr + " matching prefix " + s)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("NG"))
	p.API.LogInfo("Blocked from " + addr)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
