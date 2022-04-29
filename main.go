package main

import (
	"errors"
	"fmt"
	"net"
	"sams_helper/conf"
	"sams_helper/notice"
	"sams_helper/requests"
	"sams_helper/sams"
	"strconv"
	"strings"
	"time"
)

func main() {
	err, session := doInitStep()
	if err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}

	if session.Setting.RunMode == 2 {
		go stepSupply(&session)
	}

	if err = doBuyStep(&session); err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}
}

func doInitStep() (error, sams.Session) {
	// 初始化设置
	err, setting := conf.InitSetting()
	if err != nil {
		return err, sams.Session{}
	}

	if !(setting.RunMode == 1 || setting.RunMode == 2) {
		return conf.RunModeErr, sams.Session{}
	}
	// 初始化 requests
	request := requests.Request{}
	if err = request.InitRequest(setting); err != nil {
		return err, sams.Session{}
	}

	// 初始化用户信息
	fmt.Println("########## 初始化用户信息 ##########")
	session := sams.Session{}
	if err = session.InitSession(request, setting); err != nil {
		return err, sams.Session{}
	}

	// 设置支付方式
	if err = session.ChoosePayment(); err != nil {
		return err, sams.Session{}
	}

	// 选择收货地址
	if err = stepAddress(&session); err != nil {
		return err, sams.Session{}
	}

	// 选择优惠券
	if err = stepCoupon(&session); err != nil {
		return err, sams.Session{}
	}

	return nil, session
}

func doBuyStep(session *sams.Session) error {
stepStoreLoop:
	if err := stepStore(session); err != nil {
		goto stepStoreLoop
	}
stepCartLoop:
	if err := stepCart(session); err != nil {
		goto stepCartLoop
	}
stepCartShowLoop:
	if err := stepCartShow(session); err != nil {
		if err == conf.GotoStoreStep {
			goto stepStoreLoop
		} else if err == conf.GotoCartStep {
			goto stepCartLoop
		} else {
			goto stepCartShowLoop
		}
	}
stepGoodsLoop:
	if err := stepGoods(session); err != nil {
		if err == conf.GotoStoreStep {
			goto stepStoreLoop
		} else if err == conf.GotoCartStep {
			goto stepCartLoop
		} else if err == conf.NoMatchDeliverMode {
			return err
		} else {
			goto stepGoodsLoop
		}
	}
	capacityLoopCount := 0
	bruteTime := 1
stepCapacityLoop:
	if err := stepCapacity(session, bruteTime); err != nil {
		if err == conf.GotoCartStep {
			goto stepCartLoop
		} else if err == conf.CapacityFullErr {
			capacityLoopCount += 1
			if capacityLoopCount >= 100 {
				goto stepCartLoop
			} else {
				goto stepCapacityLoop
			}
		} else {
			goto stepCapacityLoop
		}
	}
stepOrderLoop:
	if err := stepOrder(session); err != nil {
		if err == conf.GotoStoreStep {
			goto stepStoreLoop
		} else if err == conf.GotoCartStep {
			goto stepCartLoop
		} else if err == conf.GotoCapacityStep {
			goto stepCapacityLoop
		} else if err == conf.DecreaseCapacityCountError {
			bruteTime += 1
			goto stepCapacityLoop
		} else {
			goto stepOrderLoop
		}
	}

	return nil
}

func stepAddress(session *sams.Session) error {
AddressLoop:
	fmt.Println("########## 切换购物车收货地址 ##########")
	if err := session.ChooseAddress(); err != nil {
		if _, ok := err.(net.Error); ok {
			return errors.New(fmt.Sprintf("[!] %s\n", conf.ProxyErr))
		} else {
			goto AddressLoop
		}
	}
	return nil
}

func stepCoupon(session *sams.Session) error {
	fmt.Println("########## 选择使用优惠券 ##########")
	if err := session.ChooseCoupons(); err != nil {
		fmt.Printf("[!] %s\n", err)
	}
	return nil
}

func stepStore(session *sams.Session) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 获取就近商店信息【%s】 ###########\n", time.Now().Format("15:04:05")))...)

	if err := session.GetStoreList(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.NoGoodsErr))...)
		conf.OutputBytes(c)
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepStoreSleep) * time.Millisecond)
		return err
	}

	for index, store := range session.StoreList {
		c = append(c, []byte(fmt.Sprintf("[%v] Id：%s 名称：%s, 类型 ：%d\n", index, store.StoreId, store.StoreName, store.StoreType))...)
	}
	conf.OutputBytes(c)
	return nil
}

