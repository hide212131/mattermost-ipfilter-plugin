package main

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/diegoholiveira/jsonlogic/v3"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type AttachmentInfo struct {
	ChannelUsers []*model.User
	FileInfo     *model.FileInfo
	Context      *plugin.Context
}

var initialized bool = false

func initialize() {
	jsonlogic.AddOperator("match", func(values, data interface{}) interface{} {
		v, ok := values.([]interface{})
		if ok {
			v0, ok0 := v[0].(string)
			v1, ok1 := v[1].(string)
			if ok0 && ok1 {
				re := regexp.MustCompile(v1)
				return re.MatchString(v0)
			}
		}
		return 0
	})
	initialized = true
}

func apply(policy string, info *AttachmentInfo) (bool, string) {
	if !initialized {
		initialize()
	}
	logic := strings.NewReader(policy)
	jsonb, err := json.Marshal(info)
	if err != nil {
		return false, "The JSON format of the policy definition is incorrect."
	}
	data := strings.NewReader(string(jsonb))
	var result bytes.Buffer
	err = jsonlogic.Apply(logic, data, &result)
	if err == nil {
		r := result.String()
		switch r {
		case "true\n":
			return true, ""
		case "false\n":
			return false, ""
		default:
			return false, "The JSON format of the policy definition does not return true or false."
		}
	} else {
		return false, err.Error()
	}
}
