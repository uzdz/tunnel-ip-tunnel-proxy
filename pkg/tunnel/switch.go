package tunnel

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"ip-tunnel-proxy/pkg/config"
	"ip-tunnel-proxy/pkg/config/po"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

func GetProxyDst(uid string, s bool) (dst string, e error) {
	if s == false {
		get, b := config.TunnelUserExp.Get(uid)

		if b {

			ipWithExp := get.(string)
			split := strings.Split(strings.Trim(ipWithExp, " "), "#")

			exp, err := strconv.ParseInt(split[1], 10, 64)
			if err == nil {
				if time.Now().Unix() <= exp {
					return split[0], e
				}
			}
		}
	}

	urq := po.URq{
		DeviceId: config.Number,
		Uid:      uid,
		SwitchIp: s,
		AppId:    config.AppId,
		PoolId:   config.PoolId,
	}

	if config.UserConnectLimit != nil {
		urq.UserConnectLimit = config.UserConnectLimit[uid]
	}

	e, v := reHttpProtoPost(urq, config.SwitchRemoteUrl)

	if v.GetCode() == 200 || v.Ip != "" {
		config.TunnelUserExp.Set(uid, v.Ip+"#"+string(v.ExpiredTime))
		return v.Ip, e
	}

	return "", fmt.Errorf("未申请到可代理IP或被限流！")
}

func reHttpProtoPost(rq po.URq, url string) (err error, v po.URp) {

	bytesData, err := proto.Marshal(&rq)
	if err != nil {
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/x-protobuf")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	urp := &po.URp{}
	proto.Unmarshal(respBytes, urp)

	return nil, *urp
}
