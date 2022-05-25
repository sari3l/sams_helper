package sams

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"sams_helper/tools"
)

type StoreInfoVO struct {
	StoreId           string  `json:"storeId"`
	StoreType         string  `json:"storeType"`
	StoreDeliveryAttr []int64 `json:"storeDeliveryAttr"`
}

type StoreInfo struct {
	StoreId                 string `json:"storeId"`
	StoreType               int64  `json:"storeType"`
	AreaBlockId             string `json:"areaBlockId"`
	StoreDeliveryTemplateId string `json:"-"`
	DeliveryModeId          string `json:"-"`
}

type Store struct {
	StoreId                 string  `json:"storeId"`
	StoreName               string  `json:"storeName"`
	StoreAddress            string  `json:"storeAddress"`
	StoreType               int64   `json:"storeType"`
	DeliveryModeId          string  `json:"deliveryModeId"`
	StoreDeliveryTemplateId string  `json:"storeDeliveryTemplateId"`
	AreaBlockId             string  `json:"areaBlockId"`
	AllDeliveryAttrList     []int64 `json:"allDeliveryAttrList"`
}

func (this *Store) ToStoreInfoVO() StoreInfoVO {
	return StoreInfoVO{
		StoreId:           this.StoreId,
		StoreType:         tools.Int64ToString(this.StoreType),
		StoreDeliveryAttr: this.AllDeliveryAttrList,
	}
}

func parseStore(storeData gjson.Result) (error, Store) {
	store := Store{}
	store.StoreId = storeData.Get("storeId").Str
	store.StoreName = storeData.Get("storeName").Str
	store.StoreAddress = storeData.Get("storeAddress").Str
	store.StoreType = storeData.Get("storeType").Int()
	store.DeliveryModeId = storeData.Get("storeDeliveryModeVerifyData.deliveryModeId").Str
	store.StoreDeliveryTemplateId = storeData.Get("storeRecmdDeliveryTemplateData.storeDeliveryTemplateId").Str
	store.AreaBlockId = storeData.Get("storeAreaBlockVerifyData.areaBlockId").Str
	_attrList := make([]int64, 0)
	for _, v := range storeData.Get("allDeliveryAttrList").Array() {
		_attrList = append(_attrList, v.Int())
	}
	store.AllDeliveryAttrList = _attrList
	return nil, store
}

func (session *Session) GetStoreList() error {
	data := StoreListParam{
		Longitude: session.Address.Longitude,
		Latitude:  session.Address.Latitude,
	}

	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(StoreListAPI, dataStr)
	if err != nil {
		return err
	}
	storeList := make([]Store, 0)
	for _, storeData := range result.Get("storeList").Array() {
		err, store := parseStore(storeData)
		if err != nil {
			return err
		}
		storeList = append(storeList, store)
	}
	session.StoreList = storeList
	return nil
}
