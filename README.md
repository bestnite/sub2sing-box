# sub2sing-box

将订阅/节点连接转换为 sing-box 配置的工具。

## 控制台命令

使用 `sub2sing-box <command> -h` 查看各命令的帮助信息。

## 配置

示例:

```json
{
  "subscriptions": ["订阅地址1", "订阅地址2"],
  "proxies": ["代理1", "代理2"],
  "template": "模板路径或网络地址",
  "delete": "剩余流量",
  "rename": { "原文本": "新文本" },
  "group-type": "selector",
  "sort": "tag",
  "sort-type": "asc",
  "output": "./config.json",
  "user-agent": "自定义 User-Agent"
}
```

将上述 JSON 内容保存为 `sub2sing-box.json`，执行:

```
sub2sing-box convert -c ./sub2sing-box.json
```

即可生成 sing-box 配置，无需每次重复设置参数。

补充说明：

- `subscriptions` / `proxies` 必须写成数组，分别对应命令行里的 `-s` 与 `-p`。
- `sort` 只接受 `tag`（按分组名排序）或 `num`（按分组内节点数量排序），其他值会被当作 `tag`。
- 配置键名是严格的：拼错或写了未知的键会直接报错，不会被静默忽略（`/convert` 接口的 `data` 同样适用）。

## 模板

### 默认模板

默认模板位于 `templates` 目录，使用 `tun` 配置，该配置仅供参考，请根据实际情况修改。

### 占位符

模板中可使用以下占位符:

- `<all-proxy-tags>`: 插入所有节点标签
- `<all-country-tags>`: 插入所有国家标签
- `<国家(地区)二字码>`: 插入指定国家(地区)所有节点标签，如 `<tw>`

#### 占位符使用示例:

假设有节点：

- US
- SG
- TW

```json
{
  "type": "selector",
  "tag": "节点选择",
  "outbounds": ["<all-proxy-tags>", "direct"],
  "interrupt_exist_connections": true
}

// 转换后

{
  "type": "selector",
  "tag": "节点选择",
  "outbounds": ["US", "SG", "TW", "direct"],
  "interrupt_exist_connections": true
}
```

```json
{
  "type": "selector",
  "tag": "节点选择",
  "outbounds": ["<all-country-tags>", "direct"],
  "interrupt_exist_connections": true
}

// 转换后

{
  "type": "selector",
  "tag": "节点选择",
  "outbounds": ["美国(US)", "新加坡(SG)", "台湾(TW)", "direct"],
  "interrupt_exist_connections": true
}

// 其中 "美国(US)", "新加坡(SG)", "台湾(TW)" 为策略组，分别包含 US, SG, TW 节点
```

```json
{
  "type": "selector",
  "tag": "巴哈姆特",
  "outbounds": ["<tw>", "direct"],
  "interrupt_exist_connections": true
}

// 转换后

{
  "type": "selector",
  "tag": "巴哈姆特",
  "outbounds": ["台湾(TW)", "direct"],
  "interrupt_exist_connections": true
}
```

## 支持的协议

`-p/--proxy` 与订阅（`-s/--subscription`）里的分享链接支持以下协议：

| 协议             | 链接格式                                                                            | 缺省端口 |
| ---------------- | ----------------------------------------------------------------------------------- | -------- |
| Shadowsocks      | `ss://`（SIP002 / SIP008）                                                          | 链接必填 |
| VMess            | `vmess://`                                                                          | 链接必填 |
| VLESS            | `vless://`                                                                          | 链接必填 |
| Trojan           | `trojan://`                                                                         | 链接必填 |
| Hysteria         | `hysteria://`                                                                       | 链接必填 |
| Hysteria2        | `hysteria2://`、`hy2://`                                                            | 链接必填 |
| TUIC             | `tuic://`                                                                           | 链接必填 |
| AnyTLS           | `anytls://`                                                                         | 链接必填 |
| SOCKS            | `socks://`、`socks5://`                                                             | 链接必填 |
| SSH              | `ssh://user:password@host:port?private_key=...&host_key=...#标签`                    | 22       |
| NaïveProxy       | `naive+https://`、`naive+quic://`（可选 `extra-headers`、`sni`）                     | 443      |
| HTTP(S) 代理     | `proxy-http://`、`proxy-https://`（可选 `path`、`headers`、`sni`）                   | 80 / 443 |

多个链接要写成重复的参数，不要用空格分隔（空格后面那些会被当成位置参数直接忽略）：

```
sub2sing-box convert -t sing-box-template.json -p "链接1" -p "链接2" -p "链接3" -o config.json
```

链接里如果含逗号（例如 ssh 的 `host_key_algorithms=ssh-ed25519,rsa-sha2-256`），这个逗号会被命令行框架当成参数分隔符拆开，此时请改用配置文件（`-c sub2sing-box.json`）里的 `proxies` 数组：

```json
{
  "template": "sing-box-template.json",
  "proxies": ["ssh://root:password@host:22?host_key_algorithms=ssh-ed25519,rsa-sha2-256#SSH"],
  "output": "config.json"
}
```

以下链接没有可用出站，转换时会明确说明原因：`ssr://`（sing-box 1.6.0 移除）、`wireguard://` / `wg://`（1.13.0 起改用 endpoint）、`juicity://`（内核无此出站）、`snell://`（无通用链接格式）。

## Docker 使用

```
docker run -p 8080:8080 nite07/sub2sing-box:latest
```

可以挂载目录添加自定义模板

## Server 模式 API

### GET /convert

| query | 描述                                                                                                                    |
| ----- | ----------------------------------------------------------------------------------------------------------------------- |
| data  | 同上方配置，但需要使用 [base64 URL safe 编码](<https://gchq.github.io/CyberChef/#recipe=To_Base64('A-Za-z0-9%2B/%3D')>) |
