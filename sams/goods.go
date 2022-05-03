package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
	"strconv"
)

type Goods struct {
	IsSelected bool   `json:"isSelected"`
	Quantity   int64  `json:"quantity"`
	SpuId      string `json:"spuId"`
	StoreId    string `json:"storeId"`
	StoreType  int64  `json:"storeType"`
}

type AddCartGoods struct {
	IsSelected       bool   `json:"isSelected"`
	IncreaseQuantity int64  `json:"increaseQuantity"`
	SpuId            string `json:"spuId"`
	StoreId          string `json:"storeId"`
	LabelList        string `json:"labelList"`
}

type DelCartGoods struct {
	SpuId     string `json:"spuId"`
	StoreId   string `json:"storeId"`
	Price     string `json:"price"`
	GoodsName string `json:"goodsName"`
}

func (goods ShowGoods) ToAddCartGoods(quantity int64) AddCartGoods {
	return AddCartGoods{
		IsSelected:       true,
		IncreaseQuantity: quantity,
		SpuId:            goods.SpuId,
		StoreId:          goods.StoreId,
		LabelList:        "",
	}
}

func (goods NormalGoods) ToAddCartGoods(quantity int64) AddCartGoods {
	return AddCartGoods{
		IsSelected:       true,
		IncreaseQuantity: quantity,
		SpuId:            goods.SpuId,
		StoreId:          goods.StoreId,
		LabelList:        "",
	}
}

func (goods NormalGoods) ToDelCartGoods() DelCartGoods {
	return DelCartGoods{
		SpuId:     goods.SpuId,
		StoreId:   goods.StoreId,
		Price:     fmt.Sprintf("%d.%d", goods.Price/100, goods.Price%100),
		GoodsName: goods.GoodsName,
	}
}

func (goods ShowGoods) ToNormalGoods() NormalGoods {
	return NormalGoods{
		SpuId:     goods.SpuId,
		StoreId:   goods.StoreId,
		Price:     goods.Price,
		GoodsName: goods.Title,
		BrandId:   goods.BrandId,
	}
}

type ShowGoods struct {
	SpuId         string `json:"spuId"`
	StoreId       string `json:"storeId"`
	Title         string `json:"title"`
	SubTitle      string `json:"subTitle"`
	Price         int64  `json:"price"`
	StockQuantity int64  `json:"stockQuantity"`
	BrandId       string `json:"brandId"`
}

func (this NormalGoods) ToGoods() Goods {
	return Goods{
		IsSelected: true,
		Quantity:   this.Quantity,
		SpuId:      this.SpuId,
		StoreId:    this.StoreId,
		StoreType:  this.StoreType,
	}
}

type PurchaseLimitV0 struct {
	LimitType          int64  `json:"limitType"`
	LimitNum           int64  `json:"limitNum"`
	StoreId            string `json:"storeId"`
	ResiduePurchaseNum int64  `json:"residuePurchaseNum"`
	Text               string `json:"text"`
	PopupText          string `json:"popupText"`
}

type NormalGoods struct {
	StoreId         string          `json:"storeId"`
	StoreType       int64           `json:"storeType"`
	SpuId           string          `json:"spuId"`
	SkuId           string          `json:"skuId"`
	BrandId         string          `json:"brandId"`
	GoodsName       string          `json:"goodsName"`
	Price           int64           `json:"price"`
	InvalidReason   string          `json:"invalidReason"`
	Quantity        int64           `json:"quantity"`
	PurchaseLimitV0 PurchaseLimitV0 `json:"purchaseLimitV0"`
	IsSelected      bool            `json:"isSelected"`
}

func parseShowGoods(result gjson.Result) (error, ShowGoods) {
	normalGoods := ShowGoods{}
	normalGoods.SpuId = result.Get("spuId").Str
	normalGoods.StoreId = result.Get("storeId").Str
	normalGoods.Title = result.Get("title").Str
	normalGoods.SubTitle = result.Get("subTitle").Str
	normalGoods.BrandId = result.Get("brandId").Str
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

func parsePurchaseLimitVO(result gjson.Result) (error, PurchaseLimitV0) {
	purchaseLimitV0 := PurchaseLimitV0{}
	if result.Type == 0 {
		purchaseLimitV0.LimitNum = 999
	} else {
		purchaseLimitV0.LimitType = result.Get("limitType").Int()
		purchaseLimitV0.LimitNum = result.Get("limitNum").Int()
		purchaseLimitV0.StoreId = result.Get("storeId").Str
		purchaseLimitV0.ResiduePurchaseNum = result.Get("residuePurchaseNum").Int()
		purchaseLimitV0.Text = result.Get("text").Str
		purchaseLimitV0.PopupText = result.Get("popupText").Str
	}
	return nil, purchaseLimitV0
}

func parseNormalGoods(result gjson.Result) (error, NormalGoods) {
	normalGoods := NormalGoods{}
	normalGoods.StoreId = result.Get("storeId").Str
	normalGoods.StoreType = result.Get("storeType").Int()
	normalGoods.SpuId = result.Get("spuId").Str
	normalGoods.SkuId = result.Get("skuId").Str
	normalGoods.BrandId = result.Get("brandId").Str
	normalGoods.GoodsName = result.Get("goodsName").Str
	normalGoods.Price = result.Get("price").Int()
	normalGoods.InvalidReason = result.Get("invalidReason").Str
	normalGoods.Quantity = result.Get("quantity").Int()
	normalGoods.IsSelected = result.Get("isSelected").Bool()
	_, purchaseLimitV0 := parsePurchaseLimitVO(result.Get("purchaseLimitVO"))
	normalGoods.PurchaseLimitV0 = purchaseLimitV0
	return nil, normalGoods
}

func (session *Session) CheckGoods() error {
	data := GoodsInfoParam{
		FloorId: session.FloorId,
		StoreId: "",
	}
	goodsList := make([]Goods, 0)
	for _, v := range session.Cart.FloorInfoList {
		if v.FloorId == session.FloorId {
			for _, v := range v.NormalGoodsList {
				if data.StoreId == "" {
					data.StoreId = v.StoreId
				}
				goodsList = append(goodsList, Goods{StoreId: v.StoreId, StoreType: v.StoreType, Quantity: v.Quantity, SpuId: v.SpuId})
			}
		}
	}
	data.GoodsList = goodsList
	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(GoodsInfoAPI, dataStr)
	if err != nil {
		return err
	}
	if result.Get("data.isHasException").Bool() == false {
		return nil
	} else {
		fmt.Printf("\n======== 以下商品已过期 ========\n")
		for index, v := range result.Get("data.popUpInfo.goodsList").Array() {
			_, goods := parseNormalGoods(v)
			fmt.Printf("[%v] 商品名：%s 商品ID：%s 商店ID：%v 总价：%d.%d\n", index, goods.GoodsName, goods.SpuId, goods.StoreId, goods.Price/100, goods.Price%100)
		}
		if session.Setting.IgnoreInvalid {
			return nil
		} else {
			return conf.OutOfSellErr
		}
	}
}
