package sams

import (
	"net/http"
	"net/url"
	"time"
)

type Session struct {
	AuthToken      string         `json:"auth-token"`
	FloorId        int64          `json:"floor"`
	Channel        string         `json:"channel"` // wechat alipay china_unionpay sam_coupon
	SubSaasId      string         `json:"SubSaasId"`
	Address        Address        `json:"address"`
	Uid            string         `json:"uid"`
	SettleInfo     SettleInfo     `json:"settleInfo"`
	AddressList    []Address      `json:"address-list"`
	Store          Store          `json:"store"`
	StoreList      []Store        `json:"store-list"`
	Cart           Cart           `json:"cart"`
	FloorInfo      FloorInfo      `json:"floorInfo"`
	GoodsList      []Goods        `json:"goodsList"`
	DeliveryInfoVO DeliveryInfoVO `json:"deliveryInfoVO"`
	Capacity       Capacity       `json:"capacity"`
	OrderInfo      OrderInfo      `json:"orderInfo"`
	Client         *http.Client   `json:"client"`
	Headers        *http.Header   `json:"headers"`
}

func (session *Session) InitSession(AuthToken string) error {
	session.AuthToken = AuthToken


	session.Client = &http.Client{
		Timeout:   60 * time.Second,
	}

	session.FloorId = 1
	session.Headers = &http.Header{
		"Host":            []string{"api-sams.walmartmobile.cn"},
		"content-Type":    []string{"application/json"},
		"device-type":     []string{"ios"},
		"accept":          []string{"*/*"},
		"auth-token":      []string{session.AuthToken},
		"user-agent":      []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E217 MicroMessenger/6.8.0(0x16080000) NetType/WIFI Language/en Branch/Br_trunk MiniProgramEnv/Mac"},
		"Accept-Language": []string{"zh-Hans-CN;q=1, en-CN;q=0.9, ga-IE;q=0.8"},
	}

	// 设置地址
	err := session.ChooseAddress()
	if err != nil {
		return err
	}

	// 设置支付方式
	err = session.ChoosePayment()
	if err != nil {
		return err
	}

	return nil
}
