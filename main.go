package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"net"
	"sams_helper/conf"
	"sams_helper/notice"
	"sams_helper/requests"
	"sams_helper/sams"
	"sams_helper/tools"
	"strings"
	"time"
)

func main() {
	err, session := doInitStep()
	if err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}

	go doExtendStep(&session)

	if err = doBuyStep(&session); err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}
}

func doExtendStep(session *sams.Session) {
	if session.Setting.RunMode == 2 {
		go stepSupply(session)
	}

	if session.Setting.UpdateStoreForce {
		//conf.GotoCartStep = conf.GotoStoreStep
		go checkStoreForce(session)
	}

	if session.Setting.AddGoodsFromFileSet.IsEnabled {
		go stepAddGoodsForce(session)
	}
}

func doInitStep() (error, sams.Session) {
	// 初始化设置
	fmt.Println("\n########## 配置文件检查 ##########")
	err, setting := conf.InitSetting()
	if err != nil {
		return err, sams.Session{}
	}

	// 配置检查
	if !(setting.RunMode == 1 || setting.RunMode == 2) {
		return conf.RunModeErr, sams.Session{}
	}

	if setting.AutoShardingForOrder {
		setting.AutoFixPurchaseLimitSet.IsEnabled = true
		setting.AutoFixPurchaseLimitSet.FixOnline = true
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
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepStoreSleep) * time.Millisecond)
		goto stepStoreLoop
	}
stepCartLoop:
	if err := stepCart(session); err != nil {
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCartSleep) * time.Millisecond)
		if session.Setting.UpdateStoreForce {
			goto stepStoreLoop
		} else if err == conf.GotoExit {
			return nil
		} else {
			goto stepCartLoop
		}
	}
stepCartShowLoop:
	if err := stepCartShow(session); err != nil {
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCartShowSleep) * time.Millisecond)
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
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepGoodsSleep) * time.Millisecond)
		if err == conf.GotoStoreStep {
			goto stepStoreLoop
		} else if err == conf.GotoCartStep {
			goto stepCartLoop
		} else if err == conf.NoMatchDeliverMode {
			goto stepStoreLoop
		} else {
			goto stepGoodsLoop
		}
	}
	capacityLoopCount := 0
	bruteTime := 1
stepCapacityLoop:
	if err := stepCapacity(session, bruteTime); err != nil {
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepCapacitySleep) * time.Millisecond)
		if err == conf.GotoStoreStep {
			goto stepStoreLoop
		} else if err == conf.GotoCartStep {
			goto stepCartLoop
		} else if err == conf.CapacityFullErr {
			capacityLoopCount += 1
			if capacityLoopCount >= 20 {
				goto stepStoreLoop
			} else {
				goto stepCapacityLoop
			}
		} else {
			goto stepCapacityLoop
		}
	}
stepOrderLoop:
	if err := stepOrder(session); err != nil {
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepOrderSleep) * time.Millisecond)
		if err == conf.GotoStoreStep {
			goto stepStoreLoop
		} else if err == conf.GotoCartStep {
			goto stepCartLoop
		} else if err == conf.GotoCapacityStep {
			goto stepCapacityLoop
		} else if err == conf.DecreaseCapacityCountErr || err == conf.GetDeliveryErr {
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
	var historyList = map[string]bool{}
	var isChange = false

	for _, v := range session.StoreList {
		historyList[v.StoreId] = true
	}

	if err := session.GetStoreList(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.CheckStoreErr))...)
		tools.OutputBytes(c)
		return err
	}

	c = append(c, []byte(fmt.Sprintf("########## 更新就近商店信息【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	for index, store := range session.StoreList {
		c = append(c, []byte(fmt.Sprintf("[%v] Id：%s 名称：%s, 类型 ：%d\n", index, store.StoreId, store.StoreName, store.StoreType))...)
		if !historyList[store.StoreId] {
			isChange = true
		}
	}
	if isChange {
		tools.OutputBytes(c)
		return nil
	}
	return nil
}

func stepCart(session *sams.Session) error {
	var c []byte
	if session.Setting.MoneySet.TotalCalc > session.Setting.MoneySet.TotalLimit*100 {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.TotalLimitErr))...)
		tools.OutputBytes(c)
		return conf.GotoExit
	}

	// 拆分历史恢复
	if session.Setting.AutoShardingForOrder && len(session.GoodsListFuture) > 0 {
		c = append(c, []byte(fmt.Sprintf("########## 检测到拆包遗留商品列表【%s】 ###########\n", time.Now().Format("15:04:05")))...)
		addGoodsList := make([]sams.AddCartGoods, 0)
		for _, v := range session.GoodsListFuture {
			addGoodsList = append(addGoodsList, v.ToAddCartGoods(v.Quantity))
		}
		if err := session.AddCartGoodsInfo(addGoodsList); err != nil {
			c = append(c, []byte(fmt.Sprintf("[!] 添加拆包遗留商品失败：%v", err))...)
			tools.OutputBytes(c)
			return err
		} else {
			c = append(c, []byte(fmt.Sprintf("[>] 添加拆包遗留商品成功，数量：%d\n", len(session.GoodsListFuture)))...)
			session.GoodsListFuture = make([]sams.Goods, 0)
		}
	}

	if err := session.CheckCart(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.CheckCartErr))...)
		tools.OutputBytes(c)
		return err
	}
	tools.OutputBytes(c)
	return nil
}

