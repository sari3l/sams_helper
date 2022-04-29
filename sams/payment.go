package sams

import (
	"fmt"
	"sams_helper/conf"
)

func (session *Session) ChoosePayment() error {
	var payMethod = map[int]string{0: "微信", 1: "支付宝", 2: "银联", 3: "沃尔玛礼品卡"}
	var index = 0
	fmt.Println("########## 选择支付方式 ##########")
	fmt.Println("选择说明：\n[0] 微信\n[1] 支付宝\n[2] 银联\n[3] 沃尔玛礼品卡")
	if session.Setting.AutoInputSet.IsEnabled && session.Setting.AutoInputSet.InputPayMethod <= 3 && session.Setting.AutoInputSet.InputPayMethod >= 0 {
		fmt.Printf("[>] 自动输入：%d\n", session.Setting.AutoInputSet.InputPayMethod)
		index = session.Setting.AutoInputSet.InputPayMethod
	} else {
		if session.Setting.AutoInputSet.IsEnabled {
			fmt.Println("[!] 自动输入开启，但解析 InputPayMethod 错误，请手动输入或检查配置")
		}
		index = conf.InputSelect(4)
	}
	switch index {
	case 0:
		session.Channel = "wechat"
		session.SubSaasId = "1486659732"
	case 1:
		session.Channel = "alipay"
		session.SubSaasId = "2088041177108960"
	case 2:
		session.Channel = "china_unionpay"
		session.SubSaasId = "802440354110771"
	case 3:
		session.Channel = "20200920"
		session.SubSaasId = "contract"
	}
	fmt.Printf("[>] 已选支付方式：%s\n", payMethod[index])
	return nil
}