func stepCart(session *sams.Session) error {
	var c []byte
	if err := session.CheckCart(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.NoGoodsErr))...)
		conf.OutputBytes(c)
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCartSleep) * time.Millisecond)
		return err
	}
	return nil
}
func stepCartShow(session *sams.Session) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	session.GoodsList = make([]sams.Goods, 0)
	for _, v := range session.Cart.FloorInfoList {
		if v.FloorId == session.FloorId {
			for index, goods := range v.NormalGoodsList {
				session.GoodsList = append(session.GoodsList, goods.ToGoods())
				c = append(c, []byte(fmt.Sprintf("[%v] %s 数量：%v 单价：%d.%d\n", index, goods.GoodsName, goods.Quantity, goods.Price/100, goods.Price%100))...)
			}
			session.FloorInfo = v
			session.DeliveryInfoVO = sams.DeliveryInfoVO{
				StoreDeliveryTemplateId: v.StoreInfo.StoreDeliveryTemplateId,
				DeliveryModeId:          v.StoreInfo.DeliveryModeId,
				StoreType:               v.StoreInfo.StoreType,
			}
			_amount, _ := strconv.ParseInt(session.FloorInfo.Amount, 10, 64)
			c = append(c, []byte(fmt.Sprintf("[>] 订单总价：%d.%d\n", _amount/100, _amount%100))...)
		}
	}

	if len(session.GoodsList) == 0 {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.NoGoodsErr))...)
		if session.Setting.RunMode != 2 {
			conf.OutputBytes(c)
		}
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCartShowSleep) * time.Millisecond)
		return conf.GotoCartStep
	}

	if session.Setting.AutoFixPurchaseLimitSet.IsEnabled && (session.Setting.AutoFixPurchaseLimitSet.FixOffline || session.Setting.AutoFixPurchaseLimitSet.FixOnline) {
		err, isChangedOffline, isChangedOnline := session.FixCart()
		if err != nil {
			conf.OutputBytes(c)
			return conf.GotoCartStep
		} else {
			if isChangedOffline && !isChangedOnline {
				c = append(c, []byte(fmt.Sprintln("[>] 已自动修正当前限购数量，不影响线上购物车信息，将继续执行"))...)
				conf.OutputBytes(c)
				return conf.GotoCartShowStep
			}
			if isChangedOnline {
				c = append(c, []byte(fmt.Sprintln("[>] 已自动修正限购数量，将重新检查购物车"))...)
				conf.OutputBytes(c)
				return conf.GotoCartStep
			}
		}
	}
	conf.OutputBytes(c)
	return nil
}

func stepGoods(session *sams.Session) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 开始校验当前商品【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	if err := session.CheckGoods(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		conf.OutputBytes(c)
		switch err {
		case conf.OutOfSellErr:
			return conf.GotoCartStep
		default:
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepGoodsSleep) * time.Millisecond)
			return conf.GotoGoodsStep
		}
	}

	if err := session.CheckSettleInfo(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] 校验商品失败：%s\n", err))...)
		conf.OutputBytes(c)
		switch err {
		case conf.CartGoodChangeErr:
			return conf.GotoCartStep
		case conf.LimitedErr:
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepGoodsSleep) * time.Millisecond)
			return conf.GotoGoodsStep
		case conf.NoMatchDeliverMode:
			return conf.NoMatchDeliverMode
		default:
			return conf.GotoGoodsStep
		}
	}
	return nil
}

func stepCapacity(session *sams.Session, bruteTime int) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 获取当前可用配送时间【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	err, content := session.CheckCapacity(bruteTime)
	c = append(c, content...)
	if err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		conf.OutputBytes(c)
		switch err {
		case conf.CapacityFullErr:
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCapacitySleep) * time.Millisecond)
			return err
		case conf.LimitedErr, conf.LimitedErr1:
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCapacitySleep) * time.Millisecond)
			return conf.GotoCapacityStep
		default:
			return conf.GotoCartStep
		}
	}

	if session.Setting.BruteCapacity && session.FloorInfo.StoreInfo.StoreType == 2 {
		c = append(c, []byte(fmt.Sprintln("[>] 准备爆破提交可配送时段"))...)
	} else {
		c = append(c, []byte(fmt.Sprintln("[>] 已自动选择第一条可用的配送时段"))...)
	}
	conf.OutputBytes(c)
	return nil
}