func stepCartShow(session *sams.Session) error {
	var c []byte
	var amount int64
	c = append(c, []byte(fmt.Sprintf("########## 获取购物车中商品清单【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	session.GoodsList = make([]sams.Goods, 0)
	for _, v := range session.Cart.FloorInfoList {
		var _amount int64
		if v.FloorId == session.FloorId {
			for index, goods := range v.NormalGoodsList {
				if session.Setting.CartSelectedStateSync && !goods.IsSelected {
					c = append(c, []byte(fmt.Sprintf("[未勾选] %s 数量：%v 单价：%s 重量：%fkg\n", goods.GoodsName, goods.Quantity, tools.SPrintMoney(goods.Price), goods.Weight))...)
					continue
				}
				session.GoodsList = append(session.GoodsList, goods.ToGoods())
				c = append(c, []byte(fmt.Sprintf("[%v] %s 数量：%v 单价：%s 重量：%fkg\n", index, goods.GoodsName, goods.Quantity, tools.SPrintMoney(goods.Price), goods.Weight))...)
				_amount += goods.Quantity * goods.Price
			}
			session.FloorInfo = v
			session.DeliveryInfoVO = sams.DeliveryInfoVO{
				StoreDeliveryTemplateId: v.StoreInfo.StoreDeliveryTemplateId,
				DeliveryModeId:          v.StoreInfo.DeliveryModeId,
				StoreType:               v.StoreInfo.StoreType,
			}
			c = append(c, []byte(fmt.Sprintf("[>] 订单总价：%s\n", tools.SPrintMoney(_amount)))...)
			amount += _amount
		}
		if v.IsOverWeight && !session.Setting.AutoShardingForOrder {
			c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.GoodsOverWeightErr))...)
			tools.OutputBytes(c)
			return conf.GotoCartStep
		}
	}

	if len(session.GoodsList) == 0 {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.NoGoodsErr))...)
		if session.Setting.RunMode == 2 && !session.Setting.SupplySet.ShowCartAlways {
			return conf.GotoCartStep
		}
		tools.OutputBytes(c)
		return conf.GotoCartStep
	}

	if session.Setting.AutoFixPurchaseLimitSet.IsEnabled && (session.Setting.AutoFixPurchaseLimitSet.FixOffline || session.Setting.AutoFixPurchaseLimitSet.FixOnline) {
		err, isChangedOffline, isChangedOnline := session.FixCart()
		if err != nil {
			tools.OutputBytes(c)
			return conf.GotoCartStep
		} else {
			if isChangedOffline && !isChangedOnline {
				c = append(c, []byte(fmt.Sprintln("[>] 已自动修正当前限购数量，不影响线上购物车信息，将继续执行"))...)
				tools.OutputBytes(c)
				return conf.GotoCartShowStep
			}
			if isChangedOnline {
				c = append(c, []byte(fmt.Sprintln("[>] 已自动修正限购数量，将重新检查购物车"))...)
				tools.OutputBytes(c)
				return conf.GotoCartStep
			}
		}
	}

	// 极速达过重分包
	if session.Setting.AutoShardingForOrder && session.Setting.DeliveryType == 1 {
		var isOverWeight = false
		var weightAmount int64
		var delGoodsList []sams.DelCartGoods
		var weightLimit = tools.StringToInt64(session.FloorInfo.WeightThreshold)
		var goodsListTmp = make([]sams.Goods, 0)
		for _, v := range session.GoodsList {
			if isOverWeight {
				session.GoodsListFuture = append(session.GoodsListFuture, v)
				delGoodsList = append(delGoodsList, v.ToDelCartGoods())
			} else {
				for i := int64(1); i <= v.Quantity; i++ {
					weightAmountTmp := weightAmount + i*int64(v.Weight*1000000)
					if weightAmountTmp < weightLimit {
						continue
					} else {
						v2 := sams.Goods{}
						_ = copier.Copy(&v2, &v)
						v2.Quantity -= i - 1
						session.GoodsListFuture = append(session.GoodsListFuture, v2)
						v.Quantity = i - 1
						if v.Quantity > 0 {
							_ = session.ModifyCartGoodsInfo(v)
							goodsListTmp = append(goodsListTmp, v)
						} else {
							delGoodsList = append(delGoodsList, v.ToDelCartGoods())
						}
						isOverWeight = true
						break
					}
				}
				if !isOverWeight {
					goodsListTmp = append(goodsListTmp, v)
				}
			}
			weightAmount += v.Quantity * int64(v.Weight*1000000)
		}
		if isOverWeight {
			c = append(c, []byte(fmt.Sprintf("########## 发现超重订单，执行分包【%s】 ###########\n", time.Now().Format("15:04:05")))...)
			session.GoodsList = goodsListTmp
			for index, v := range session.GoodsList {
				c = append(c, []byte(fmt.Sprintf("[%v] %s 数量：%v 单价：%s\n", index, v.GoodsName, v.Quantity, tools.SPrintMoney(v.Price)))...)
			}
			_ = session.DelCartGoodsInfo(delGoodsList)
		}

		if len(session.GoodsListFuture) > 0 {
			c = append(c, []byte(fmt.Sprintf("########## 分包下批次待购商品【%s】 ###########\n", time.Now().Format("15:04:05")))...)
			for _, goods := range session.GoodsListFuture {
				c = append(c, []byte(fmt.Sprintf("[~] %s 数量：%v 单价：%s\n", goods.GoodsName, goods.Quantity, tools.SPrintMoney(goods.Price)))...)
			}
		}
	}

	if amount < session.Setting.MoneySet.AmountMin*100 {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.MoneyMinErr))...)
		tools.OutputBytes(c)
		return conf.GotoCartStep
	} else if amount > session.Setting.MoneySet.AmountMax*100 {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.MoneyMaxErr))...)
		tools.OutputBytes(c)
		return conf.GotoCartStep
	}

	tools.OutputBytes(c)
	return nil
}

