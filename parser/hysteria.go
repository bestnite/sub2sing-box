package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/nitezs/sub2sing-box/constant"
	"github.com/nitezs/sub2sing-box/model"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/byteformats"
)

func ParseHysteria(proxy string) (model.Outbound, error) {
	if !strings.HasPrefix(proxy, constant.HysteriaPrefix) {
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
			Type:    ErrInvalidPort,
			Message: err.Error(),
			Raw:     proxy,
		}
	}

	query := link.Query()

	protocol, auth, insecure, upmbps, downmbps, obfs, alpnStr := query.Get("protocol"), query.Get("auth"), query.Get("insecure"), query.Get("upmbps"), query.Get("downmbps"), query.Get("obfs"), query.Get("alpn")
	insecureBool, err := strconv.ParseBool(insecure)
	enableTLS := insecureBool || alpnStr != ""
	if err != nil {
		insecureBool = false
	}

	var alpn []string
	alpnStr = strings.TrimSpace(alpnStr)
	if alpnStr != "" {
		alpn = strings.Split(alpnStr, ",")
	}

	remarks := link.Fragment
	if remarks == "" {
		remarks = fmt.Sprintf("%s:%s", server, portStr)
	}
	remarks = strings.TrimSpace(remarks)

	up := &byteformats.NetworkBytesCompat{}
	err = json.Unmarshal(fmt.Appendf(nil, `"%s"`, upmbps), up)
	if err != nil {
		return model.Outbound{}, &ParseError{
			Type:    ErrInvalidNetworkBytesCompat,
			Message: err.Error(),
			Raw:     proxy,
		}
	}
	down := &byteformats.NetworkBytesCompat{}
	err = json.Unmarshal(fmt.Appendf(nil, `"%s"`, downmbps), down)
	if err != nil {
		return model.Outbound{}, &ParseError{
			Type:    ErrInvalidNetworkBytesCompat,
			Message: err.Error(),
			Raw:     proxy,
		}
	}
	outboundOptions := option.HysteriaOutboundOptions{
		ServerOptions: option.ServerOptions{
			Server:     server,
			ServerPort: port,
		},
		Up:      up,
		Down:    down,
		Auth:    []byte(auth),
		Obfs:    obfs,
		Network: option.NetworkList(protocol),
		OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
			TLS: &option.OutboundTLSOptions{
				Enabled:  enableTLS,
				Insecure: insecureBool,
				ALPN:     alpn,
			},
		},
	}

	return model.Outbound{
		Type:    "hysteria",
		Tag:     remarks,
		Options: outboundOptions,
	}, nil
}
