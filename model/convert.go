package model

type ConvertRequest struct {
	Subscriptions []string          `form:"subscriptions" json:"subscriptions"`
	Proxies       []string          `form:"proxy" json:"proxy"`
	Template      string            `form:"template" json:"template"`
	Delete        string            `form:"delete" json:"delete"`
	Rename        map[string]string `form:"rename" json:"rename"`
	Group         bool              `form:"group" json:"group"`
	GroupType     string            `form:"group-type" json:"group-type"`
	SortKey       string            `form:"sort" json:"sort"`
	SortType      string            `form:"sort-type" json:"sort-type"`
	Output        string            `json:"output"`
	GroupRules    string            `form:"group-rules" json:"group-rules"`
}