func stepGoods(session *sams.Session) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 开始校验当前商品、优惠券【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	if err := session.CheckGoods(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		tools.OutputBytes(c)
		switch err {
		case conf.OutOfSellErr:
			return conf.GotoCartStep
		default:
			return conf.GotoGoodsStep
		}
	}

	if err := session.CheckSettleInfo(); err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] 校验商品失败：%s\n", err))...)
		tools.OutputBytes(c)
		switch err {
		case conf.CartGoodChangeErr:
			return conf.GotoCartStep
		case conf.LimitedErr:
			return conf.GotoGoodsStep
		case conf.NoMatchDeliverMode:
			return conf.NoMatchDeliverMode
		default:
			return conf.GotoGoodsStep
		}
	} else {
		c = append(c, []byte(fmt.Sprintf("[>] 校验商品成功\n[>] 优惠券抵扣：%s, 最终总金额：%s\n", tools.SPrintMoneyStr(session.SettleInfo.CouponFee), tools.SPrintMoneyStr(session.SettleInfo.TotalAmount)))...)
		session.DeliveryInfoVO = sams.DeliveryInfoVO{
			StoreDeliveryTemplateId: session.SettleInfo.SettleDelivery.StoreDeliveryTemplateId,
			DeliveryModeId:          session.SettleInfo.SettleDelivery.DeliveryModeIdList[0],
			StoreType:               session.Setting.StoreType,
		}
	}
	tools.OutputBytes(c)
	return nil
}

func stepCapacity(session *sams.Session, bruteTime int) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 获取当前可用配送时间【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	err, content := session.CheckCapacity(bruteTime)
	c = append(c, content...)
	if err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		tools.OutputBytes(c)
		switch err {
		case conf.CapacityFullErr:
			return err
		case conf.LimitedErr, conf.LimitedErr1:
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
	tools.OutputBytes(c)
	return nil
}

