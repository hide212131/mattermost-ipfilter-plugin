package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{}
	//configuration := &configuration{AllowedIps: "192.0.*, 127.*"}
	configuration := &configuration{AllowedIps: "10.*"}
	plugin.setConfiguration(configuration)

	api := &plugintest.API{}
	plugin.SetAPI(api)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.NotNil(result)
	defer result.Body.Close()
	bodyBytes, err := ioutil.ReadAll(result.Body)
	assert.Nil(err)
	bodyString := string(bodyBytes)

	assert.Equal("OK", bodyString)
}
