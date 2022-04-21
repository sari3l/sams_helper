package sams

import (
	"github.com/tidwall/gjson"
	"strconv"
)

func parseNormalGoodsV2(result gjson.Result) (error, NormalGoodsV2) {
	normalGoods := NormalGoodsV2{}
	normalGoods.SpuId = result.Get("spuId").Str
	normalGoods.StoreId = result.Get("storeId").Str
	normalGoods.Title = result.Get("title").Str
	normalGoods.SubTitle = result.Get("subTitle").Str
	for _, v := range result.Get("priceInfo").Array() {
		if priceStr := v.Get("priceTypeName").Str; priceStr == "销售价" || priceStr == "锁价" {
			price, _ := strconv.ParseInt(v.Get("price").Str, 10, 64)
			normalGoods.Price = price
		}
	}
	stockQuantity, _ := strconv.ParseInt(result.Get("stockInfo.stockQuantity").Str, 10, 64)
	normalGoods.StockQuantity = stockQuantity
	return nil, normalGoods
}

func (session *Session) GetGuaranteedSupplyGoods() (error, []NormalGoodsV2) {
	var goodsList = make([]NormalGoodsV2, 0)
	err, result := session.GetPageData("1187641882302384150")
	if err != nil {
		return err, nil
	}
	for _, v := range result.PageModuleVOList.Array() {
		if v.Get("moduleSign").Str == "goodsModule" {
			for _, v2 := range v.Get("renderContent.goodsList").Array() {
				_, goods := parseNormalGoodsV2(v2)
				goodsList = append(goodsList, goods)
			}
		}
	}
	return nil, goodsList
}