func stepOrder(session *sams.Session) error {
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	err, order := session.CommitPay()
	if err == nil {
		c = append(c, []byte(fmt.Sprintf("[>] 抢购成功，订单号 %s，请前往app付款！\n", order.OrderNo))...)
		session.Setting.MoneySet.TotalCalc += order.PayInfo.TotalAmt
		tools.OutputBytes(c)
		err = notice.Do(session.Setting.NoticeSet)
		if err != nil {
			c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
			tools.OutputBytes(c)
		}
		if session.Setting.RunUnlimited || len(session.GoodsListFuture) > 0 {
			return conf.GotoCartStep
		} else {
			return nil
		}
	} else {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		switch err {
		case conf.LimitedErr, conf.LimitedErr1:
			c = append(c, []byte(fmt.Sprintln("[!] 立即重试..."))...)
			tools.OutputBytes(c)
			return conf.GotoOrderStep
		case conf.CloseOrderTimeExceptionErr, conf.NotDeliverCapCityErr:
			tools.OutputBytes(c)
			return conf.GotoCapacityStep
		case conf.DecreaseCapacityCountErr:
			tools.OutputBytes(c)
			return conf.DecreaseCapacityCountErr
		case conf.GetDeliveryErr:
			tools.OutputBytes(c)
			return conf.GotoCartStep
		case conf.OutOfSellErr:
			tools.OutputBytes(c)
			return conf.GotoCartStep
		case conf.StoreHasClosedErr:
			tools.OutputBytes(c)
			return conf.GotoStoreStep
		case conf.GoodsExceedLimitErr:
			tools.OutputBytes(c)
			return conf.GotoCartStep
		default:
			tools.OutputBytes(c)
			return conf.GotoCapacityStep
		}
	}
}

func stepSupply(session *sams.Session) {
	orderAlready := make(map[string]bool)
GetSupplyGoodsLoop:
	var trigger = 1
	var c []byte
	c = append(c, []byte(fmt.Sprintf("########## 获取保供商品【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	validGoods := sams.ShowGoods{}
	err, goodsList := session.GetGuaranteedSupplyGoodsAll()
	if err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] 保供监控错误：%s\n", err))...)
		tools.OutputBytes(c)
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepSupplySleep) * time.Millisecond)
		goto GetSupplyGoodsLoop
	} else {
		if len(goodsList) == 0 {
			c = append(c, []byte(fmt.Sprintln("[!] 未上架保供商品"))...)
			tools.OutputBytes(c)
			time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepSupplySleep) * time.Millisecond)
			goto GetSupplyGoodsLoop
		}
	}

	for index, v := range goodsList {
		if orderAlready[v.SpuId] {
			c = append(c, []byte(fmt.Sprintf("[已添加此商品] %s 数量：%v 单价：%s 详情：%s\n", v.Title, v.StockQuantity, tools.SPrintMoney(v.Price), v.SubTitle))...)
			continue
		}
		if session.Setting.SupplySet.ParseSet.IsEnabled {
			isBlack := false
			isWhite := false
			for _, keyWord := range session.Setting.SupplySet.ParseSet.KeyWords {
				if len(keyWord) > 0 && session.Setting.SupplySet.ParseSet.Mode == 2 && strings.Contains(v.Title, keyWord) {
					isBlack = true
					break
				} else if len(keyWord) > 0 && session.Setting.SupplySet.ParseSet.Mode == 1 && (strings.Contains(v.Title, keyWord) || strings.Contains(v.SubTitle, keyWord)) {
					isWhite = true
					break
				}
			}
			if (session.Setting.SupplySet.ParseSet.Mode == 2 && isBlack) || (session.Setting.SupplySet.ParseSet.Mode == 1 && !isWhite) {
				c = append(c, []byte(fmt.Sprintf("[已忽略此商品] %s 数量：%v 单价：%s 详情：%s\n", v.Title, v.StockQuantity, tools.SPrintMoney(v.Price), v.SubTitle))...)
				continue
			}
		}

		c = append(c, []byte(fmt.Sprintf("[%v] %s 数量：%v 单价：%s 详情：%s\n", index, v.Title, v.StockQuantity, tools.SPrintMoney(v.Price), v.SubTitle))...)
		if session.Setting.SupplySet.AddForce || v.StockQuantity > 0 {
			validGoods = v
			if session.Setting.SupplySet.AddForce {
				c = append(c, []byte(fmt.Sprintf("[>] 强制添加保供商品: %s, 后台尝试下单\n", validGoods.Title))...)
			} else {
				c = append(c, []byte(fmt.Sprintf("[+] 发现可购保供商品: %s, 即将自动添加并下单\n", validGoods.Title))...)
			}
			// 自动添加购物车
			_goodList := []sams.AddCartGoods{validGoods.ToNormalGoods().ToAddCartGoods(1)}
			if err = session.AddCartGoodsInfo(_goodList); err != nil {
				c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
			} else {
				trigger += 1
				orderAlready[v.SpuId] = true
			}
		}
	}
	tools.OutputBytes(c)
	time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepSupplySleep*trigger) * time.Millisecond)
	goto GetSupplyGoodsLoop
}

