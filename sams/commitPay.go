package sams

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"strconv"
)

type CommitPayParam struct {
	GoodsList        []Goods               `json:"goodsList"`
	InvoiceInfo      map[int64]interface{} `json:"invoiceInfo"`
	CartDeliveryType int64                 `json:"cartDeliveryType"`
	FloorId          int64                 `json:"floorId"`

	PurchaserName      string             `json:"purchaserName"`
	SettleDeliveryInfo SettleDeliveryInfo `json:"settleDeliveryInfo"`
	PayType            int64              `json:"payType"`
	Currency           string             `json:"currency"`
	Channel            string             `json:"channel"`
	ShortageId         int64              `json:"shortageId"`
	OrderType          int64              `json:"orderType"`
	Uid                string             `json:"uid"`
	AppId              string             `json:"appId"`
	AddressId          string             `json:"addressId"`
	DeliveryInfoVO     DeliveryInfoVO     `json:"deliveryInfoVO"`
	Remark             string             `json:"remark"`
	StoreInfo          StoreInfo          `json:"storeInfo"`
	ShortageDesc       string             `json:"shortageDesc"`
	PayMethodId        string             `json:"payMethodId"`
}

type IOSCommitPayParam struct {
	CommitPayParam
	Amount       string `json:"amount"`
	TradeType    string `json:"tradeType"`
	PurchaserId  string `json:"purchaserId"`
	IsSelfPickup int64  `json:"isSelfPickup"`
}

type MiniProgramCommitPayParam struct {
	CommitPayParam
	Amount                int64    `json:"amount"`
	LabelList             string   `json:"labelList"`
	IsSelectShoppingNotes bool     `json:"isSelectShoppingNotes"`
	CouponList            []string `json:"couponList"`
	SaasId                string   `json:"saasId"`
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
	DeliveryType         int64  `json:"deliveryType"`
	DeliveryDesc         string `json:"deliveryDesc"`
	DeliveryName         string `json:"deliveryName"`
	ExpectArrivalTime    string `json:"expectArrivalTime"`
	ExpectArrivalEndTime string `json:"expectArrivalEndTime"`
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
	_data := CommitPayParam{
		DeliveryInfoVO:   session.DeliveryInfoVO,
		StoreInfo:        session.FloorInfo.StoreInfo,
		Channel:          session.Channel,
		InvoiceInfo:      make(map[int64]interface{}),
		AddressId:        session.Address.AddressId,
		CartDeliveryType: session.Setting.DeliveryType,
		FloorId:          session.FloorId,
		GoodsList:        session.GoodsList,

		Currency:           "CNY",
		OrderType:          0,
		PayMethodId:        session.SubSaasId,
		PayType:            0,
		Remark:             "",
		ShortageDesc:       "其他商品继续配送（缺货商品直接退款）",
		ShortageId:         1,
		SettleDeliveryInfo: session.SettleDeliveryInfo,
		Uid:                session.Uid,
		AppId:              "wx111",
	}

	var dataStr []byte
	// 为了对照数据包，特意按设备类型排序观察
	switch session.Setting.DeviceType {
	case 2:
		_amount, _ := strconv.ParseInt(session.FloorInfo.Amount, 10, 64)
		data := MiniProgramCommitPayParam{
			CommitPayParam:        _data,
			Amount:                _amount,
			IsSelectShoppingNotes: true,
			CouponList:            []string{},
			LabelList:             "",
			SaasId:                "1818",
		}
		dataStr, _ = json.Marshal(data)
	default: // ios
		data := IOSCommitPayParam{
			CommitPayParam: _data,
			Amount:         session.FloorInfo.Amount,
			TradeType:      "APP",
			PurchaserId:    "",
			IsSelfPickup:   0,
		}
		dataStr, _ = json.Marshal(data)
	}

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
