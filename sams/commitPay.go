package sams

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"sams_helper/tools"
)

type OrderInfo struct {
	IsSuccess bool    `json:"isSuccess"`
	OrderNo   string  `json:"orderNo"`
	PayAmount string  `json:"payAmount"`
	Channel   string  `json:"channel"`
	PayInfo   PayInfo `json:"PayInfo"`
}

type PayInfo struct {
	PayInfo    string `json:"payInfo"`
	OutTradeNo string `json:"outTradeNo"`
	TotalAmt   int64  `json:"totalAmt"`
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
		IsSuccess: result.Get("isSuccess").Bool(),
		OrderNo:   result.Get("orderNo").Str,
		PayAmount: result.Get("payAmount").Str,
		Channel:   result.Get("channel").Str,
		PayInfo: PayInfo{
			PayInfo:    result.Get("payInfo.PayInfo").Str,
			OutTradeNo: result.Get("payInfo.OutTradeNo").Str,
			TotalAmt:   result.Get("payInfo.TotalAmt").Int(),
		},
	}
	return nil, order
}

func (session *Session) CommitPay() (error, OrderInfo) {
	data := CommitPayParam{
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
		CouponList:         make([]CouponInfo, 0),
	}
	if len(session.CouponList) > 0 {
		for _, v := range session.CouponList {
			data.CouponList = append(data.CouponList, CouponInfo{PromotionId: v.RuleId, StoreId: session.FloorInfo.StoreInfo.StoreId})
		}
	}

	var dataStr []byte
	// 为了对照数据包，特意按设备类型排序观察
	switch session.Setting.DeviceType {
	case 2:
		amount := tools.StringToInt64(session.SettleInfo.TotalAmount)
		data := MiniProgramCommitPayParam{
			CommitPayParam:        data,
			Amount:                amount,
			IsSelectShoppingNotes: true,
			LabelList:             "",
			SaasId:                session.Setting.SassId,
			AppId:                 "",
		}
		dataStr, _ = json.Marshal(data)
	default: // ios
		data := IOSCommitPayParam{
			CommitPayParam: data,
			Amount:         session.SettleInfo.TotalAmount,
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
