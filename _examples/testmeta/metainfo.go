package testmeta

import (
	"github.com/thinkgos/aliyun-iot/infra"
)

// just for test
const (
	ProductKey    = "a1QR3GD1Db3"
	ProductSecret = "mvngTYBlX9Z9l1V0"
	DeviceName    = "dynamic"
	DeviceSecret  = "9690f9da431078f105b7969b23e05762"
)

func MetaInfo() infra.MetaInfo {
	return infra.MetaInfo{
		ProductKey,
		ProductSecret,
		DeviceName,
		DeviceSecret,
	}
}
