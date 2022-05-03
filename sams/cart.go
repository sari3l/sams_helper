package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
	"strconv"
)

type FloorInfo struct {
	FloorId         int64         `json:"floorId"`
	Amount          string        `json:"amount"`
	Quantity        int64         `json:"quantity"`
	StoreInfo       StoreInfo     `json:"storeInfo"`
	NormalGoodsList []NormalGoods `json:"normalGoodsList"`
}

type Cart struct {
	DeliveryAddress Address     `json:"deliveryAddress"`
	FloorInfoList   []FloorInfo `json:"floorInfoList"`
}

func parseFloorInfo(result gjson.Result) (error, FloorInfo) {
	floorInfo := FloorInfo{}
	floorInfo.FloorId = result.Get("floorId").Int()
	floorInfo.Amount = result.Get("amount").Str
	floorInfo.Quantity = result.Get("quantity").Int()
	floorInfo.StoreInfo = StoreInfo{
		StoreId:                 result.Get("storeInfo.storeId").Str,
		StoreType:               result.Get("storeInfo.storeType").Int(),
		AreaBlockId:             result.Get("storeInfo.areaBlockId").Str,
		StoreDeliveryTemplateId: result.Get("storeInfo.storeDeliveryTemplateId").Str,
		DeliveryModeId:          result.Get("storeInfo.deliveryModeId").Str,
	}

	// 普通商品
	for _, v := range result.Get("normalGoodsList").Array() {
		_, normalGoods := parseNormalGoods(v)
		floorInfo.NormalGoodsList = append(floorInfo.NormalGoodsList, normalGoods)
	}

	// 促销商品
	for _, promotionGoodsList := range result.Get("promotionFloorGoodsList").Array() {
		for _, promotionGoods := range promotionGoodsList.Get("promotionGoodsList").Array() {
			_, p := parseNormalGoods(promotionGoods)
			floorInfo.NormalGoodsList = append(floorInfo.NormalGoodsList, p)
		}
	}

	// 有时间返回的 amount 为 “0”，为了方便显示按订单重新计算
	amount, _ := strconv.ParseInt(floorInfo.Amount, 10, 64)
	if amount == 0 {
		for _, v := range floorInfo.NormalGoodsList {
			amount += v.Quantity * v.Price
		}
		floorInfo.Amount = strconv.FormatInt(amount, 10)
	}

	return nil, floorInfo
}

func (session *Session) parseMiniProgramGoodsInfo(result gjson.Result) (error, FloorInfo) {
	floorInfo := FloorInfo{}
	floorInfo.FloorId = session.FloorId
	floorInfo.Amount = result.Get("selectedAmount").Str
	for _, v := range result.Get("normalGoodsList").Array() {
		_, normalGoods := parseNormalGoods(v)
		floorInfo.NormalGoodsList = append(floorInfo.NormalGoodsList, normalGoods)
		for _, s := range session.StoreList {
			if normalGoods.StoreId == s.StoreId && (floorInfo.StoreInfo.StoreType != session.Setting.StoreType) {
				floorInfo.StoreInfo = StoreInfo{
					StoreId:                 s.StoreId,
					StoreType:               s.StoreType,
					AreaBlockId:             s.AreaBlockId,
					StoreDeliveryTemplateId: s.StoreDeliveryTemplateId,
					DeliveryModeId:          s.DeliveryModeId,
				}
			}
		}
	}
	return nil, floorInfo
}

func (session *Session) SetCartInfo(result gjson.Result) error {
	cart := Cart{}
	cart.FloorInfoList = make([]FloorInfo, 0)
	switch session.Setting.DeviceType {
	case 1:
		for _, v := range result.Get("data.floorInfoList").Array() {
			if "失效商品" == v.Get("floorName").Str {
				continue
			}
			_, floor := parseFloorInfo(v)
			cart.FloorInfoList = append(cart.FloorInfoList, floor)
		}
	case 2:
		for _, v := range result.Get("data.miniProgramGoodsInfo").Array() {
			_, floor := session.parseMiniProgramGoodsInfo(v)
			floor.Amount = result.Get("data.selectedAmount").Str
			floor.Quantity = result.Get("data.selectedNumber").Int()
			cart.FloorInfoList = append(cart.FloorInfoList, floor)
		}
	default:
		return conf.DeliveryTypeErr
	}
	session.Cart = cart
	return nil
}

