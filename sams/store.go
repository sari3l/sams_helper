package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type Store struct {
	StoreId                 string `json:"storeId"`
	StoreName               string `json:"storeName"`
	StoreAddress            string `json:"storeAddress"`
	StoreType               int64  `json:"storeType"`
	DeliveryModeId          string `json:"deliveryModeId"`
	StoreDeliveryTemplateId string `json:"storeDeliveryTemplateId"`
	AreaBlockId             string `json:"areaBlockId"`
}

func (session *Session) parseStore(storeData gjson.Result) (error, Store) {
	store := Store{}
	store.StoreId = storeData.Get("storeId").Str
	store.StoreName = storeData.Get("storeName").Str
	store.StoreAddress = storeData.Get("storeAddress").Str
	store.StoreType = storeData.Get("storeType").Int()
	store.DeliveryModeId = storeData.Get("storeDeliveryModeVerifyData.deliveryModeId").Str
	store.StoreDeliveryTemplateId = storeData.Get("storeRecmdDeliveryTemplateData.storeDeliveryTemplateId").Str
	store.AreaBlockId = storeData.Get("storeAreaBlockVerifyData.areaBlockId").Str
	return nil, store
}

func (session *Session) GetStoreList() error {
	fmt.Printf("\n########## 获取就近商店信息 ###########\n")
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
	for _, storeData := range result.Get("data.storeList").Array() {
		err, store := session.parseStore(storeData)
		if err != nil {
			return err
		}
		storeList = append(storeList, store)
	}
	session.StoreList = storeList
	return nil
}
