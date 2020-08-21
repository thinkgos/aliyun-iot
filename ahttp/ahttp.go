// Package ahttp 实现http client 上传数据. 授权方式为自动调用授权,可手动调用,也可以直接调用发送数据接口
package ahttp

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thinkgos/aliyun-iot/clog"
)

const (
	signMethodHMACSHA1   = "hmacsha1"
	signMethodHMACMD5    = "hmacmd5"
	defaultAuthLimitTime = time.Minute * 15 // 当授权通过后,在15分钟内不可再授权,防止授权频繁
)

// AuthRequest 鉴权请求
type AuthRequest struct {
	Version    string `json:"version"`
	ClientID   string `json:"clientId"`
	SignMethod string `json:"signmethod"`
	Sign       string `json:"sign"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
	// 校验时间戳15分钟内的请求有效。时间戳格式为数值，
	// 值为自GMT 1970年1月1日0时0分到当前时间点所经过的毫秒数。
	Timestamp int64 `json:"timestamp"`
}

// AuthResponse 鉴权回复
type AuthResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Info    struct {
		Token string `json:"token"`
	} `json:"info"`
}

// Client 客户端
type Client struct {
	productKey   string
	deviceName   string
	deviceSecret string
	host         string
	version      string
	signMethod   string

	token    atomic.Value
	whenAuth time.Time
	mu       sync.Mutex

	httpc *http.Client
	*clog.Clog
}

// Option client option
type Option func(c *Client)

// WithHTTPClient with custom http.Client
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpc = c
	}
}

// WithHost 设置远程主机,
func WithHost(h string) Option {
	return func(c *Client) {
		if !strings.Contains(h, "://") {
			h = "http://" + h
		}

		if h != "" {
			c.host = h
		}
	}
}

// WithSignMethod 设置签名方法,目前支持hmacMD5和hmacSHA1
func WithSignMethod(method string) Option {
	return func(c *Client) {
		if method == signMethodHMACMD5 || method == signMethodHMACSHA1 {
			c.signMethod = method
		} else {
			c.signMethod = signMethodHMACMD5
		}
	}
}

// New 新建alink http client
// 默认hmacmd5加签算法
// 默认上海host
// 请求超时2秒
func New(opts ...Option) *Client {
	c := &Client{
		host:       "https://iot-as-http.cn-shanghai.aliyuncs.com",
		version:    "default",
		signMethod: signMethodHMACMD5,
		httpc:      http.DefaultClient,
		Clog:       clog.New(clog.WithLogger(clog.NewLogger(log.New(os.Stderr, "alink http --> ", log.LstdFlags)))),
	}
	c.token.Store("")
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// SetDeviceMetaInfo 设置设备三元组信息
func (sf *Client) SetDeviceMetaInfo(productKey, deviceName, deviceSecret string) *Client {
	sf.productKey = productKey
	sf.deviceName = deviceName
	sf.deviceSecret = deviceSecret
	return sf
}

// sendAuth 鉴权
// TODO: use cache
func (sf *Client) sendAuth() error {
	if sf.productKey == "" ||
		sf.deviceName == "" ||
		sf.deviceSecret == "" {
		return errors.New("invalid device meta info")
	}

	sf.mu.Lock()
	defer sf.mu.Unlock()
	// 如果刚在15分钟内刚授权过,不用再授权了. 直接返回
	if time.Since(sf.whenAuth) < defaultAuthLimitTime {
		return nil
	}
	authPy := AuthRequest{
		Version:    sf.version,
		ClientID:   sf.productKey + "." + sf.deviceName,
		SignMethod: sf.signMethod,
		ProductKey: sf.productKey,
		DeviceName: sf.deviceName,
		Timestamp:  time.Now().Unix() * 1000,
	}

	if err := authPy.generateSign(sf.deviceSecret); err != nil {
		return err
	}

	b, err := json.Marshal(&authPy)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, sf.host+"/auth", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := sf.httpc.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	rspPy := AuthResponse{}
	if err = json.NewDecoder(response.Body).Decode(&rspPy); err != nil {
		return err
	}

	if rspPy.Code != CodeSuccess {
		err = NewCodeError(rspPy.Code, rspPy.Message)
		sf.Debug("auth failed, %+v", err)
		return err
	}
	sf.token.Store(rspPy.Info.Token)
	sf.whenAuth = time.Now()
	sf.Debug("auth success!")
	return nil
}

// DataResponse 上报数据回复
type DataResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Info    struct {
		MessageID int64 `json:"messageID"`
	} `json:"info"`
}

func (sf *Client) publish(uri string, payload interface{}) (int64, error) {
	token := sf.token.Load().(string)
	if token == "" {
		return 0, NewCodeError(CodeTokenIsNull, "token is null")
	}

	var buf *bytes.Buffer
	switch v := payload.(type) {
	case string:
		buf = bytes.NewBufferString(v)
	case []byte:
		buf = bytes.NewBuffer(v)
	default:
		return 0, errors.New("Unknown payload type, must be string or []byte")
	}

	request, err := http.NewRequest(http.MethodPost, sf.host+uri, buf)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("password", token)
	response, err := sf.httpc.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	rspPy := DataResponse{}
	if err = json.NewDecoder(response.Body).Decode(&rspPy); err != nil {
		return 0, err
	}
	sf.Debug("publish response, %+v", rspPy)
	if rspPy.Code == 0 {
		return rspPy.Info.MessageID, nil
	}
	return 0, NewCodeError(rspPy.Code, rspPy.Message)
}

// Publish 数据推送
func (sf *Client) Publish(uri string, payload interface{}) error {
	_, err := sf.publish(uri, payload)
	if err != nil {
		var pErr *CodeError
		if errors.As(err, &pErr) &&
			(pErr.Code() == CodeTokenExpired ||
				pErr.Code() == CodeTokenCheckFailed ||
				pErr.Code() == CodeTokenIsNull) {
			if err = sf.sendAuth(); err != nil {
				return err
			}
			_, err = sf.publish(uri, payload)
		} else {
			sf.Error("send data failed, %#v", err)
		}
	}
	return err
}

func (sf *AuthRequest) generateSign(deviceSecret string) error {
	var f func() hash.Hash

	if sf.SignMethod == signMethodHMACSHA1 {
		f = sha1.New
	} else {
		f = md5.New
		sf.SignMethod = signMethodHMACMD5
	}
	signSource := fmt.Sprintf("clientId%sdeviceName%sproductKey%stimestamp%d",
		sf.ClientID, sf.DeviceName, sf.ProductKey, sf.Timestamp)
	h := hmac.New(f, []byte(deviceSecret))
	if _, err := h.Write([]byte(signSource)); err != nil {
		return err
	}

	sf.Sign = hex.EncodeToString(h.Sum(nil))
	return nil
}