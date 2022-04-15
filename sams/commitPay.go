package sams

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type CommitPayPram struct {
	GoodsList          []Goods               `json:"goodsList"`
	InvoiceInfo        map[int64]interface{} `json:"invoiceInfo"`
	CartDeliveryType   int64                 `json:"cartDeliveryType"`
	FloorId            int64                 `json:"floorId"`
	Amount             string                `json:"amount"`
	PurchaserName      string                `json:"purchaserName"`
	SettleDeliveryInfo SettleDeliveryInfo    `json:"settleDeliveryInfo"`
	TradeType          string                `json:"tradeType"` //"APP"
	PurchaserId        string                `json:"purchaserId"`
	PayType            int64                 `json:"payType"`
	Currency           string                `json:"currency"`     // CNY
	Channel            string                `json:"channel"`      // wechat
	ShortageId         int64                 `json:"shortageId"`   //1
	IsSelfPickup       int64                 `json:"isSelfPickup"` //0
	OrderType          int64                 `json:"orderType"`    //0
	Uid                string                `json:"uid"`
	AppId              string                `json:"appId"`
	AddressId          string                `json:"addressId"`
	DeliveryInfoVO     DeliveryInfoVO        `json:"deliveryInfoVO"`
	Remark             string                `json:"remark"`
	StoreInfo          StoreInfo             `json:"storeInfo"`
	ShortageDesc       string                `json:"shortageDesc"`
	PayMethodId        string                `json:"payMethodId"`
	LabelList          string                `json:"labelList"`
}

type OrderInfo struct {
	IsSuccess bool    `json:"isSuccess"`
	OrderNo   string  `json:"orderNo"`
	PayAmount string  `json:"payAmount"`
	Channel   string  `json:"channel"`
	PayInfo   PayInfo `json:"PayInfo"`
}

type PayInfo struct {
	PayInfo    string `json:"PayInfo"`
	OutTradeNo string `json:"OutTradeNo"`
	TotalAmt   int64  `json:"TotalAmt"`
}

type SettleDeliveryInfo struct {
	DeliveryType         int64  `json:"deliveryType"`         //默认0
	ExpectArrivalTime    string `json:"expectArrivalTime"`    //配送时间: 1649922300000
	ExpectArrivalEndTime string `json:"expectArrivalEndTime"` //配送时间
	ArrivalTimeStr       string `json:"-"`
}

func (session *Session) GetOrderInfo(result gjson.Result) (error, OrderInfo) {
	order := OrderInfo{
		IsSuccess: result.Get("data.isSuccess").Bool(),
		OrderNo:   result.Get("data.orderNo").Str,
		PayAmount: result.Get("data.payAmount").Str,
		Channel:   result.Get("data.channel").Str,
		PayInfo: PayInfo{
			PayInfo:    result.Get("data.PayInfo.PayInfo").Str,
			OutTradeNo: result.Get("data.PayInfo.OutTradeNo").Str,
			TotalAmt:   result.Get("data.PayInfo.TotalAmt").Int(),
		},
	}
	return nil, order
}

func (session *Session) CommitPay() (error, OrderInfo) {
	data := CommitPayPram{
		GoodsList:          session.GoodsList,
		InvoiceInfo:        make(map[int64]interface{}),
		CartDeliveryType:   session.Setting.DeliveryType, // 1,急速到达 2,全城配送
		FloorId:            session.FloorId,
		Amount:             "13123", //测试没用但必须有
		PurchaserName:      "",
		SettleDeliveryInfo: session.SettleDeliveryInfo,
		TradeType:          "APP",
		PurchaserId:        "",
		PayType:            0,
		Currency:           "CNY",
		Channel:            session.Channel,
		ShortageId:         1,
		IsSelfPickup:       0,
		OrderType:          0,
		LabelList:          "",    // 小程序模式必须有
		Uid:                "123", //s.Uid,
		AppId:              fmt.Sprintf("123"),
		AddressId:          session.Address.AddressId,
		DeliveryInfoVO:     session.DeliveryInfoVO,
		Remark:             "",
		StoreInfo:          session.FloorInfo.StoreInfo,
		ShortageDesc:       "其他商品继续配送（缺货商品直接退款）",
		PayMethodId:        session.SubSaasId,
	}

	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(CommitPayAPI, dataStr)
	if err != nil {
		return err, OrderInfo{}
	}
	if session.Setting.DeviceType == 2 {
		return session.GetOrderInfo(result)
	} else if session.Setting.DeviceType == 1 && result.Get("data.isSuccess").Bool() {
		return session.GetOrderInfo(result)
	} else {
		return errors.New(result.Get("data.failReason").Str), OrderInfo{}
	}
}
