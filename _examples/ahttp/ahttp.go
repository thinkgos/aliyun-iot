package main

import (
	"fmt"
	"log"
	"time"

	"github.com/thinkgos/aliyun-iot/ahttp"
	"github.com/thinkgos/aliyun-iot/dm"
)

const (
	productKey    = "a1iJcssSlPC"
	productSecret = "lw3QzKHNfh7XvOxO"
	deviceName    = "rawtest"
	deviceSecret  = "ld9Xf2BtKGfdEC7G9nSMe1wYfgllvi3Q"
)

// 透传
func main() {
	var err error

	client := ahttp.New().SetDeviceMetaInfo(productKey, deviceName, deviceSecret)
	client.
		LogMode(true)

	uri := dm.URICOAPHTTPPrePrefix + fmt.Sprintf(dm.URISysPrefix, productKey, deviceName) + dm.URIThingModelUpRaw
	bPayload := []byte{0x00, 0x00, 0x00, 0x04, 0xeb, 0x04, 0x76, 0x73, 0x6e, 0x30,
		0x04, 0x67, 0x73, 0x6e, 0x30, 0x04, 0x6d, 0x73, 0x6e, 0x30, 0x05, 0x61,
		0x63, 0x73, 0x6e, 0x30, 0x00, 0x41, 0x00, 0x41, 0x00, 0x41, 0x00, 0x41,
		0x00, 0x41, 0x00, 0x41, 0x00, 0x41, 0x00, 0x41, 0x00, 0x00, 0x04, 0xeb,
		0x00, 0x17, 0x84, 0x7a, 0x09, 0xc2, 0x5a, 0x06, 0x84, 0x47, 0x37, 0x6e,
		0x57, 0x53, 0x25, 0x18, 0x19, 0x40, 0x28, 0x3b, 0x21, 0x24, 0x01, 0xbb,
		0x00, 0x11, 0x00, 0xd9, 0x01, 0x9b, 0x01, 0x46, 0x01, 0x31, 0x01, 0xa1,
		0x00, 0xca, 0x59, 0x05, 0x4a, 0x51, 0x17, 0x16, 0x0a, 0x2e, 0x1e, 0x19,
		0x01, 0x98, 0x24, 0x3e, 0x9a, 0x5c, 0x27, 0x08, 0x24, 0x27, 0x0f, 0xcb,
		0x0a, 0x55, 0x1b, 0x76, 0xfe, 0x7b, 0x53, 0x57, 0x55, 0x37, 0x11, 0x12,
		0x5f, 0x42, 0x18, 0x4d, 0x4c, 0x1d, 0x16, 0x55, 0xf2, 0xb9, 0x21, 0x67,
		0x60, 0x27, 0x08, 0x4d, 0x4c, 0x74, 0xa8, 0x5e, 0x01, 0x08, 0x5c, 0x08,
		0x4b, 0x63, 0x55, 0x18, 0x5e, 0x0e, 0x1f, 0x97, 0x99, 0x00, 0x80, 0x61,
		0x1f, 0x4c, 0x24, 0x4f, 0x00, 0x04, 0x57, 0x3d, 0x2a, 0x49, 0x2f, 0xc1,
		0x06, 0x15, 0x2a, 0x50, 0x2d, 0x1f, 0x1a, 0x38, 0x14, 0x0c, 0x59, 0x58,
		0x5a, 0x01, 0x60, 0x26, 0x0d, 0x14, 0x31, 0x01, 0xe1, 0x01, 0x3b, 0x3a,
		0x20, 0x5b, 0x01, 0x45, 0x02, 0x1a, 0x57, 0x23, 0x62, 0x0f, 0x17, 0x00,
		0x00, 0x04, 0xeb, 0x05, 0xe7, 0x6a, 0x7b, 0xcb, 0xd5, 0x17, 0x01, 0x56,
		0x00, 0x92, 0x00, 0xb0, 0x01, 0xf9, 0x01, 0x03, 0x00, 0x23, 0x00, 0x23,
		0x00, 0x43, 0x5b, 0x7c, 0x80, 0x42, 0xa2, 0x2e, 0x56, 0x68, 0x5c, 0x41,
		0x53, 0x51, 0x09, 0x51, 0x2b, 0x4e, 0x61, 0x43, 0x09, 0x2a, 0x14, 0x4d,
		0x42, 0x1f, 0x47, 0x38, 0x52, 0x47}
	for {
		err = client.Publish(uri, bPayload)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 10)
	}
}