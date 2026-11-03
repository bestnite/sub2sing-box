package handler

import (
	"encoding/json"

	"github.com/bestnite/sub2sing-box/common"
	"github.com/bestnite/sub2sing-box/model"
	"github.com/bestnite/sub2sing-box/util"

	"github.com/gin-gonic/gin"
)

func Convert(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	if c.Query("data") == "" {
		c.JSON(400, gin.H{
			"error": "Missing data parameter",
		})
		return
	}
	j, err := util.DecodeBase64(c.Query("data"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid data",
		})
		return
	}
	data, err := model.ParseConvertRequest([]byte(j))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if data.Proxies == nil && data.Subscriptions == nil {
		c.JSON(400, gin.H{
			"error": "Must provide at least one subscription or proxy",
		})
		return
	}
	groupRules := make(map[string][]string)
	if data.GroupRules != "" {
		err = json.Unmarshal([]byte(data.GroupRules), &groupRules)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	result, err := common.Convert(
		data.Subscriptions,
		data.Proxies,
		data.Template,
		data.Delete,
		data.Rename,
		data.Group,
		data.GroupType,
		data.SortKey,
		data.SortType,
		groupRules,
		data.UserAgent,
	)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.String(200, result)
}
