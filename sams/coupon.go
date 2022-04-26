package sams

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"sams_helper/conf"
)

type Coupon struct {
	Code        string `json:"code"`
	Remark      string `json:"remark"`
	Name        string `json:"name"`
	ExpireStart string `json:"expireStart"`
	ExpireEnd   string `json:"expireEnd"`
}

func parseCoupon(result gjson.Result) (error, Coupon) {
	coupon := Coupon{}
	coupon.Code = result.Get("code").Str
	coupon.Remark = result.Get("remark").Str
	coupon.Name = result.Get("name").Str
	coupon.ExpireStart = result.Get("expireStart").Str
	coupon.ExpireEnd = result.Get("expireEnd").Str
	return nil, coupon
}

func parseCouponList(result gjson.Result) (error, []Coupon) {
	var couponList []Coupon
	for _, v := range result.Get("data.couponInfoList").Array() {
		_, coupon := parseCoupon(v)
		couponList = append(couponList, coupon)
	}
	return nil, couponList
}

func (session *Session) GetCoupon() (error, []Coupon) {
	var total int64 = 20
	var status = "1"        // 1->有效 | 3->已过期
	var page int64 = 1      // 初始页数
	var pageSize int64 = 10 // 默认 20
	couponList := make([]Coupon, 0)
	for (page-1)*pageSize <= total {
		data := CouponListParam{
			Uid:      session.Uid,
			Status:   status,
			PageSize: pageSize,
			PageNum:  page,
		}
		dataStr, _ := json.Marshal(data)
		err, result := session.Request.POST(CouponListAPI, dataStr)
		if err != nil {
			return err, nil
		}
		total = result.Get("data.total").Int()
		page += 1
		_, couponListTmp := parseCouponList(result)
		couponList = append(couponList, couponListTmp...)
	}
	return nil, couponList
}

func (session *Session) ChooseCoupons() error {
	err, couponList := session.GetCoupon()
	if err != nil {
		return err
	}
	if len(couponList) == 0 {
		return conf.NoValidCouponErr
	}
	for i, addr := range couponList {
		fmt.Printf("[%2v] 名称：%-15s 有效期：%s - %s 简介：%-30q\n", i, addr.Name, conf.UnixToTime(addr.ExpireStart), conf.UnixToTime(addr.ExpireEnd), addr.Remark)
	}

	indexes := conf.InputString(len(couponList))
	if len(indexes) == 0 {
		fmt.Printf("[>] %s\n", conf.NoChooseCouponErr)
	} else {
		fmt.Printf("[>] 已选取优惠券 %v", indexes)
		for _, v := range indexes {
			session.CouponList = append(session.CouponList, couponList[v])
		}
	}
	return nil
}
