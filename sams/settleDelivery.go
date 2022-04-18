package sams

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type SettleDelivery struct {
	DeliveryType            int64    `json:"deliveryType"` // 1,极速达 2, 全城配 3, 物流配送
	DeliveryName            string   `json:"deliveryName"`
	DeliveryDesc            string   `json:"deliveryDesc"`
	ExpectArrivalTime       string   `json:"expectArrivalTime"`
	ExpectArrivalEndTime    string   `json:"expectArrivalEndTime"`
	StoreDeliveryTemplateId string   `json:"storeDeliveryTemplateId"`
	DeliveryModeIdList      []string `json:"deliveryModeIdList"`
	AreaBlockId             string   `json:"areaBlockId"`
	AreaBlockName           string   `json:"areaBlockName"`
	FirstPeriod             int64    `json:"firstPeriod"`
}

func parseSettleDelivery(result gjson.Result) (error, SettleDelivery) {
	r := SettleDelivery{
		DeliveryType:            result.Get("deliveryType").Int(),
		DeliveryName:            result.Get("deliveryName").Str,
		DeliveryDesc:            result.Get("deliveryDesc").Str,
		ExpectArrivalTime:       result.Get("expectArrivalTime").Str,
		ExpectArrivalEndTime:    result.Get("expectArrivalEndTime").Str,
		StoreDeliveryTemplateId: result.Get("storeDeliveryTemplateId").Str,
		AreaBlockId:             result.Get("AreaBlockId").Str,
		AreaBlockName:           result.Get("areaBlockName").Str,
		FirstPeriod:             result.Get("firstPeriod").Int(),
	}

	for _, v := range result.Get("deliveryModeIdList").Array() {
		r.DeliveryModeIdList = append(r.DeliveryModeIdList, v.Str)
	}
	return nil, r
}

type SettleInfo struct {
	SaasId          string         `json:"saasId"`
	Uid             string         `json:"uid"`
	FloorId         int64          `json:"floorId"`
	FloorName       string         `json:"floorName"`
	SettleDelivery  SettleDelivery `json:"settleDelivery"`
	DeliveryAddress Address        `json:"deliveryAddress"`
}

func (session *Session) GetSettleInfo(result gjson.Result) error {
	r := SettleInfo{}

	for _, v := range result.Get("data.settleDelivery").Array() {
		_, settleDelivery := parseSettleDelivery(v)
		r.SettleDelivery = settleDelivery

	}
	r.SaasId = result.Get("data.saasId").Str
	r.Uid = result.Get("data.uid").Str
	r.FloorId = result.Get("data.floorId").Int()
	r.FloorName = result.Get("data.floorName").Str
	err, address := parseAddress(result.Get("data.deliveryAddress"))
	if err == nil {
		r.DeliveryAddress = address
	}

	session.SettleInfo = r
	return nil
}

type StoreInfo struct {
	StoreId                 string `json:"storeId"`
	StoreType               int64  `json:"storeType"`
	AreaBlockId             string `json:"areaBlockId"`
	StoreDeliveryTemplateId string `json:"-"`
	DeliveryModeId          string `json:"-"`
}

type DeliveryInfoVO struct {
	StoreDeliveryTemplateId string `json:"storeDeliveryTemplateId"`
	DeliveryModeId          string `json:"deliveryModeId"`
	StoreType               int64  `json:"storeType"`
}

func (session *Session) CheckSettleInfo() error {
	data := SettleParam{
		Uid:              session.Uid,
		AddressId:        session.Address.AddressId,
		DeliveryInfoVO:   session.DeliveryInfoVO,
		CartDeliveryType: session.Setting.DeliveryType,
		StoreInfo:        session.FloorInfo.StoreInfo,
		CouponList:       make([]string, 0),
		IsSelfPickup:     0,
		FloorId:          session.FloorId,
		GoodsList:        session.GoodsList,
	}

	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(SettleInfoAPI, dataStr)
	if err != nil {
		return err
	}
	return session.GetSettleInfo(result)
}
