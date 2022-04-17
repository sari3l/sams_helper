package conf

import "errors"

var AuthTokenErr = errors.New("auth-token 可能未设置，请检查")
var AuthFailErr = errors.New("鉴权失败 auth-token 过期")
var CartGoodChangeErr = errors.New("购物车商品发生变化，请返回购物车页面重新结算")
var LimitedErr = errors.New("服务器正忙,请稍后再试")
var LimitedErr1 = errors.New("当前购物火爆，请稍后再试")
var OOSErr = errors.New("部分商品已缺货")
var CapacityFullErr = errors.New("当前无剩余运力，重新检测是否释放")
var FAILErr = errors.New("未知失败")
var NoMatchDeliverMode = errors.New("当前区域不支持配送，请重新选择地址")
var CloseOrderTimeExceptionErr = errors.New("尊敬的会员，您选择的配送时间已失效，请重新选择")
var NotDeliverCapCityErr = errors.New("当前配送时间段已约满，请重新选择配送时段")
var DecreaseCapacityCountError = errors.New("扣减运力失败")
var StoreHasClosedError = errors.New("门店已打烊")
var DeliveryTypeErr = errors.New("未知设备类型")
