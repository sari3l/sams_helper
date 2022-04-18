package sams

import (
	"fmt"
	"sams_helper/conf"
)

func (session *Session) ChoosePayment() error {
	fmt.Println("########## 选择支付方式 ##########")
	fmt.Println("选择说明：\n[0] 微信\n[1] 支付宝\n[2] 银联\n[3] 沃尔玛礼品卡")
	index := conf.InputSelect(4)
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
	return nil
}
