package sams

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
)

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
		DeliveryInfoVO:     session.DeliveryInfoVO,
		StoreInfo:          session.FloorInfo.StoreInfo,
		Channel:            session.Channel,
		InvoiceInfo:        make(map[int64]interface{}),
		AddressId:          session.Address.AddressId,
		CartDeliveryType:   session.Setting.DeliveryType,
		FloorId:            session.FloorId,
		GoodsList:          session.GoodsList,
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
		CouponList:         make([]string, 0),
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
			LabelList:             "",
			SaasId:                session.Setting.SassId,
			AppId:                 "",
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
	return session.GetOrderInfo(result)
}
