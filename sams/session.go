package sams

import (
	"sams_helper/conf"
	"sams_helper/requests"
)

type Session struct {
	AuthToken          string             `json:"auth-token"`
	FloorId            int64              `json:"floor"`
	Channel            string             `json:"channel"` // wechat alipay china_unionpay sam_coupon
	SubSaasId          string             `json:"SubSaasId"`
	Address            Address            `json:"address"`
	Uid                string             `json:"uid"`
	SettleInfo         SettleInfo         `json:"settleInfo"`
	StoreList          []Store            `json:"store-list"`
	Cart               Cart               `json:"cart"`
	FloorInfo          FloorInfo          `json:"floorInfo"`
	GoodsList          []Goods            `json:"goodsList"`
	SettleDeliveryInfo SettleDeliveryInfo `json:"settleDeliveryInfo"`
	DeliveryInfoVO     DeliveryInfoVO     `json:"deliveryInfoVO"`
	Capacity           Capacity           `json:"capacity"`
	Request            requests.Request   `json:"request"`
	Setting            conf.Setting       `json:"setting"`
}

func (session *Session) InitSession(request requests.Request, setting conf.Setting) error {
	session.Request = request
	session.Setting = setting
	session.FloorId = setting.FloorId

	return session.CheckSession()
}

func (session *Session) CheckSession() error {
	if len(session.Setting.AuthToken) < 64 {
		return conf.AuthTokenErr
	}

	err, _ := session.Request.GET(AddressListAPI)
	if err != nil {
		return err
	}

	return nil
}
