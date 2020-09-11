package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
)

// ConfigGetParams 获取配置的参数域
type ConfigGetParams struct {
	ConfigScope string `json:"configScope"` // 配置范围, 目前只支持产品维度配置. 取值: product
	GetType     string `json:"getType"`     // 获取配置类型. 目前支持文件类型,取值: file
}

// ConfigParamsData 配置获取回复数据域或配置推送参数域
type ConfigParamsData struct {
	ConfigID   string `json:"configId"`   // 配置文件的ID
	ConfigSize int64  `json:"configSize"` // 配置文件大小,按字节计算
	Sign       string `json:"sign"`       // 签名
	SignMethod string `json:"signMethod"` // 签名方法,仅支持Sha256
	URL        string `json:"url"`        // 存储配置文件的对象存储（OSS）地址
	GetType    string `json:"getType"`    // 同ConfigGetParams.GetType
}

// ConfigGetResponse 配置获取的回复
type ConfigGetResponse struct {
	ID      uint             `json:"id,string"`
	Code    int              `json:"code"`
	Data    ConfigParamsData `json:"data"`
	Message string           `json:"message,omitempty"`
}

// ConfigPushRequest 配置推送的请求
type ConfigPushRequest struct {
	ID      uint             `json:"id,string"`
	Version string           `json:"version"`
	Params  ConfigParamsData `json:"params"`
	Method  string           `json:"method"`
}

// ThingConfigGet 获取配置参数
// request:  /sys/{productKey}/{deviceName}/thing/config/get
// response: /sys/{productKey}/{deviceName}/thing/config/get_reply
func (sf *Client) ThingConfigGet(devID int) (*Entry, error) {
	if !sf.isGateway {
		return nil, ErrNotSupportFeature
	}
	if devID < 0 {
		return nil, ErrInvalidParameter
	}
	node, err := sf.SearchNode(devID)
	if err != nil {
		return nil, err
	}

	uri := sf.URIService(URISysPrefix, URIThingConfigGet, node.ProductKey(), node.DeviceName())
	id := sf.RequestID()
	err = sf.SendRequest(uri, id, MethodConfigGet, ConfigGetParams{"product", "file"})
	if err != nil {
		return nil, err
	}
	sf.debugf("upstream thing <config>: get,@%d", id)
	return sf.Insert(id), nil
}

// ProcThingConfigGetReply 处理获取配置的应答
// 上行
// request:   /sys/{productKey}/{deviceName}/thing/config/get
// response:  /sys/{productKey}/{deviceName}/thing/config/get_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/config/get_reply
func ProcThingConfigGetReply(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}

	rsp := ConfigGetResponse{}
	err := json.Unmarshal(payload, &rsp)
	if err != nil {
		return err
	}

	if rsp.Code != infra.CodeSuccess {
		err = infra.NewCodeError(rsp.Code, rsp.Message)
	}

	c.done(rsp.ID, err, nil)
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream thing <config>: get reply,@%d,payload@%+v", rsp.ID, rsp)
	return c.eventProc.EvtThingConfigGetReply(c, err, pk, dn, rsp.Data)
}

// ProcThingConfigPush 处理配置推送
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/config/push
// response:  /sys/{productKey}/{deviceName}/thing/config/push_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/config/push
func ProcThingConfigPush(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 6) {
		return ErrInvalidURI
	}
	req := ConfigPushRequest{}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	err := c.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
	if err != nil {
		return err
	}
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream thing <config>: push request")
	return c.eventProc.EvtThingConfigPush(c, pk, dn, req.Params)
}