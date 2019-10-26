package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/aliIOT"
	"github.com/thinkgos/aliIOT/infra"
	"github.com/thinkgos/aliIOT/model"
	"github.com/thinkgos/aliIOT/sign"
)

const (
	productKey    = "a1iJcssSlPC"
	productSecret = "lw3QzKHNfh7XvOxO"
	deviceName    = "dyncreg"
	deviceSecret  = "irqurH8zaIg1ChoeaBjLHiqBXEZnlVq8"
)

func main() {
	bPayload := []byte{0x0, 0x6, 0x0, 0x0, 0x0, 0x56, 0x32, 0x32, 0x32, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x47, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x4d, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x41, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x0, 0x1, 0x0, 0x2, 0x0, 0x3, 0x0, 0x4, 0x0, 0x5, 0x0, 0x6, 0x0, 0x7, 0x0, 0x8, 0x0, 0xf, 0x42, 0x40, 0x1, 0x2, 0x3, 0xc, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x1, 0x2d, 0x1, 0x2e, 0x1, 0x2f, 0x1, 0x30, 0x1, 0x31, 0x1, 0x32, 0x1, 0x33, 0x1, 0x34, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64, 0x1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xb, 0xc, 0xd, 0x3, 0xe9, 0x3, 0xea, 0x3, 0xeb, 0xd, 0x6e, 0xb, 0xc, 0xd, 0xe, 0xf, 0x3, 0xeb, 0x3, 0xeb, 0x10, 0x11, 0x12, 0x13, 0x14, 0xfc, 0x16, 0xfc, 0x15, 0x10, 0x11, 0x12, 0x13, 0x14, 0x10, 0x11, 0x12, 0x13, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64, 0xe, 0xf, 0x2, 0x0, 0xb, 0xc, 0xd, 0xe, 0xf, 0x3, 0xe8, 0x10, 0x11, 0x12, 0x13, 0x14, 0x0, 0x10, 0x11, 0x12, 0x13, 0x14, 0xa, 0x10, 0x11, 0x12, 0x13, 0x14, 0x0, 0x1, 0x2, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x3, 0xe6, 0x3, 0xe7, 0x60, 0x61, 0x62, 0x63, 0x64, 0x3, 0xe7, 0x60, 0x61, 0x62, 0x63, 0x64, 0x0, 0x0, 0x0, 0x64, 0x7, 0x0, 0x0, 0x0, 0x64, 0x0, 0x64, 0x1, 0x91, 0x1, 0x91, 0x1, 0x91, 0x1, 0x91, 0x1, 0x91, 0x1, 0x91, 0x1, 0x91, 0x1, 0x91, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64}

	signs, err := sign.NewMQTTSign().SetSDKVersion(infra.IOTSDKVersion).Generate(&sign.MetaInfo{
		ProductKey:    productKey,
		ProductSecret: productSecret,
		DeviceName:    deviceName,
		DeviceSecret:  deviceSecret,
	}, sign.CloudRegionShangHai)
	if err != nil {
		panic(err)
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", signs.HostName, signs.Port))
	opts.SetClientID(signs.ClientID)
	opts.SetUsername(signs.UserName)
	opts.SetPassword(signs.Password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		log.Println("mqtt client connection success")
	})
	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		log.Println("mqtt client connection lost, ", err)
	})
	client := mqtt.NewClient(opts)
	dmopt := model.NewOption(productKey, deviceName, deviceSecret).Valid()
	manage := aliIOT.NewWithMQTT(dmopt, client)

	client.Connect().Wait()
	manage.LogMode(true)
	_ = manage.Subscribe(manage.URIServiceItself(model.URISysPrefix, model.URIThingModelUpRawReply), model.ProcThingModelUpRawReply)

	//b, _ := hex.DecodeString(payload)
	//log.Printf("%#v", b)
	for {
		err = manage.UpstreamThingModelUpRaw(model.DevItself, bPayload)
		if err != nil {
			log.Printf("error: %#v", err)
		} else {
			log.Printf("success")
		}
		time.Sleep(time.Second * 10)
	}
}
