package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
	"time"
)

type CapCityResponse struct {
	StrDate        string `json:"strDate"`
	DeliveryDesc   string `json:"deliveryDesc"`
	DeliveryDescEn string `json:"deliveryDescEn"`
	DateISFull     bool   `json:"dateISFull"`
	List           []List `json:"list"`
}

type List struct {
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	TimeISFull    bool   `json:"timeISFull"`
	Disabled      bool   `json:"disabled"`
	CloseDate     string `json:"closeDate"`
	CloseTime     string `json:"closeTime"`
	StartRealTime string `json:"startRealTime"` //1649984400000
	EndRealTime   string `json:"endRealTime"`   //1650016800000
}

type Capacity struct {
	Data                      string            `json:"data"`
	CapCityResponseList       []CapCityResponse `json:"capcityResponseList"`
	PortalPerformanceTemplate string            `json:"getPortalPerformanceTemplateResponse"`
}

func parseCapacity(result gjson.Result) (error, CapCityResponse) {
	var list []List
	for _, v := range result.Get("list").Array() {
		list = append(list, List{
			StartTime:     v.Get("startTime").Str,
			EndTime:       v.Get("endTime").Str,
			TimeISFull:    v.Get("timeISFull").Bool(),
			Disabled:      v.Get("disabled").Bool(),
			CloseDate:     v.Get("closeDate").Str,
			CloseTime:     v.Get("closeTime").Str,
			StartRealTime: v.Get("startRealTime").Str,
			EndRealTime:   v.Get("endRealTime").Str,
		})
	}
	capacity := CapCityResponse{
		StrDate:        result.Get("strDate").Str,
		DeliveryDesc:   result.Get("deliveryDesc").Str,
		DeliveryDescEn: result.Get("deliveryDescEn").Str,
		DateISFull:     result.Get("dateISFull").Bool(),
		List:           list,
	}
	return nil, capacity
}

func (session *Session) GetCapacity(result gjson.Result) error {
	var capCityResponseList []CapCityResponse
	for _, v := range result.Get("data.capcityResponseList").Array() {
		_, product := parseCapacity(v)
		capCityResponseList = append(capCityResponseList, product)
	}
	session.Capacity = Capacity{
		Data:                      result.String(),
		CapCityResponseList:       capCityResponseList,
		PortalPerformanceTemplate: result.Get("data.getPortalPerformanceTemplateResponse").Str,
	}
	return nil
}

func (session *Session) SetCapacity() error {
	session.SettleDeliveryInfo = SettleDeliveryInfo{}
	isSet := false
	for _, caps := range session.Capacity.CapCityResponseList {
		for _, v := range caps.List {
			fmt.Printf("配送时间： %s %s - %s, 是否可用：%v\n", caps.StrDate, v.StartTime, v.EndTime, !v.TimeISFull && !v.Disabled)
			if v.TimeISFull == false && v.Disabled == false && session.SettleDeliveryInfo.ExpectArrivalTime == "" {
				session.SettleDeliveryInfo.ExpectArrivalTime = v.StartRealTime
				session.SettleDeliveryInfo.ExpectArrivalEndTime = v.EndRealTime
				isSet = true
				break
			}
		}
		if isSet {
			session.SettleDeliveryInfo.DeliveryType = session.Setting.DeliveryType
			session.SettleDeliveryInfo.DeliveryName = caps.StrDate
			break
		}
	}
	if isSet {
		return nil
	}
	return conf.CapacityFullErr
}

func (session *Session) CheckCapacity() error {
	data := make(map[string]interface{})
	var perDateList []string
	for i := 0; i <= session.Setting.PerDateLen; i++ {
		perDateList = append(perDateList, time.Now().AddDate(0, 0, i).Format("2006-01-02"))
	}
	data["perDateList"] = perDateList
	data["storeDeliveryTemplateId"] = session.Cart.FloorInfoList[0].StoreInfo.StoreDeliveryTemplateId
	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(CapacityDataAPI, dataStr)
	if err != nil {
		return nil
	}

	err = session.GetCapacity(result)
	if err != nil {
		return err
	}

	err = session.SetCapacity()
	if err != nil {
		return err
	}
	return nil
}
