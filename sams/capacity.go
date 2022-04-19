package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
	"strconv"
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

func unixToTime(timestamp string) string {
	_time, _ := strconv.ParseInt(timestamp, 0, 64)
	tm := time.Unix(_time/1000, _time%1000)
	return fmt.Sprintf("%s", tm.Format("2006-01-02 03:04:05 PM"))
}

func (session *Session) SetCapacity() error {
	session.SettleDeliveryInfo = SettleDeliveryInfo{}
	isSet := false
	bruteStatus := false
	if session.Setting.BruteCapacity && session.FloorInfo.StoreInfo.StoreType == 2 {
		var _end []string
		session.SettleDeliveryInfo.DeliveryType = session.Setting.DeliveryType
		session.SettleDeliveryInfo.DeliveryName = session.Capacity.CapCityResponseList[0].StrDate
		session.SettleDeliveryInfo.ExpectArrivalTime = session.Capacity.CapCityResponseList[0].List[0].StartRealTime
		for _, caps := range session.Capacity.CapCityResponseList {
			for _, v := range caps.List {
				_end = append(_end, v.EndRealTime)
			}
		}
		if len(_end) >= 3 {
			session.SettleDeliveryInfo.ExpectArrivalEndTime = _end[len(_end)-2]
			fmt.Printf("爆破配送时间范围：%s - %s\n", unixToTime(session.SettleDeliveryInfo.ExpectArrivalTime), unixToTime(session.SettleDeliveryInfo.ExpectArrivalEndTime))
			bruteStatus = true
			isSet = true
		}
	}
	if !bruteStatus {
		for _, caps := range session.Capacity.CapCityResponseList {
			for _, v := range caps.List {
				fmt.Printf("配送时间： %s %s - %s, 是否可用：%v\n", caps.StrDate, v.StartTime, v.EndTime, !v.TimeISFull && !v.Disabled)
				if v.TimeISFull == false && v.Disabled == false && !isSet {
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
	}
	if isSet {
		return nil
	}
	return conf.CapacityFullErr
}

func (session *Session) CheckCapacity() error {
	var dataStr []byte
	var perDateList []string
	for i := 0; i <= session.Setting.PerDateLen; i++ {
		perDateList = append(perDateList, time.Now().AddDate(0, 0, i).Format("2006-01-02"))
	}
	_data := CapacityDataParam{}
	_data.PerDateList = perDateList
	_data.StoreDeliveryTemplateId = session.Cart.FloorInfoList[0].StoreInfo.StoreDeliveryTemplateId
	switch session.Setting.DeviceType {
	case 2:
		data := MiniProgramCapacityDataParam{
			CapacityDataParam: _data,
			Uid:               session.Uid,
			AppId:             "",
			SassId:            session.Setting.SassId,
		}
		dataStr, _ = json.Marshal(data)
	default: // ios
		data := IOSCapacityDataParam{
			CapacityDataParam: _data,
		}
		dataStr, _ = json.Marshal(data)
	}
	err, result := session.Request.POST(CapacityDataAPI, dataStr)
	if err != nil {
		return nil
	}

	if err = session.GetCapacity(result); err != nil {
		return err
	}

	if err = session.SetCapacity(); err != nil {
		return err
	}
	return nil
}
