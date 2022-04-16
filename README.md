# sams_helper

山姆下单助手

修改自： https://github.com/robGoods/sams

## 使用方式

修改 `config.yaml` 文件中以下内容：

```yaml
authToken: "<authToken>"	# 山姆会员token，从 Header 中获取
deviceType: 1			# 1->移动端模拟，2->小程序模拟
deliveryType: 2                 # 1->极速达，2->全城配
ignoreInvalid: true             # 是否忽略无效商品
noticeType: 1			# 0->不通知, 1->bark, 2->mac sound
barkToken: "<barkToken>"	# 若 noticeType 设为 1，需要同时将此参数设置为 Bark 通知 token
```

## 注意

1. `小程序模式`与`移动端模式`没有本质不同，随意选择即可。
2. 如果购物车里同时存在`极速达`、`全城配`商品，若`deliveryType = 2`，则会尝试通过全城配购买极速达货物（可能无货）。
3. `ignoreInvalid`
   - true：在检测有无效商品时仍会尝试提交订单（可能会因为金额减少增加运费）。
   - false：会一直检测等待至商品全部有效才会尝试提交订单。
4. `noticeType = 2` 即 `Mac Sound` 提醒，仅对 OSX 系统有效。

## 更新：

- 2022年04月16日 添加`小程序模拟`模式，修复逻辑问题
- 2022年04月15日 选择地址后自动更新购物车
- 2022年04月14日 修正一些函数问题

## 效果

<img src="https://github.com/sari3l/sams_helper/blob/main/pics/pic_1.jpeg" width="50%"/><br/>
<img src="https://github.com/sari3l/sams_helper/blob/main/pics/pic_2.png" width="50%"/><br/>

## 声明

本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！
