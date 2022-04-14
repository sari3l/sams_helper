package main

import (
	"SAMS_buyer/sams"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	// 初始化用户信息
	session := sams.Session{}
	err := session.InitSession(
		"<token>",
		2)

	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	for true {
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
		fmt.Printf("\n########## 开始校验当前商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		if err = session.CheckGoods(); err != nil {
			fmt.Println(err)
			switch err {
			case sams.OOSErr:
				goto CartLoop
			default:
				goto GoodsLoop
			}
		}
		if err = session.CheckSettleInfo(); err != nil {
			fmt.Println(err)
			switch err {
			case sams.CartGoodChangeErr:
				goto CartLoop
			case sams.LimitedErr:
				goto GoodsLoop
			default:
				goto GoodsLoop
			}
		}

	CapacityLoop:
		fmt.Printf("\n########## 获取当前可用配送时间【%s】 ###########\n", time.Now().Format("15:04:05"))
		err = session.CheckCapacity()
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		dateISFull := true
		for _, capCityResponse := range session.Capacity.CapCityResponseList {
			if capCityResponse.DateISFull == false && dateISFull {
				dateISFull = false
				fmt.Printf("发现可用的配送时段:%s!\n", capCityResponse.StrDate)
			}
		}

		if dateISFull {
			fmt.Println("当前无可用配送时间段")
			time.Sleep(1 * time.Second)
			goto CapacityLoop
		}
	OrderLoop:
		err = session.CommitPay()
		fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
		switch err {
		case nil:
			fmt.Println("抢购成功，请前往app付款！")
			return
		case sams.LimitedErr1:
			fmt.Printf("[%s] 立即重试...\n", err)
			goto OrderLoop
		default:
			goto CartLoop
		}
	}

	for _ = range [3]int{} {
		exec.Command("say", "抢到啦 快去付款", "--voice=Ting-ting").Run()
	}
}
