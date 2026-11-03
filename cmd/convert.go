package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bestnite/sub2sing-box/common"
	"github.com/bestnite/sub2sing-box/model"

	"github.com/spf13/cobra"
)

var (
	subscriptions []string
	proxies       []string
	template      string
	output        string
	delete        string
	rename        map[string]string
	group         bool
	groupType     string
	sortKey       string
	sortType      string
	config        string
	groupRules    string
	userAgent     string
)

func init() {
	convertCmd.Flags().StringSliceVarP(&subscriptions, "subscription", "s", make([]string, 0), "subscription URLs")
	convertCmd.Flags().StringSliceVarP(&proxies, "proxy", "p", make([]string, 0), "proxies share links")
	convertCmd.Flags().StringVarP(&template, "template", "t", "", "template file path or URL")
	convertCmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	convertCmd.Flags().StringVarP(&delete, "delete", "d", "", "delete proxies with regex")
	convertCmd.Flags().StringToStringVarP(&rename, "rename", "r", make(map[string]string), "rename proxies with regex")
	convertCmd.Flags().BoolVarP(&group, "group", "g", false, "grouping proxies by country")
	convertCmd.Flags().StringVarP(&groupType, "group-type", "G", "", "group type, selector or urltest")
	convertCmd.Flags().StringVarP(&sortKey, "sort", "S", "", "sort key, tag or num")
	convertCmd.Flags().StringVarP(&sortType, "sort-type", "T", "", "sort type, asc or desc")
	convertCmd.Flags().StringVarP(&config, "config", "c", "", "configuration file path")
	convertCmd.Flags().StringVarP(&groupRules, "group-rules", "R", "", "group rules")
	convertCmd.Flags().StringVarP(&userAgent, "user-agent", "u", "", "custom User-Agent for fetching subscriptions and remote templates")
	RootCmd.AddCommand(convertCmd)
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Long:  "Convert common proxy to sing-box proxy",
	Short: "Convert common proxy to sing-box proxy",
	Run:   convertRun,
}

func convertRun(cmd *cobra.Command, args []string) {
	if err := loadConfig(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	groupRulesMap := make(map[string][]string)
	if groupRules != "" {
		err := json.Unmarshal([]byte(groupRules), &groupRulesMap)
		if err != nil {
			fmt.Println("Error parsing group rules:", err)
			os.Exit(1)
		}
	}
	result, err := common.Convert(
		subscriptions,
		proxies,
		template,
		delete,
		rename,
		group,
		groupType,
		sortKey,
		sortType,
		groupRulesMap,
		userAgent,
	)
	if err != nil {
		fmt.Println("Conversion error:", err)
		os.Exit(1)
	}
	if output != "" {
		err = os.WriteFile(output, []byte(result), 0666)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			os.Exit(1)
		}
		fmt.Println("Config has been saved in:", output)
	} else {
		fmt.Println(result)
	}
}

func loadConfig() error {
	if config == "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting working directory: %w", err)
		}
		config = filepath.Join(wd, "sub2sing-box.json")
		if _, err := os.Stat(config); os.IsNotExist(err) {
			return nil
		}
	}

	bytes, err := os.ReadFile(config)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	cfg, err := model.ParseConvertRequest(bytes)
	if err != nil {
		return fmt.Errorf("parsing config JSON: %w", err)
	}

	mergeConfig(cfg)
	return nil
}

func mergeConfig(cfg model.ConvertRequest) {
	subscriptions = append(subscriptions, cfg.Subscriptions...)
	proxies = append(proxies, cfg.Proxies...)
	if template == "" {
		template = cfg.Template
	}
	if delete == "" {
		delete = cfg.Delete
	}
	for k, v := range cfg.Rename {
		rename[k] = v
	}
	if groupType == "" {
		if cfg.GroupType != "" {
			groupType = cfg.GroupType
		} else {
			groupType = "selector"
		}
	}
	if sortKey == "" {
		if cfg.SortKey != "" {
			sortKey = cfg.SortKey
		} else {
			sortKey = "tag"
		}
	}
	if sortType == "" {
		if cfg.SortType != "" {
			sortType = cfg.SortType
		} else {
			sortType = "asc"
		}
	}
	if output == "" {
		output = cfg.Output
	}
	if userAgent == "" {
		userAgent = cfg.UserAgent
	}
}
