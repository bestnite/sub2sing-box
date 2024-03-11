package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sub2sing-box/model"
)

func ParseVless(proxy string) (model.Proxy, error) {
	if !strings.HasPrefix(proxy, "vless://") {
		return model.Proxy{}, fmt.Errorf("invalid vless Url")
	}
	parts := strings.SplitN(strings.TrimPrefix(proxy, "vless://"), "@", 2)
	if len(parts) != 2 {
		return model.Proxy{}, fmt.Errorf("invalid vless Url")
	}
	serverInfo := strings.SplitN(parts[1], "#", 2)
	serverAndPortAndParams := strings.SplitN(serverInfo[0], "?", 2)
	serverAndPort := strings.SplitN(serverAndPortAndParams[0], ":", 2)
	params, err := url.ParseQuery(serverAndPortAndParams[1])
	if err != nil {
		return model.Proxy{}, err
	}
	if len(serverAndPort) != 2 {
		return model.Proxy{}, fmt.Errorf("invalid vless")
	}
	port, err := strconv.Atoi(strings.TrimSpace(serverAndPort[1]))
	if err != nil {
		return model.Proxy{}, err
	}
	remarks := ""
	if len(serverInfo) == 2 {
		if strings.Contains(serverInfo[1], "|") {
			remarks = strings.SplitN(serverInfo[1], "|", 2)[1]
		} else {
			remarks, err = url.QueryUnescape(serverInfo[1])
			if err != nil {
				return model.Proxy{}, err
			}
		}
	} else {
		remarks, err = url.QueryUnescape(serverAndPort[0])
		if err != nil {
			return model.Proxy{}, err
		}
	}
	server := strings.TrimSpace(serverAndPort[0])
	uuid := strings.TrimSpace(parts[0])
	network := params.Get("type")
	result := model.Proxy{
		Type: "vless",
		VLESS: model.VLESS{
			Tag:        remarks,
			Server:     server,
			ServerPort: uint16(port),
			UUID:       uuid,
			Network:    network,
			Flow:       params.Get("flow"),
		},
	}
	if params.Get("security") == "tls" {
		result.VLESS.TLS = &model.OutboundTLSOptions{
			Enabled:  true,
			ALPN:     strings.Split(params.Get("alpn"), ","),
			Insecure: params.Get("allowInsecure") == "1",
		}
	}
	if params.Get("security") == "reality" {
		result.VLESS.TLS = &model.OutboundTLSOptions{
			Enabled:    true,
			ServerName: params.Get("sni"),
			UTLS: &model.OutboundUTLSOptions{
				Enabled:     params.Get("fp") != "",
				Fingerprint: params.Get("fp"),
			},
			Reality: &model.OutboundRealityOptions{
				Enabled:   true,
				PublicKey: params.Get("pbk"),
				ShortID:   params.Get("sid"),
			},
			ALPN: strings.Split(params.Get("alpn"), ","),
		}
	}
	if params.Get("type") == "ws" {
		result.VLESS.Transport = &model.V2RayTransportOptions{
			Type: "ws",
			WebsocketOptions: model.V2RayWebsocketOptions{
				Path: params.Get("path"),
				Headers: map[string]string{
					"Host": params.Get("host"),
				},
			},
		}
	}
	if params.Get("type") == "quic" {
		result.VLESS.Transport = &model.V2RayTransportOptions{
			Type:        "quic",
			QUICOptions: model.V2RayQUICOptions{},
		}
	}
	if params.Get("type") == "grpc" {
		result.VLESS.Transport = &model.V2RayTransportOptions{
			Type: "grpc",
			GRPCOptions: model.V2RayGRPCOptions{
				ServiceName: params.Get("serviceName"),
			},
		}
	}
	if params.Get("type") == "http" {
		host, err := url.QueryUnescape(params.Get("host"))
		if err != nil {
			return model.Proxy{}, err
		}
		result.VLESS.Transport = &model.V2RayTransportOptions{
			Type: "http",
			HTTPOptions: model.V2RayHTTPOptions{
				Host: strings.Split(host, ","),
			},
		}
	}
	return result, nil
}