func stepOrder(session *sams.Session) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	err, order := session.CommitPay()
	if err == nil {
		c = append(c, []byte(fmt.Sprintf("[>] 抢购成功，订单号 %s，请前往app付款！", order.OrderNo))...)
		conf.OutputBytes(c)
		err = notice.Do(session.Setting.NoticeSet)
		if err != nil {
			c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
			conf.OutputBytes(c)
		}
		if session.Setting.RunUnlimited {
			return conf.GotoCartStep
		} else {
			return nil
		}
	} else {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		switch err {
		case conf.LimitedErr, conf.LimitedErr1:
			c = append(c, []byte(fmt.Sprintln("[!] 立即重试..."))...)
			conf.OutputBytes(c)
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepOrderSleep) * time.Millisecond)
			return conf.GotoOrderStep
		case conf.CloseOrderTimeExceptionErr, conf.NotDeliverCapCityErr:
			conf.OutputBytes(c)
			return conf.GotoCapacityStep
		case conf.DecreaseCapacityCountError:
			conf.OutputBytes(c)
			return conf.DecreaseCapacityCountError
		case conf.OutOfSellErr:
			conf.OutputBytes(c)
			return conf.GotoCartStep
		case conf.StoreHasClosedError:
			conf.OutputBytes(c)
			return conf.GotoStoreStep
		default:
			conf.OutputBytes(c)
			return conf.GotoCapacityStep
		}
	}
}

func stepSupply(session *sams.Session) {
	orderAlready := make(map[string]bool)
GetGoodsLoop:
	var trigger = 1
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 获取保供商品【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	validGoods := sams.NormalGoodsV2{}
	err, goodsList := session.GetGuaranteedSupplyGoods()
	if err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] 保供监控错误：%s\n", err))...)
		conf.OutputBytes(c)
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepSupplySleep) * time.Millisecond)
		goto GetGoodsLoop
	} else {
		if len(goodsList) == 0 {
			c = append(c, []byte(fmt.Sprintln("[!] 未上架保供商品"))...)
			conf.OutputBytes(c)
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepSupplySleep) * time.Millisecond)
			goto GetGoodsLoop
		}
	}

	for index, v := range goodsList {
		if orderAlready[v.SpuId] {
			c = append(c, []byte(fmt.Sprintf("[已添加此商品] %s 数量：%v 单价：%d.%d 详情：%s\n", v.Title, v.StockQuantity, v.Price/100, v.Price%100, v.SubTitle))...)
			continue
		}
		if session.Setting.SupplySet.IsEnabled {
			isBlack := false
			isWhite := false
			for _, keyWord := range session.Setting.SupplySet.KeyWords {
				if len(keyWord) > 0 && session.Setting.SupplySet.Mode == 2 && strings.Contains(v.Title, keyWord) {
					isBlack = true
					break
				} else if len(keyWord) > 0 && session.Setting.SupplySet.Mode == 1 && strings.Contains(v.Title, keyWord) {
					isWhite = true
					break
				}
			}
			if (session.Setting.SupplySet.Mode == 2 && isBlack) || (session.Setting.SupplySet.Mode == 1 && !isWhite) {
				c = append(c, []byte(fmt.Sprintf("[已忽略此商品] %s 数量：%v 单价：%d.%d 详情：%s\n", v.Title, v.StockQuantity, v.Price/100, v.Price%100, v.SubTitle))...)
				continue
			}
		}

		c = append(c, []byte(fmt.Sprintf("[%v] %s 数量：%v 单价：%d.%d 详情：%s\n", index, v.Title, v.StockQuantity, v.Price/100, v.Price%100, v.SubTitle))...)
		if v.StockQuantity > 0 {
			validGoods = v
			c = append(c, []byte(fmt.Sprintf("[>] 发现可购保供商品: %s, 即将自动添加并下单\n", validGoods.Title))...)
			// 自动添加购物车
			_goodList := []sams.AddCartGoods{validGoods.ToAddCartGoods()}
			if err = session.AddCartGoodsInfo(_goodList); err != nil {
				c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
			} else {
				trigger += 1
				orderAlready[v.SpuId] = true
			}
		}
	}
	conf.OutputBytes(c)
	time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepSupplySleep*trigger) * time.Millisecond)
	goto GetGoodsLoop
}
