package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
)

type Address struct {
	AddressId       string `json:"addressId"`
	Mobile          string `json:"mobile"`          // 手机号
	Name            string `json:"name"`            // 用户姓名
	CountryName     string `json:"countryName"`     // 国家
	ProvinceName    string `json:"provinceName"`    // 省份
	CityName        string `json:"cityName"`        // 城市：上海市
	DistrictName    string `json:"districtName"`    // 区域：长宁区
	ReceiverAddress string `json:"receiverAddress"` // 小区：绿园一村
	DetailAddress   string `json:"detailAddress"`   // 楼栋：XX幢XXX室
	IsDefault       int64  `json:"isDefault"`       // 优先级
	AddressTag      string `json:"addressTag"`      // 住址标签
	Latitude        string `json:"latitude"`        // 经度
	Longitude       string `json:"longitude"`       // 维度
	CreateTime      string `json:"createTime"`      // 创建时间
	UpdateTime      string `json:"updateTime"`      // 更新时间
}

func parseAddress(addressData gjson.Result) (error, Address) {
	address := Address{}
	address.AddressId = addressData.Get("addressId").Str
	address.Mobile = addressData.Get("mobile").Str
	address.Name = addressData.Get("name").Str
	address.CountryName = addressData.Get("countryName").Str
	address.ProvinceName = addressData.Get("provinceName").Str
	address.CityName = addressData.Get("cityName").Str
	address.DistrictName = addressData.Get("districtName").Str
	address.ReceiverAddress = addressData.Get("receiverAddress").Str
	address.DetailAddress = addressData.Get("detailAddress").Str
	address.IsDefault = addressData.Get("isDefault").Int()
	address.AddressTag = addressData.Get("addressTag").Str
	address.Latitude = addressData.Get("latitude").Str
	address.Longitude = addressData.Get("longitude").Str
	address.CreateTime = addressData.Get("createTimeCreateTime").Str
	address.UpdateTime = addressData.Get("updateTime").Str
	return nil, address
}

func (session *Session) GetAddress() (error, []Address) {
	err, result := session.Request.GET(AddressListAPI)
	if err != nil {
		return err, nil
	}
	var addressList = make([]Address, 0)
	validAddress := result.Get("data.addressList").Array()
	for _, addressData := range validAddress {
		err, address := parseAddress(addressData)
		if err != nil {
			return err, nil
		}
		addressList = append(addressList, address)
	}
	return nil, addressList
}

func (session *Session) SetAddress(address Address) error {
	session.Address = address
	data := SetAddressParam{
		AddressId: session.Address.AddressId,
	}
	dataStr, _ := json.Marshal(data)
	err, _ := session.Request.POST(SetAddressAPI, dataStr)
	return err
}

func (session *Session) ChooseAddress() error {
	err, addressList := session.GetAddress()
	if err != nil {
		return err
	}
	if len(addressList) == 0 {
		return conf.NoValidAddressErr
	}
	for i, addr := range addressList {
		fmt.Printf("[%v] %s %s %s %s %s\n", i, addr.Name, addr.DistrictName, addr.ReceiverAddress, addr.DetailAddress, addr.Mobile)
	}
	index := conf.InputSelect(len(addressList))
	if err = session.SetAddress(addressList[index]); err != nil {
		return err
	}
	return nil
}
