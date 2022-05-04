package sams

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"strings"
)

func parseGoodsList(g gjson.Result) (error, []ShowGoods) {
	goodsList := make([]ShowGoods, 0)
	for _, v := range g.Get("data.dataList").Array() {
		_, goods := parseShowGoods(v)
		goodsList = append(goodsList, goods)
	}
	return nil, goodsList
}

func goodsTitleMatch(goods []ShowGoods, keyword string) (error, []ShowGoods) {
	goodsList := make([]ShowGoods, 0)
	for _, v := range goods {
		if strings.Contains(v.Title, keyword) {
			goodsList = append(goodsList, v)
		}
	}
	return nil, goodsList
}

func (session *Session) GetGoodsFromSearch(keyword string) (error, []ShowGoods) {
	var goodsList = make([]ShowGoods, 0)
	var total int64 = 20
	var page int64 = 1      // 初始页数
	var pageSize int64 = 20 // 默认 20
	for (page-1)*pageSize <= total {
		data := GoodsPortalSearchParam{
			Filter:         make([]string, 0),
			Uid:            session.Uid,
			UidType:        3,
			PageSize:       pageSize,
			Sort:           0,
			Keyword:        keyword,
			UserUid:        session.Uid,
			AddressVO:      session.Address.ToAddressVO(),
			IsFastDelivery: false,
			PageNum:        page,
		}
		for _, v := range session.StoreList {
			data.StoreInfoVOList = append(data.StoreInfoVOList, v.ToStoreInfoVO())
		}

		dataStr, _ := json.Marshal(data)
		err, result := session.Request.POST(GoodsPortalSearchAPI, dataStr)
		if err != nil {
			return err, nil
		}
		total = result.Get("data.totalCount").Int()
		page += 1
		_, goodsListTmp := parseGoodsList(result)
		_, goodsListTmp = goodsTitleMatch(goodsListTmp, keyword)
		goodsList = append(goodsList, goodsListTmp...)
	}
	return nil, goodsList
}