func stepAddGoodsForce(session *sams.Session) {
	var fileMd5 string
	var fileName = tools.GetFilePath("goodsList.yaml")
	var cartGoodsListHistory []sams.ShowGoods
HotStartLoop:
	var c []byte
	var goodsList map[string]int64

	c = append(c, []byte(fmt.Sprintf("########## 获取 goodsList.yaml 中商品描述内容【%s】 ###########\n", time.Now().Format("15:04:05")))...)
	err, newFileMd5 := tools.FileMd5Calc(fileName)
	if err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		tools.OutputBytes(c)
		time.Sleep(1000 * time.Millisecond)
		goto HotStartLoop
	}

	err = tools.ReadFromYaml(fileName, &goodsList)
	if err != nil {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
		tools.OutputBytes(c)
		time.Sleep(1000 * time.Millisecond)
		goto HotStartLoop
	}

	if len(session.StoreList) == 0 {
		c = append(c, []byte(fmt.Sprintf("[!] %s\n", conf.NoStoreInitErr))...)
		tools.OutputBytes(c)
		time.Sleep(1000 * time.Millisecond)
		goto HotStartLoop
	}

	if fileMd5 == newFileMd5 {
		c = append(c, []byte(fmt.Sprintln("[>] 预期商品文件列表未变化"))...)
	} else {
		c = append(c, []byte(fmt.Sprintf("########## 检测待获取商品列表更新【%s】 ###########\n", time.Now().Format("15:04:05")))...)

		if len(cartGoodsListHistory) > 0 {
			c = append(c, []byte(fmt.Sprintln("[>] 删除历史购物车信息"))...)
			delGoodsList := make([]sams.DelCartGoods, 0)
			for _, v := range cartGoodsListHistory {
				delGoodsList = append(delGoodsList, v.ToNormalGoods().ToDelCartGoods())
			}
			if err = session.DelCartGoodsInfo(delGoodsList); err != nil {
				c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
				tools.OutputBytes(c)
				time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepGoodsHotModeSleep) * time.Millisecond)
				goto HotStartLoop
			}
			cartGoodsListHistory = make([]sams.ShowGoods, 0)
		}

		fileMd5 = newFileMd5
		for goodsName, goodsQuantity := range goodsList {
			if goodsQuantity == 0 {
				goodsQuantity = 1
			}
			c = append(c, []byte(fmt.Sprintf("[>] 正在搜索商品关键字：%s\n", goodsName))...)
			err, result := session.GetGoodsFromSearch(goodsName)
			if err != nil {
				c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
				continue
			}
			cartGoodsListHistory = result
			if len(result) == 0 {
				c = append(c, []byte(fmt.Sprintln("[!] 未查询到相关产品信息"))...)
			} else {
				c = append(c, []byte(fmt.Sprintf("[>] 发现商品数量：%d\n", len(result)))...)
				addGoodsList := make([]sams.AddCartGoods, 0)
				for _, v := range result {
					if session.Setting.AddGoodsFromFileSet.ShowGoodsInfo {
						c = append(c, []byte(fmt.Sprintf("[+] 准备添加商品：%s，数量：%v\n", v.Title, goodsQuantity))...)
					}
					addGoodsList = append(addGoodsList, v.ToNormalGoods().ToAddCartGoods(goodsQuantity))
				}
				if err = session.AddCartGoodsInfo(addGoodsList); err != nil {
					c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
				} else {
					c = append(c, []byte(fmt.Sprintln("[>] 添加商品成功"))...)
					for _, v := range result {
						if session.Setting.AddGoodsFromFileSet.ShowGoodsInfo {
							c = append(c, []byte(fmt.Sprintf("[x] 修改商品数量：%s，数量：%v\n",
								v.Title, goodsQuantity))...)
						}

						modifyGoodsQuantity := v.ToNormalGoods().ToGoods()
						modifyGoodsQuantity.Quantity = goodsQuantity
						if err = session.ModifyCartGoodsInfo(modifyGoodsQuantity); err != nil {
							c = append(c, []byte(fmt.Sprintf("[!] %s\n", err))...)
						} else {
							c = append(c, []byte(fmt.Sprintln("[>] 修改商品数量成功"))...)
						}
					}
				}
			}
		}
	}
	tools.OutputBytes(c)
	if session.Setting.AddGoodsFromFileSet.Mode == 2 {
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepGoodsHotModeSleep) * time.Millisecond)
		goto HotStartLoop
	}
	return
}

func checkStoreForce(session *sams.Session) {
	for true {
		time.Sleep(time.Duration(session.Setting.SleepTimeSet.StepUpdateStoreForceSleep) * time.Millisecond)
		_ = stepStore(session)
	}
}
