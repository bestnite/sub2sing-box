package model

import (
	"bytes"
	"encoding/json"
)

type ConvertRequest struct {
	Subscriptions []string          `form:"subscription" json:"subscriptions"`
	Proxies       []string          `form:"proxy" json:"proxies"`
	Template      string            `form:"template" json:"template"`
	Delete        string            `form:"delete" json:"delete"`
	Rename        map[string]string `form:"rename" json:"rename"`
	Group         bool              `form:"group" json:"group"`
	GroupType     string            `form:"group-type" json:"group-type"`
	SortKey       string            `form:"sort" json:"sort"`
	SortType      string            `form:"sort-type" json:"sort-type"`
	Output        string            `json:"output"`
	GroupRules    string            `form:"group-rules" json:"group-rules"`
	UserAgent     string            `form:"user-agent" json:"user-agent"`
}

// ParseConvertRequest decodes a convert request.
//
// Unknown fields are rejected on purpose: a misspelled key used to be dropped
// silently, which turned a wrong "subscription" into an empty subscriptions
// list and produced a config without any node.
func ParseConvertRequest(data []byte) (ConvertRequest, error) {
	var cfg ConvertRequest
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return ConvertRequest{}, err
	}
	return cfg, nil
}
