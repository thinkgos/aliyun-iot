package dm

import (
	"fmt"
	"strings"
)

// 分隔符定义
const (
	SEP = "/"
)

// URI 前缀定义
const (
	URICOAPHTTPPrePrefix = "/topic"

	URISysPrefix             = "/sys/%s/%s/"
	URIExtSessionPrefix      = "/ext/session/%s/%s/"
	URIExtNtpPrefix          = "/ext/ntp/%s/%s/"
	URIExtErrorPrefix        = "/ext/error/%s/%s"
	URIOtaDeviceInformPrefix = "/ota/device/inform/%s/%s"

	URIReplySuffix = "reply"
)

// RRPC URI定义
const (
	//  系统RRPC调用
	URIRRPCRequestSingleWildcard = "rrpc/request/+"
	URIRRPCResponse              = "rrpc/response/%s"

	// 自定义RRPC
	URIExtRRPCPrefix         = "/ext/rrpc/%s/"
	URIExtRRPCSingleWildcard = "/ext/rrpc/+/"
)

// 设备URI 定义
const (
	// 透传数据上行,下行云端
	URIThingModelDownRaw      = "thing/model/down_raw"
	URIThingModelDownRawReply = "thing/model/down_raw_reply"
	URIThingModelUpRaw        = "thing/model/up_raw"
	URIThingModelUpRawReply   = "thing/model/up_raw_reply"

	// 事件上行,下行云端
	URIThingEventPropertyPost            = "thing/event/property/post"
	URIThingEventPropertyPostReply       = "thing/event/property/post_reply"
	URIThingEventPost                    = "thing/event/%s/post"
	URIThingEventPostReply               = "thing/event/%s/post_reply"
	URIThingEventPostReplySingleWildcard = "thing/event/+/post_reply"

	// 设备信息上行,下行云端
	URIThingDeviceInfoUpdate      = "thing/deviceinfo/update"
	URIThingDeviceInfoUpdateReply = "thing/deviceinfo/update_reply"
	URIThingDeviceInfoDelete      = "thing/deviceinfo/delete"
	URIThingDeviceInfoDeleteReply = "thing/deviceinfo/delete_reply"

	// 期望属性值上行,下行云端
	URIThingDesiredPropertyGet         = "thing/property/desired/get"
	URIThingDesiredPropertyGetReply    = "thing/property/desired/get_reply"
	URIThingDesiredPropertyDelete      = "thing/property/desired/delete"
	URIThingDesiredPropertyDeleteReply = "thing/property/desired/delete_reply"

	// 服务调用上行,下行云端
	URIThingServicePropertySet           = "thing/service/property/set"
	URIThingServicePropertySetReply      = "thing/service/property/set_reply"
	URIThingServiceRequest               = "thing/service/%s"
	URIThingServiceResponse              = "thing/service/%s_reply"
	URIThingServiceRequestSingleWildcard = "thing/service/+"
	URIThingServiceRequestMultiWildcard2 = "thing/service/#"

	/* dsl template From Local To Cloud Request And Response */
	URIThingDslTemplateGet      = "thing/dsltemplate/get"
	URIThingDslTemplateGetReply = "thing/dsltemplate/get_reply"

	/* dynamic tsl From Local To Cloud Request And Response */
	URIThingDynamicTslGet      = "thing/dynamicTsl/get"
	URIThingDynamicTslGetReply = "thing/dynamicTsl/get_reply"

	/* ntp From Local To Cloud Request And Response */
	URINtpRequest  = "request"
	URINtpResponse = "response"

	// config
	URIThingConfigGet       = "thing/config/get"
	URIThingConfigGetReply  = "thing/config/get_reply"
	URIThingConfigPush      = "thing/config/push"
	URIThingConfigPushReply = "thing/config/push_reply"
)

// 网关URI定义
const (
	// 子设备动态注册
	URIThingSubDevRegister      = "thing/sub/register"
	URIThingSubDevRegisterReply = "thing/sub/register_reply"

	// 子设备登录
	URISubDevCombineLogin      = "combine/login"
	URISubDevCombineLoginReply = "combine/login_reply"
	URISubDevCombineLogout     = "combine/logout"
	SubDevCombineLogoutReply   = "combine/logout_reply"

	// 网关网络拓扑
	URIThingTopoAdd         = "thing/topo/add"
	URIThingTopoAddReply    = "thing/topo/add_reply"
	URIThingTopoDelete      = "thing/topo/delete"
	URIThingTopoDeleteReply = "thing/topo/delete_reply"
	URIThingTopoGet         = "thing/topo/get"
	URIThingTopoGetReply    = "thing/topo/get_reply"
)

// URIServiceSelf 生成URI, inName的作用是有一些需要格式化到name里.
func URIService(prefix, name, productKey, deviceName string, inName ...string) string {
	str := strings.Builder{}
	str.Grow(len(prefix) + len(name) + len(productKey) + len(deviceName))
	if prefix != "" {
		str.WriteString(fmt.Sprintf(prefix, productKey, deviceName))
	}
	if name != "" {
		if len(inName) > 0 {
			str.WriteString(fmt.Sprintf(name, inName[0]))
		} else {
			str.WriteString(name)
		}
	}
	return str.String()
}

// URIServiceReplyWithRequestURI 根据请求的URI生成应答的URI
func URIServiceReplyWithRequestURI(uri string) string {
	return uri + "_" + URIReplySuffix
}

// URIServiceSpilt 分隔URI,会去除左边的SEP分隔符
func URIServiceSpilt(uri string) []string {
	return strings.Split(strings.TrimLeft(uri, SEP), SEP)
}
