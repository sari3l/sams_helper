package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
	"sams_helper/tools"
)

type Goods struct {
	IsSelected bool   `json:"isSelected"`
	Quantity   int64  `json:"quantity"`
	SpuId      string `json:"spuId"`
	StoreId    string `json:"storeId"`
	StoreType  int64  `json:"storeType"`
	GoodsName  string
	Price      int64
	Weight     float64
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

func (goods Goods) ToAddCartGoods(quantity int64) AddCartGoods {
	return AddCartGoods{
		IsSelected:       true,
		IncreaseQuantity: quantity,
		SpuId:            goods.SpuId,
		StoreId:          goods.StoreId,
		LabelList:        "",
	}
}

func (goods Goods) ToDelCartGoods() DelCartGoods {
	return DelCartGoods{
		SpuId:     goods.SpuId,
		StoreId:   goods.StoreId,
		Price:     fmt.Sprintf("%d.%d", goods.Price/100, goods.Price%100),
		GoodsName: goods.GoodsName,
	}
}

func (goods NormalGoods) ToAddCartGoods(quantity int64) AddCartGoods {
	return goods.ToGoods().ToAddCartGoods(quantity)
}

func (goods NormalGoods) ToDelCartGoods() DelCartGoods {
	return goods.ToGoods().ToDelCartGoods()
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
	SpuId         string  `json:"spuId"`
	StoreId       string  `json:"storeId"`
	Title         string  `json:"title"`
	SubTitle      string  `json:"subTitle"`
	Price         int64   `json:"price"`
	StockQuantity int64   `json:"stockQuantity"`
	BrandId       string  `json:"brandId"`
	Weight        float64 `json:"weight"`
}

func (goods NormalGoods) ToGoods() Goods {
	return Goods{
		IsSelected: true,
		Quantity:   goods.Quantity,
		SpuId:      goods.SpuId,
		StoreId:    goods.StoreId,
		StoreType:  goods.StoreType,
		GoodsName:  goods.GoodsName,
		Price:      goods.Price,
		Weight:     goods.Weight,
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
	StockQuantity   int64           `json:"stockQuantity"`
	Weight          float64         `json:"weight"`
}

func parseShowGoods(result gjson.Result) (error, ShowGoods) {
	showGoods := ShowGoods{}
	showGoods.SpuId = result.Get("spuId").Str
	showGoods.StoreId = result.Get("storeId").Str
	showGoods.Title = result.Get("title").Str
	showGoods.SubTitle = result.Get("subTitle").Str
	showGoods.BrandId = result.Get("brandId").Str
	showGoods.Weight = result.Get("weight").Num
	for _, v := range result.Get("priceInfo").Array() {
		if priceStr := v.Get("priceTypeName").Str; priceStr == "销售价" || priceStr == "锁价" {
			price := tools.StringToInt64(v.Get("price").Str)
			showGoods.Price = price
		}
	}
	stockQuantity := tools.StringToInt64(result.Get("stockInfo.stockQuantity").Str)
	showGoods.StockQuantity = stockQuantity
	return nil, showGoods
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
	normalGoods.StockQuantity = result.Get("stockQuantity").Int()
	normalGoods.IsSelected = result.Get("isSelected").Bool()
	normalGoods.Weight = result.Get("weight").Num
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
	if result.Get("isHasException").Bool() == false {
		return nil
	} else {
		fmt.Printf("\n======== 以下商品已过期 ========\n")
		for index, v := range result.Get("popUpInfo.goodsList").Array() {
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

func (session *Session) QueryGoodsDetail(spuId string) (error, ShowGoods) {
	data := QueryDetailParam{
		SpuId: spuId,
	}
	for _, v := range session.StoreList {
		data.StoreInfoVOList = append(data.StoreInfoVOList, v.ToStoreInfoVO())
	}

	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(QueryDetailAPI, dataStr)
	if err != nil {
		return err, ShowGoods{}
	}
	return parseShowGoods(result)
}
