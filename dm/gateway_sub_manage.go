package dm

import (
	"encoding/json"

	"github.com/thinkgos/aliyun-iot/infra"
)

// ProcThingGwDisable 禁用子设备
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/disable
// response:  /sys/{productKey}/{deviceName}/thing/disable_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/disable
func ProcThingGwDisable(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 5) {
		return ErrInvalidURI
	}

	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream <thing>: disable >> %s - %s", pk, dn)

	req := &Request{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return err
	}

	if err = c.SetDevAvailByPkDN(pk, dn, false); err != nil {
		c.warnf("<thing> disable failed, %+v", err)
	}
	if err = c.eventGwProc.EvtThingDisable(c, pk, dn); err != nil {
		c.warnf("<thing> disable, ipc send message failed, %+v", err)
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}

// ProcThingGwEnable 启用子设备
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/enable
// response:  /sys/{productKey}/{deviceName}/thing/enable_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/enable
func ProcThingGwEnable(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 5) {
		return ErrInvalidURI
	}
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream <thing>: enable >> %s - %s", pk, dn)

	req := &Request{}
	err := json.Unmarshal(payload, req)
	if err != nil {
		return err
	}

	if err = c.SetDevAvailByPkDN(pk, dn, true); err != nil {
		c.warnf("<thing> enable failed, %+v", err)
	}

	if err = c.eventGwProc.EvtThingEnable(c, pk, dn); err != nil {
		c.warnf("<thing> enable, ipc send message failed, %+v", err)
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}

// ProcThingGwDelete 子设备删除,网关类型设备
// 下行
// request:   /sys/{productKey}/{deviceName}/thing/delete
// response:  /sys/{productKey}/{deviceName}/thing/delete_reply
// subscribe: /sys/{productKey}/{deviceName}/thing/delete
func ProcThingGwDelete(c *Client, rawURI string, payload []byte) error {
	uris := URIServiceSpilt(rawURI)
	if len(uris) < (c.uriOffset + 5) {
		return ErrInvalidURI
	}
	pk, dn := uris[c.uriOffset+1], uris[c.uriOffset+2]
	c.debugf("downstream <thing>: delete >> %s - %s", pk, dn)

	req := &Request{}
	if err := json.Unmarshal(payload, req); err != nil {
		return err
	}

	c.DeleteByPkDn(pk, dn)
	if err := c.eventGwProc.EvtThingDelete(c, pk, dn); err != nil {
		c.warnf("<thing> delete, ipc send message failed, %+v", err)
	}
	return c.SendResponse(URIServiceReplyWithRequestURI(rawURI), req.ID, infra.CodeSuccess, "{}")
}