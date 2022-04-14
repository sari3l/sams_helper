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

type FloorInfo struct {
	FloorId         int64         `json:"floorId"`
	Amout           string        `json:"amout"`
	Quantity        int64         `json:"quantity"`
	StoreInfo       StoreInfo     `json:"storeInfo"`
	NormalGoodsList []NormalGoods `json:"normalGoodsList"`
}

type Cart struct {
	DeliveryAddress Address     `json:"deliveryAddress"`
	FloorInfoList   []FloorInfo `json:"floorInfoList"`
}

func (session *Session) parseFloorInfo(result gjson.Result) (error, FloorInfo) {
	floorInfo := FloorInfo{}
	floorInfo.FloorId = result.Get("floorId").Int()
	floorInfo.Amout = result.Get("amount").Str
	floorInfo.Quantity = result.Get("quantity").Int()
	floorInfo.StoreInfo = StoreInfo{
		StoreId:                 result.Get("storeInfo.storeId").Str,
		StoreType:               fmt.Sprintf("%d", result.Get("storeInfo.storeType").Int()),
		AreaBlockId:             result.Get("storeInfo.areaBlockId").Str,
		StoreDeliveryTemplateId: result.Get("storeInfo.storeDeliveryTemplateId").Str,
		DeliveryModeId:          result.Get("storeInfo.deliveryModeId").Str,
	}

	for _, v := range result.Get("normalGoodsList").Array() {
		_, normalGoods := session.parseNormalGoods(v)
		floorInfo.NormalGoodsList = append(floorInfo.NormalGoodsList, normalGoods)
	}

	return nil, floorInfo
}

func (session *Session) GetCartInfo(result gjson.Result) error {
	cart := Cart{}
	cart.FloorInfoList = make([]FloorInfo, 0)
	for _, v := range result.Get("data.floorInfoList").Array() {
		_, floor := session.parseFloorInfo(v)
		cart.FloorInfoList = append(cart.FloorInfoList, floor)
	}
	session.Cart = cart
	return nil
}

func (session *Session) CheckCart() error {
	urlPath := CartAPI
	data := CartParam{
		StoreList: session.StoreList,
		AddressId: session.Address.AddressId,
		Uid:       "",
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
			return session.GetCartInfo(result)
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