func (session *Session) ModifyCartGoodsInfo(goods Goods) error {
	data := ModifyCartGoodsInfoParam{
		CartGoodsInfo: goods,
		Uid:           session.Uid,
	}
	dataStr, _ := json.Marshal(data)
	err, _ := session.Request.POST(ModifyCartGoodsInfoAPI, dataStr)
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) AddCartGoodsInfo(addGoodsList []AddCartGoods) error {
	data := AddCartGoodsInfoParam{
		CartGoodsInfoList: addGoodsList,
		Uid:               session.Uid,
	}
	dataStr, _ := json.Marshal(data)
	err, _ := session.Request.POST(AddCartGoodsInfoAPI, dataStr)
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) DelCartGoodsInfo(delGoodsList []DelCartGoods) error {
	data := DelCartGoodsInfoParam{
		CartGoodsList: delGoodsList,
		Uid:           session.Uid,
	}
	dataStr, _ := json.Marshal(data)
	err, _ := session.Request.POST(DelCartGoodsInfoAPI, dataStr)
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) FixCart() (error, bool, bool) {
	isChangedOffline := false
	isChangedOnline := false
	var removeQuantity int64 = 0
	var removeAmount int64 = 0
	for index, v := range session.Cart.FloorInfoList {
		for index2, v2 := range v.NormalGoodsList {
			if v2.PurchaseLimitV0.LimitNum < v2.Quantity {
				// offline
				fmt.Printf("[!] 校验发现限购商品：%s，限购数量：%d，现有数量：%d，正在修正中...\n", v2.GoodsName, v2.PurchaseLimitV0.LimitNum, v2.Quantity)
				if session.Setting.AutoFixPurchaseLimitSet.FixOffline && !session.Setting.AutoFixPurchaseLimitSet.FixOnline {
					removeQuantity += v2.Quantity - v2.PurchaseLimitV0.LimitNum
					removeAmount += (v2.Quantity - v2.PurchaseLimitV0.LimitNum) * v2.Price
					v2.Quantity = v2.PurchaseLimitV0.LimitNum
					v.NormalGoodsList[index2] = v2
					isChangedOffline = true
				}
				// online
				if session.Setting.AutoFixPurchaseLimitSet.FixOnline {
					_goods := v2.ToGoods()
					_goods.Quantity = v2.PurchaseLimitV0.LimitNum
					if err := session.ModifyCartGoodsInfo(_goods); err != nil {
						return conf.FixCartErr, isChangedOffline, true
					}
					isChangedOnline = true
				}

			}
		}
		session.Cart.FloorInfoList[index].Quantity -= removeQuantity
		_amount, _ := strconv.ParseInt(session.Cart.FloorInfoList[index].Amount, 10, 64)
		session.Cart.FloorInfoList[index].Amount = strconv.FormatInt(_amount-removeAmount, 10)
	}
	return nil, isChangedOffline, isChangedOnline
}

func (session *Session) CheckCart() error {
	session.Cart = Cart{}
	data := CartParam{
		StoreList:         session.StoreList,
		AddressId:         session.Address.AddressId,
		Uid:               session.Uid,
		DeliveryType:      fmt.Sprintf("%d", session.Setting.DeliveryType),
		HomePagelatitude:  session.Address.Latitude,
		HomePagelongitude: session.Address.Longitude,
	}
	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(CartAPI, dataStr)
	if err != nil {
		return err
	}
	return session.SetCartInfo(result)
}
