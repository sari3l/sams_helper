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
	urlPath := StoreListAPI
	data := StoreListParam{
		Longitude: session.Address.Longitude,
		Latitude:  session.Address.Latitude,
	}

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

		default:
			return errors.New(result.Get("msg").Str)
		}

	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}

//func (session *Session) ChooseStore() error {
//	fmt.Printf("########## 根据商品类型分配商店 ###########\n")
//	err := session.GetStoreList()
//	if err != nil {
//		return err
//	}
//	for index, store := range session.StoreList {
//		fmt.Printf("[%v] Id：%s 名称：%s, 类型 ：%d\n", index, store.StoreId, store.StoreName, store.StoreType)
//	}
//	var index int
//	for true {
//		fmt.Println("\n选择说明:\n类型 2 普通商店\n类型 8 保税店\n类型 32 全球购\n\n请选择商店序号（0, 1, 2...)：")
//		stdin := bufio.NewReader(os.Stdin)
//		_, err := fmt.Fscanln(stdin, &index)
//		if err != nil {
//			fmt.Printf("输入有误：%s!\n", err)
//		} else if index >= len(session.AddressList) {
//			fmt.Println("输入有误：超过最大序号！")
//		} else {
//			break
//		}
//	}
//	session.Store = session.StoreList[index]
//	return nil
//}
