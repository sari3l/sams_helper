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

type Goods struct {
	IsSelected bool   `json:"isSelected"`
	Quantity   int64  `json:"quantity"`
	SpuId      string `json:"spuId"`
	StoreId    string `json:"storeId"`
}

func (this NormalGoods) ToGoods() Goods {
	return Goods{
		IsSelected: true,
		Quantity:   this.Quantity,
		SpuId:      this.SpuId,
		StoreId:    this.StoreId,
	}
}

type NormalGoods struct {
	StoreId       string `json:"storeId"`
	StoreType     int64  `json:"storeType"`
	SpuId         string `json:"spuId"`
	SkuId         string `json:"skuId"`
	BrandId       string `json:"brandId"`
	GoodsName     string `json:"goodsName"`
	Price         int64  `json:"price"`
	InvalidReason string `json:"invalidReason"`
	Quantity      int64  `json:"quantity"`
}

func (session *Session) parseNormalGoods(result gjson.Result) (error, NormalGoods) {
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
	return nil, normalGoods
}

func (session *Session) CheckGoods() error {
	urlPath := GoodsInfoAPI
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
				goodsList = append(goodsList, Goods{StoreId: v.StoreId, Quantity: v.Quantity, SpuId: v.SpuId})
			}
		}
	}
	data.GoodsList = goodsList
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
			if result.Get("data.isHasException").Bool() == false {
				return nil
			} else {
				fmt.Printf("\n======== 以下商品已过期 ========\n")
				for index, v := range result.Get("data.popUpInfo.goodsList").Array() {
					_, goods := session.parseNormalGoods(v)
					fmt.Printf("[%v] %s 数量：%v 总价：%d\n", index, goods.SpuId, goods.StoreId, goods.Price)
				}
				return OOSErr
			}
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
