package parser

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bestnite/sub2sing-box/constant"
	"github.com/bestnite/sub2sing-box/model"
	"github.com/bestnite/sub2sing-box/util"
	"github.com/sagernet/sing-box/option"
)

func ParseSocks(proxy string) (model.Outbound, error) {
	if !strings.HasPrefix(proxy, constant.SocksPrefix) {
		return model.Outbound{}, &ParseError{Type: ErrInvalidPrefix, Raw: proxy}
	}
	link, err := url.Parse(proxy)
	if err != nil {
		return model.Outbound{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "url parse error",
			Raw:     proxy,
		}
	}
	server := link.Hostname()
	if server == "" {
		return model.Outbound{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server host",
			Raw:     proxy,
		}
	}
	portStr := link.Port()
	if portStr == "" {
		return model.Outbound{}, &ParseError{
			Type:    ErrInvalidStruct,
			Message: "missing server port",
			Raw:     proxy,
		}
	}
	port, err := ParsePort(portStr)
	if err != nil {
		return model.Outbound{}, &ParseError{
			Type: ErrInvalidPort,
			Raw:  portStr,
		}
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	encodeStr := link.User.Username()
	var username, password string
	if encodeStr != "" {
		decodeStr, err := util.DecodeBase64(encodeStr)
		splitStr := strings.Split(decodeStr, ":")
		if err != nil {
			return model.Outbound{}, &ParseError{
				Type:    ErrInvalidStruct,
				Message: "url parse error",
				Raw:     proxy,
			}
		}
		username = splitStr[0]
		if len(splitStr) == 2 {
			password = splitStr[1]
		}
	}

	outboundOptions := option.SOCKSOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     server,
			ServerPort: port,
		},
		Username: username,
		Password: password,
	}

	return model.Outbound{
		Type:    "socks",
		Tag:     remarks,
		Options: outboundOptions,
	}, nil
}
