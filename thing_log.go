// Copyright 2020 thinkgos (thinkgo@aliyun.com).  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package aiot

import (
	"encoding/json"

	"github.com/things-go/aliyun-iot/infra"
	"github.com/things-go/aliyun-iot/uri"
)

// @see https://help.aliyun.com/document_detail/159857.html?spm=a2c4g.11186623.6.714.39cb741aXw9Osf

// Log level,从高到低
const (
	LogFatal = "FATAL"
	LogError = "ERROR"
	LogWarn  = "WARN"
	LogInfo  = "INFO"
	LogDebug = "DEBUG"
	LogOther = "OTHER"
)

// ConfigLogParam 设备获取日志配置参数域
type ConfigLogParam struct {
	// 配置范围,目前日志只有设备维度配置,默认为device
	ConfigScope string `json:"configScope"`
	// 获取内容类型,默认为content.因日志配置内容较少,默认直接返回内容
	GetType string `json:"getType"`
}

// ThingConfigLogGet 获取日志配置
// request： /sys/${productKey}/${deviceName}/thing/config/Log/get
// response：/sys/${productKey}/${deviceName}/thing/config/Log/get_reply
func (sf *Client) ThingConfigLogGet(pk, dn string, _ ConfigLogParam) (*Token, error) {
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingConfigLogGet, pk, dn)
	return sf.SendRequest(_uri, infra.MethodConfigLogGet, ConfigLogParam{
		"device",
		"content",
	})
}

// LogParam 日志内容参数域
type LogParam struct {
	// 日志的采集时间,为设备本地UTC时间,
	// 包含时区信息,以毫秒计,格式为 yyyy-MM-dd'T'HH:mm:ss.SSSZ
	UtcTime string `json:"utcTime"`
	// FATAL,ERROR,WARN,INFO,DEBUG
	LogLevel string `json:"logLevel"`
	// 当设备端使用Android SDK时,模块名称为ALK-LK.
	// 当设备端使用C SDK时,需自定义模块名称.
	// 当设备端使用自行开发的SDK时,需自定义模块名称
	Module string `json:"module"`
	// 当设备端使用Android SDK时,请参见Android SDK错误码.
	// 当设备端使用C SDK时,请参见C SDK状态码.
	// 当设备端使用自行开发的SDK时,可以自定义结果状态码,也可以为空.
	Code string `json:"code"`
	// 可选参数,上下文跟踪内容,设备端使用Alink协议消息的id,App端使用TraceId(追踪ID)
	TraceContext string `json:"traceContext,omitempty"`
	// 日志详细内容
	LogContent string `json:"logContent"`
}

// ThingLogPost 设备上报日志内容
// request： /sys/${productKey}/${deviceName}/thing/config/Log/post
// response：/sys/${productKey}/${deviceName}/thing/config/Log/post_reply
func (sf *Client) ThingLogPost(pk, dn string, lp []LogParam) (*Token, error) {
	if !sf.IsActive(pk, dn) {
		return nil, ErrNotActive
	}
	if len(lp) == 0 {
		return nil, ErrInvalidParameter
	}
	_uri := uri.URI(uri.SysPrefix, uri.ThingLogPost, pk, dn)
	return sf.SendRequest(_uri, infra.MethodLogPost, lp)
}

// ConfigLogMode 日志配置的日志上报模式
type ConfigLogMode struct {
	Mode int `json:"mode"` // 0 表示不上报, 1 表示上报
}

// ConfigLogParamData 日志配置的参数域或配置域
type ConfigLogParamData struct {
	GetType string        `json:"getType"`
	Content ConfigLogMode `json:"content"`
}

// ConfigLogResponse 日志配置回复
type ConfigLogResponse struct {
	ID      uint               `json:"id,string"`
	Code    int                `json:"code"`
	Data    ConfigLogParamData `json:"data"`
	Message string             `json:"message,omitempty"`
}

// ProcThingConfigLogGetReply 处理获取日志配置应答
// request：  /sys/${productKey}/${deviceName}/thing/config/Log/get
// response： /sys/${productKey}/${deviceName}/thing/config/Log/get_reply
// subscribe: /sys/${productKey}/${deviceName}/thing/config/Log/get_reply
func ProcThingConfigLogGetReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}
	rsp := &ConfigLogResponse{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, rsp.Data, err})
	c.Log.Debugf("thing.config.log.get.reply @%d", rsp.ID)
	pk, dn := uris[1], uris[2]
	return c.cb.ThingConfigLogGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingLogPostReply 处理日志上报应答
// request： /sys/${productKey}/${deviceName}/thing/Log/post
// response：/sys/${productKey}/${deviceName}/thing/Log/post_reply
func ProcThingLogPostReply(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 6 {
		return ErrInvalidURI
	}
	rsp := &Response{}
	err := json.Unmarshal(payload, rsp)
	if err != nil {
		return err
	}
	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}
	c.signalPending(Message{rsp.ID, nil, err})
	c.Log.Debugf("thing.log.post.reply @%d", rsp.ID)
	pk, dn := uris[1], uris[2]
	return c.cb.ThingLogPostReply(c, err, pk, dn)
}

// ConfigLogPush 日志配置推送
type ConfigLogPush struct {
	ID      uint               `json:"id,string"`
	Version string             `json:"version"`
	Params  ConfigLogParamData `json:"params"`
}

// ProcThingConfigLogPush 处理日志配置推送
// subscribe: /sys/${productKey}/${deviceName}/thing/config/Log/push
func ProcThingConfigLogPush(c *Client, rawURI string, payload []byte) error {
	uris := uri.Spilt(rawURI)
	if len(uris) < 7 {
		return ErrInvalidURI
	}

	req := &ConfigLogPush{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return err
	}

	c.Log.Debugf("thing.config.log.push @%d", req.ID)
	pk, dn := uris[1], uris[2]
	return c.cb.ThingConfigLogPush(c, pk, dn, req.Params)
}
