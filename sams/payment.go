package sams

import (
	"bufio"
	"fmt"
	"os"
)

func (session *Session) ChoosePayment() error {
	var index int
	fmt.Println("\n########## 选择支付方式 ##########\n")
	for true {
		fmt.Println("选择说明：\n0 微信\n1 支付宝\n2 银联\n3 沃尔玛礼品卡\n\n请输入支付方式序号：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index == 0 {
			session.Channel = "wechat"
			session.SubSaasId = "1486659732"
			break
		} else if index == 1 {
			session.Channel = "alipay"
			session.SubSaasId = "2088041177108960"
			break

		} else if index == 2 {
			session.Channel = "china_unionpay"
			session.SubSaasId = "802440354110771"
			break

		} else if index == 3 {
			session.Channel = "20200920"
			break
		} else {
			fmt.Println("输入有误：序号无效！")
		}
	}
	return nil
}
