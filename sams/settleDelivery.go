package sams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
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

func parseSettleDelivery(g gjson.Result) (error, SettleDelivery) {
	r := SettleDelivery{
		DeliveryType:            g.Get("deliveryType").Int(),
		DeliveryName:            g.Get("deliveryName").Str,
		DeliveryDesc:            g.Get("deliveryDesc").Str,
		ExpectArrivalTime:       g.Get("expectArrivalTime").Str,
		ExpectArrivalEndTime:    g.Get("expectArrivalEndTime").Str,
		StoreDeliveryTemplateId: g.Get("storeDeliveryTemplateId").Str,
		AreaBlockId:             g.Get("AreaBlockId").Str,
		AreaBlockName:           g.Get("areaBlockName").Str,
		FirstPeriod:             g.Get("firstPeriod").Int(),
	}

	for _, v := range g.Get("deliveryModeIdList").Array() {
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
	StoreType               string `json:"storeType"`
	AreaBlockId             string `json:"areaBlockId"`
	StoreDeliveryTemplateId string `json:"-"`
	DeliveryModeId          string `json:"-"`
}

type DeliveryInfoVO struct {
	StoreDeliveryTemplateId string `json:"storeDeliveryTemplateId"`
	DeliveryModeId          string `json:"deliveryModeId"`
	StoreType               string `json:"storeType"`
}

type SettleParam struct {
	Uid              string         `json:"uid"`
	AddressId        string         `json:"addressId"`
	DeliveryInfoVO   DeliveryInfoVO `json:"deliveryInfoVO"`
	CartDeliveryType int64          `json:"cartDeliveryType"`
	StoreInfo        StoreInfo      `json:"storeInfo"`
	CouponList       []string       `json:"couponList"`
	IsSelfPickup     int64          `json:"isSelfPickup"`
	FloorId          int64          `json:"floorId"`
	GoodsList        []Goods        `json:"goodsList"`
}

func (session *Session) CheckSettleInfo() error {
	urlPath := SettleInfoAPI

	data := SettleParam{
		Uid:              session.Uid,
		AddressId:        session.Address.AddressId,
		DeliveryInfoVO:   session.DeliveryInfoVO,
		CartDeliveryType: 2,
		StoreInfo:        session.FloorInfo.StoreInfo,
		CouponList:       make([]string, 0),
		IsSelfPickup:     0,
		FloorId:          session.FloorId,
		GoodsList:        session.GoodsList,
	}

	dataStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", urlPath, bytes.NewReader(dataStr))
	req.Header = *session.Headers

	resp, err := session.Client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {
		result := gjson.Parse(string(body))
		switch result.Get("code").Str {
		case "Success":
			return session.GetSettleInfo(result)
		case "LIMITED":
			return LimitedErr
		case "CART_GOOD_CHANGE":
			return CartGoodChangeErr
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
