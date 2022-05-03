package conf

import "errors"

var GotoStoreStep = errors.New("跳转商店展示")
var GotoCartStep = errors.New("跳转购物车获取")
var GotoCartShowStep = errors.New("跳转购物车展示")
var GotoGoodsStep = errors.New("跳转商品校验")
var GotoCapacityStep = errors.New("跳转运力检查")
var GotoOrderStep = errors.New("跳转下订购单")
var GotoExit = errors.New("程序退出")

var RunModeErr = errors.New("运行模式错误，请检查配置")
var AuthTokenErr = errors.New("auth-token 可能未设置，请检查")
var ProxyErr = errors.New("网络错误，请检查是否设置错误代理")
var NoValidAddressErr = errors.New("没有有效的收货地址，请前往 APP 添加或者检查 Auth-Token 是否正确")
var CheckStoreErr = errors.New("检查商店信息失败")
var CheckCartErr = errors.New("检查购物车失败")
var NoGoodsErr = errors.New("当前购物车中无有效商品")
var FixCartErr = errors.New("修正购物车列表限购商品数量失败")
var NoValidCouponErr = errors.New("没有查询到有效的优惠券，跳过此步骤")
var NoChooseCouponErr = errors.New("输入异常或未选取优惠券，将继续执行")
var NoStoreInitErr = errors.New("尚未获取到商店信息，将重新检查")

var MoneyMinErr = errors.New("检测购物车金额尚未达到单次订单金额下限，将重新检查购物车")
var MoneyMaxErr = errors.New("检测购物车金额超过单次订单金额上限，将手动修改购物车")
var TotalLimitErr = errors.New("检测下单总金额已超累计金额上限，程序将退出！")

var AuthFailErr = errors.New("鉴权失败 auth-token 过期")
var CartGoodChangeErr = errors.New("购物车商品发生变化，请返回购物车页面重新结算")
var LimitedErr = errors.New("服务器正忙,请稍后再试")
var LimitedErr1 = errors.New("当前购物火爆，请稍后再试")
var OutOfSellErr = errors.New("部分商品已缺货")
var CapacityFullErr = errors.New("当前无剩余运力，重新检测是否释放")
var FAILErr = errors.New("未知失败")
var GetDeliveryErr = errors.New("获取履约配送信息异常，即将重试")
var NoMatchDeliverMode = errors.New("当前区域不支持配送，将重新读取商店信息")
var CloseOrderTimeExceptionErr = errors.New("尊敬的会员，您选择的配送时间已失效，请重新选择")
var NotDeliverCapCityErr = errors.New("当前配送时间段已约满，请重新选择配送时段")
var DecreaseCapacityCountError = errors.New("扣减运力失败，即将重试")
var StoreHasClosedError = errors.New("门店已打烊")
var DeliveryTypeErr = errors.New("未知设备类型")
var NotCheckShopPendingErr = errors.New("请阅读并勾选《购物须知》")
var RequestErr = errors.New("请求异常")
var CartFullErr = errors.New("购物车最多99件商品")

// 429 429 {"message":"Requests rate limited. stage:service"}
