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

type CommitPayPram struct {
	GoodsList          []Goods                `json:"goodsList"`
	InvoiceInfo        map[int64]interface{}  `json:"invoiceInfo"`
	CartDeliveryType   int64                  `json:"cartDeliveryType"`
	FloorId            int64                  `json:"floorId"`
	Amount             string                 `json:"amount"`
	PurchaserName      string                 `json:"purchaserName"`
	SettleDeliveryInfo map[string]interface{} `json:"settleDeliveryInfo"`
	TradeType          string                 `json:"tradeType"` //"APP"
	PurchaserId        string                 `json:"purchaserId"`
	PayType            int64                  `json:"payType"`
	Currency           string                 `json:"currency"`     // CNY
	Channel            string                 `json:"channel"`      // wechat
	ShortageId         int64                  `json:"shortageId"`   //1
	IsSelfPickup       int64                  `json:"isSelfPickup"` //0
	OrderType          int64                  `json:"orderType"`    //0
	Uid                string                 `json:"uid"`
	AppId              string                 `json:"appId"`
	AddressId          string                 `json:"addressId"`
	DeliveryInfoVO     DeliveryInfoVO         `json:"deliveryInfoVO"`
	Remark             string                 `json:"remark"`
	StoreInfo          StoreInfo              `json:"storeInfo"`
	ShortageDesc       string                 `json:"shortageDesc"`
	PayMethodId        string                 `json:"payMethodId"`
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
	TotalAmt   int    `json:"TotalAmt"`
}

func (session *Session) GetOrderInfo(result gjson.Result) error {
	session.OrderInfo = OrderInfo{
		IsSuccess: result.Get("data.isSuccess").Bool(),
		OrderNo:   result.Get("data.orderNo").Str,
		PayAmount: result.Get("data.payAmount").Str,
		Channel:   result.Get("data.channel").Str,
		PayInfo: PayInfo{
			PayInfo:    result.Get("data.PayInfo.PayInfo").Str,
			OutTradeNo: result.Get("data.PayInfo.OutTradeNo").Str,
			TotalAmt:   int(result.Get("data.PayInfo.TotalAmt").Num),
		},
	}
	return nil
}

func (session *Session) CommitPay() error {
	urlPath := CommitPayAPI

	data := CommitPayPram{
		GoodsList:          session.GoodsList,
		InvoiceInfo:        make(map[int64]interface{}),
		CartDeliveryType:   2, // 1,急速到达 2,全城配送
		FloorId:            0,
		Amount:             "13123", //测试没用但必须有
		PurchaserName:      "",
		SettleDeliveryInfo: map[string]interface{}{"deliveryType": 0},
		//SettleDeliveryInfo: map[string]interface{}{"deliveryType": 0, "expectArrivalTime": nil, "expectArrivalEndTime": nil},
		TradeType:      "APP",
		PurchaserId:    "",
		PayType:        0,
		Currency:       "CNY",
		Channel:        session.Channel,
		ShortageId:     1,
		IsSelfPickup:   0,
		OrderType:      0,
		Uid:            "123", //s.Uid,
		AppId:          fmt.Sprintf("123"),
		AddressId:      session.Address.AddressId,
		DeliveryInfoVO: session.DeliveryInfoVO,
		Remark:         "",
		StoreInfo:      session.FloorInfo.StoreInfo,
		ShortageDesc:   "其他商品继续配送（缺货商品直接退款）",
		PayMethodId:    session.SubSaasId,
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
			if result.Get("data.isSuccess").Bool() {
				return session.GetOrderInfo(result)
			}
			return errors.New(result.Get("data.failReason").Str)
		case "LIMITED":
			return LimitedErr1
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
