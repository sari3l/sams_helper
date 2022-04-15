package main

import (
	"SAMS_buyer/conf"
	"SAMS_buyer/notice"
	"SAMS_buyer/requests"
	"SAMS_buyer/sams"
	"fmt"
	"net"
	"time"
)

func main() {
	// 初始化设置
	err, setting := conf.InitSetting()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// 初始化 requests
	request := requests.Request{}
	err = request.InitRequest(setting)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// 初始化用户信息
	session := sams.Session{}
	err = session.InitSession(request, setting)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	for true {
	AddressLoop:
		// 选取收货地址
		fmt.Println("\n########## 切换购物车收货地址 ##########\n")
		err = session.ChooseAddress()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				fmt.Println("\n网络错误，请检查是否设置错误代理!")
				return
			} else {
				goto AddressLoop
			}

		} else {
			fmt.Println("\n切换成功!")
			fmt.Printf("%s %s %s %s %s \n", session.Address.Name, session.Address.DistrictName, session.Address.ReceiverAddress, session.Address.DetailAddress, session.Address.Mobile)
		}

	StoreLoop:
		// 获取门店
		err = session.GetStoreList()
		if err != nil {
			fmt.Printf("%s", err)
			goto StoreLoop
		}

		for index, store := range session.StoreList {
			fmt.Printf("[%v] Id：%s 名称：%s, 类型 ：%d\n", index, store.StoreId, store.StoreName, store.StoreType)
		}

	CartLoop:
		// 商品列表获取，与地址挂钩
		fmt.Printf("\n########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		session.CheckCart()
		for _, v := range session.Cart.FloorInfoList {
			if v.FloorId == session.FloorId {
				for index, goods := range v.NormalGoodsList {
					session.GoodsList = append(session.GoodsList, goods.ToGoods())
					fmt.Printf("[%v] %s 数量：%v 总价：%d\n", index, goods.GoodsName, goods.Quantity, goods.Price)
				}
				session.FloorInfo = v
				session.DeliveryInfoVO = sams.DeliveryInfoVO{
					StoreDeliveryTemplateId: v.StoreInfo.StoreDeliveryTemplateId,
					DeliveryModeId:          v.StoreInfo.DeliveryModeId,
					StoreType:               v.StoreInfo.StoreType,
				}
			} else {
				// 无效商品
				//for index, goods := range v.NormalGoodsList {
				//	fmt.Printf("----[%v] %s 数量：%v 总价：%d\n", index, goods.SpuId, goods.StoreId, goods.Price)
				//}
			}
		}
		if len(session.GoodsList) == 0 {
			fmt.Println("当前购物车中无有效商品")
			goto CartLoop
		}

	GoodsLoop:
		// 商品检查
		fmt.Printf("\n########## 开始校验当前商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err = session.CheckGoods(); err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			switch err {
			case conf.OOSErr:
				goto CartLoop
			default:
				goto GoodsLoop
			}
		}
		if err = session.CheckSettleInfo(); err != nil {
			fmt.Printf("校验商品失败：%s\n", err)
			time.Sleep(1 * time.Second)
			switch err {
			case conf.CartGoodChangeErr:
				goto CartLoop
			case conf.LimitedErr:
				goto GoodsLoop
			case conf.NoMatchDeliverMode:
				goto AddressLoop
			default:
				goto GoodsLoop
			}
		}

	CapacityLoop:
		// 运力获取
		fmt.Printf("\n########## 获取当前可用配送时间【%s】 ###########\n", time.Now().Format("15:04:05"))
		err = session.CheckCapacity()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			goto CapacityLoop
		}

		if session.SettleDeliveryInfo.ArrivalTimeStr != "" {
			fmt.Printf("发现可用的配送时段::%s!\n", session.SettleDeliveryInfo.ArrivalTimeStr)
		} else {
			fmt.Println("当前无可用配送时间段")
			time.Sleep(1 * time.Second)
			goto CapacityLoop
		}

	OrderLoop:
		// 下订单操作
		err, order := session.CommitPay()
		fmt.Printf("\n########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
		switch err {
		case nil:
			fmt.Printf("抢购成功，订单号 %s，请前往app付款！", order.OrderNo)
			err = notice.Do(setting.NoticeSet)
			if err != nil {
				fmt.Printf("%s", err)
			}
			return
		case conf.LimitedErr1:
			fmt.Println("立即重试...")
			goto OrderLoop
		case conf.CloseOrderTimeExceptionErr, conf.DecreaseCapacityCountError, conf.NotDeliverCapCityErr:
			goto CapacityLoop
		case conf.OOSErr:
			goto CartLoop
		case conf.StoreHasClosedError:
			goto StoreLoop
		default:
			goto CapacityLoop
		}
	}
}
