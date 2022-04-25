package main

import (
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
	// 初始化设置
	err, setting := conf.InitSetting()
	if err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}

	// 初始化 requests
	request := requests.Request{}
	if err = request.InitRequest(setting); err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}

	// 初始化用户信息
	fmt.Println("########## 初始化用户信息 ##########")
	session := sams.Session{}
	if err = session.InitSession(request, setting); err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}

	// 设置支付方式
	if err = session.ChoosePayment(); err != nil {
		fmt.Printf("[!] %s\n", err)
		return
	}
AddressLoop:
	// 选取收货地址
	fmt.Println("########## 切换购物车收货地址 ##########")
	if err = session.ChooseAddress(); err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("[!] %s\n", conf.ProxyErr)
			return
		} else {
			goto AddressLoop
		}
	} else {
		fmt.Println("[>] 切换成功!")
		fmt.Printf("[>] %s %s %s %s %s\n", session.Address.Name, session.Address.DistrictName, session.Address.ReceiverAddress, session.Address.DetailAddress, session.Address.Mobile)
	}

ModeEnter:
	if session.Setting.RunMode == 1 || session.Setting.RunMode == 99 {
	StoreLoop:
		// 获取门店
		fmt.Printf("########## 获取就近商店信息 ###########\n")
		if err = session.GetStoreList(); err != nil {
			fmt.Printf("[!] %s\n", conf.NoGoodsErr)
			time.Sleep(1 * time.Second)
			goto StoreLoop
		}

		for index, store := range session.StoreList {
			fmt.Printf("[%v] Id：%s 名称：%s, 类型 ：%d\n", index, store.StoreId, store.StoreName, store.StoreType)
		}

	CartLoop:
		// 商品列表获取，与地址挂钩
		if err = session.CheckCart(); err != nil {
			fmt.Printf("[!] %s\n", conf.NoGoodsErr)
			time.Sleep(1 * time.Second)
			goto CartLoop
		}
	CartShowLoop:
		fmt.Printf("########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		session.GoodsList = make([]sams.Goods, 0)
		_supplyCheck := false
		for _, v := range session.Cart.FloorInfoList {
			if v.FloorId == session.FloorId {
				for index, goods := range v.NormalGoodsList {
					if session.Setting.RunMode == 99 && goods.SpuId == session.Setting.GoodSpuId {
						_supplyCheck = true
					}
					session.GoodsList = append(session.GoodsList, goods.ToGoods())
					fmt.Printf("[%v] %s 数量：%v 单价：%d.%d\n", index, goods.GoodsName, goods.Quantity, goods.Price/100, goods.Price%100)
				}
				session.FloorInfo = v
				session.DeliveryInfoVO = sams.DeliveryInfoVO{
					StoreDeliveryTemplateId: v.StoreInfo.StoreDeliveryTemplateId,
					DeliveryModeId:          v.StoreInfo.DeliveryModeId,
					StoreType:               v.StoreInfo.StoreType,
				}
				_amount, _ := strconv.ParseInt(session.FloorInfo.Amount, 10, 64)
				fmt.Printf("[>] 订单总价：%d.%d\n", _amount/100, _amount%100)
			}
		}

		if session.Setting.RunMode == 99 && !_supplyCheck {
			session.Setting.RunMode = 2
			goto ModeEnter
		}
		if len(session.GoodsList) == 0 {
			fmt.Printf("[!] %s\n", conf.NoGoodsErr)
			time.Sleep(1 * time.Second)
			goto CartLoop
		}

		if session.Setting.AutoFixPurchaseLimitSet.IsEnabled && (session.Setting.AutoFixPurchaseLimitSet.FixOffline || session.Setting.AutoFixPurchaseLimitSet.FixOnline) {
			err, isChangedOffline, isChangedOnline := session.FixCart()
			if err != nil {
				goto CartLoop
			} else {
				if isChangedOffline && !isChangedOnline {
					fmt.Println("[>] 已自动修正当前限购数量，不影响线上购物车信息，将继续执行")
					goto CartShowLoop
				}
				if isChangedOnline {
					fmt.Println("[>] 已自动修正限购数量，将重新检查购物车")
					goto CartLoop
				}
			}
		}

	GoodsLoop:
		// 商品检查
		fmt.Printf("########## 开始校验当前商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err = session.CheckGoods(); err != nil {
			fmt.Printf("[!] %s\n", err)
			switch err {
			case conf.OutOfSellErr:
				goto CartLoop
			default:
				time.Sleep(500 * time.Millisecond)
				goto GoodsLoop
			}
		}
		if err = session.CheckSettleInfo(); err != nil {
			fmt.Printf("[!] 校验商品失败：%s\n", err)
			switch err {
			case conf.CartGoodChangeErr:
				goto CartLoop
			case conf.LimitedErr:
				time.Sleep(500 * time.Millisecond)
				goto GoodsLoop
			case conf.NoMatchDeliverMode:
				goto AddressLoop
			default:
				goto GoodsLoop
			}
		}

	CapacityLoop:
		// 运力获取
		fmt.Printf("########## 获取当前可用配送时间【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err = session.CheckCapacity(); err != nil {
			fmt.Printf("[!] %s\n", err)
			switch err {
			case conf.CapacityFullErr, conf.LimitedErr, conf.LimitedErr1:
				time.Sleep(500 * time.Millisecond)
				goto CapacityLoop
			default:
				goto CartLoop
			}
		}
		if session.Setting.BruteCapacity && session.FloorInfo.StoreInfo.StoreType == 2 {
			fmt.Printf("[>] 准备爆破提交可配送时段\n")
		} else {
			fmt.Printf("[>] 已自动选择第一条可用的配送时段\n")
		}

	OrderLoop:
		// 下订单操作
		fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
		err, order := session.CommitPay()
		if err == nil {

			fmt.Printf("抢购成功，订单号 %s，请前往app付款！", order.OrderNo)
			err = notice.Do(setting.NoticeSet)
			if err != nil {
				fmt.Printf("[!] %s\n", err)
			}
			return
		} else {
			fmt.Printf("[!] %s\n", err)
			switch err {
			case conf.LimitedErr, conf.LimitedErr1:
				time.Sleep(100 * time.Millisecond)
				fmt.Println("[!] 立即重试...")
				goto OrderLoop
			case conf.CloseOrderTimeExceptionErr, conf.DecreaseCapacityCountError, conf.NotDeliverCapCityErr:
				goto CapacityLoop
			case conf.OutOfSellErr:
				goto CartLoop
			case conf.StoreHasClosedError:
				goto StoreLoop
			default:
				goto CapacityLoop
			}
		}
	} else if session.Setting.RunMode == 2 {
	GetGoodsLoop:
		fmt.Printf("########## 获取保供商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		validGoods := sams.NormalGoodsV2{}
		err, goodsList := session.GetGuaranteedSupplyGoods()
		if err != nil {
			fmt.Printf("[!] 保供监控错误：%s\n", err)
			time.Sleep(1 * time.Second)
			goto GetGoodsLoop
		} else {
			if len(goodsList) == 0 {
				fmt.Println("[!] 未上架保供商品")
				time.Sleep(1 * time.Second)
				goto GetGoodsLoop
			}
		}

		for index, v := range goodsList {
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
					fmt.Printf("[已忽略此商品] %s 数量：%v 单价：%d.%d 详情：%s\n", v.Title, v.StockQuantity, v.Price/100, v.Price%100, v.SubTitle)
					continue
				}
			}

			fmt.Printf("[%v] %s 数量：%v 单价：%d.%d 详情：%s\n", index, v.Title, v.StockQuantity, v.Price/100, v.Price%100, v.SubTitle)
			if v.StockQuantity > 0 {
				validGoods = v
				break
			}
		}

		if validGoods.Title == "" {
			time.Sleep(1 * time.Second)
			goto GetGoodsLoop
		} else {
			fmt.Printf("[>] 发现可购保供商品: %s, 即将自动添加并下单\n", validGoods.Title)
			// 自动添加购物车
			_goodList := []sams.AddCartGoods{validGoods.ToAddCartGoods()}
			if err = session.AddCartGoodsInfo(_goodList); err != nil {
				fmt.Printf("[!] %s\n", err)
			}
			session.Setting.GoodSpuId = validGoods.SpuId
			session.Setting.RunMode = 99
			goto ModeEnter
		}

	} else {
		fmt.Printf("[!] %s\n", conf.RunModeErr)
	}
}
