package sams

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
)

type CapCityResponse struct {
	StrDate        string `json:"strDate"`
	DeliveryDesc   string `json:"deliveryDesc"`
	DeliveryDescEn string `json:"deliveryDescEn"`
	DateISFull     bool   `json:"dateISFull"`
}

type Capacity struct {
	Data                      string            `json:"data"`
	CapCityResponseList       []CapCityResponse `json:"capcityResponseList"`
	PortalPerformanceTemplate string            `json:"getPortalPerformanceTemplateResponse"`
}

func parseCapacity(result gjson.Result) (error, CapCityResponse) {
	var sizes []map[string]interface{}
	for _, size := range result.Get("sizes").Array() {
		sizes = append(sizes, size.Value().(map[string]interface{}))
	}
	capacity := CapCityResponse{
		StrDate:        result.Get("strDate").Str,
		DeliveryDesc:   result.Get("deliveryDesc").Str,
		DeliveryDescEn: result.Get("deliveryDescEn").Str,
		DateISFull:     result.Get("dateISFull").Bool(),
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

func (session *Session) CheckCapacity() error {
	urlPath := CapacityDataAPI

	data := make(map[string]interface{})
	data["perDateList"] = []string{
		time.Now().Format("2006-01-02"),
		time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		time.Now().AddDate(0, 0, 3).Format("2006-01-02"),
		time.Now().AddDate(0, 0, 4).Format("2006-01-02"),
		time.Now().AddDate(0, 0, 5).Format("2006-01-02"),
		time.Now().AddDate(0, 0, 6).Format("2006-01-02"),
	}
	data["storeDeliveryTemplateId"] = session.Cart.FloorInfoList[0].StoreInfo.StoreDeliveryTemplateId
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
			return session.GetCapacity(result)
		default:
			return errors.New(result.Get("msg").Str)
		}
	} else {
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
