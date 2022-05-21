# sams_helper

<img src="https://github.com/sari3l/sams_helper/blob/main/pics/sams_helper.png"/>

借鉴自： https://github.com/robGoods/sams v1.0

## 使用方式

### 0x01 运行

1. 从 [Releases](https://github.com/sari3l/sams_helper/releases) 页面下载对应系统版本。
2. 命令行运行可执行文件（初次运行会自动释放配置文件）。

### 0x02 配置

主要修改 `config.yaml` 文件中以下内容：

```yaml
# 山姆会员 Token
authToken: "74xxxxxxxxxxxx"
# 运行模式
# 1->山姆抢购 | 2->保供抢购
runMode: 1
# 配送方式
# 1->极速达 | 2->全城配
deliveryType: 1
```

其他配置（尤其是通知相关配置）请查看文件内注释自行修改。

## 功能支持 

| config                | 简介                                   |
|-----------------------|--------------------------------------|
| addGoodsFromFileSet   | 自动搜索添加商品（即使无货），允许热加载随时更新             |
| autoFixPurchaseLimit  | 对限购、库存数量不足的商品自动修正                    |
| autoInputSet          | 自动选择支付方式、收货地址、优惠券                    |
| autoShardingForOrder  | 超重订单自动拆分下单（暂未考虑运费、重量最优解）             |
| bruteCapacity         | 运力爆破（非常规爆破），全城配下很好用的功能，增大下单成功率       |
| cartSelectedStateSync | 全平台勾选状态同步，只会下单勾选商品                   |
| moneySet              | 控制单次订单金额上下限，累计金额上限（会影响超重拆分功能）        |
| runUnlimited          | 成功下单后程序不退出持续运行                       |
| sleepTimeSet          | 自定义各步骤休眠时间                           |
| supplySet             | 即时获取保供清单，可强制添加无货保供商品，同时可设置黑白名单限制监控清单 |
| updateStoreForce      | 强制刷新商店，避免店面突然上下线引起的异常（早中抢购建议开启）      |
| proxySet              | 代理设置，方便抓包、调试，或切换异地避免当地高峰网络堵塞         |

## 注意事项

1. `runUnlimited: true`永不停歇模式开启
   - 下单成功会回到购物车检查状态，如果没有设置提醒可能错过文字消息。
   - 保供商品已下单的套餐不会再次购买，在`保供抢购`下不建议程序持续开启超过一天。
2. `runMode: 2` 即`保供自动抢单`模式，可以配置黑白名单进行筛选，具体内容在`config.yaml - supplySet`自行配置
3. `bruteCapacity: true` 开启运力暴力模式 
   - 此模式下主要以尽力下单为目标，最终配送时间跨度可能较大
   - 暴力模式会通过`全城配`来下单，可能导致最终订单内取消部分`极速达`商品（具体可通过查看移动端购物车`全城配`即为最终订单商品列表）。
4. 如果购物车里同时存在`极速达`、`全城配`商品，若`deliveryType = 2`，则会尝试通过全城配购买极速达货物（可能无货）。
5. `ignoreInvalid` 忽略无效商品
   - `true` 在检测有无效商品时自动剔除，继续提交订单（可能会因为金额减少增加运费）。
   - `false` 会一直检测等待至商品全部有效才会尝试提交订单。
6. `autoFixPurchaseLimit` 限购数量自动修正
   - `isEnabled: true` 开启修正
   - `fixOffline: true` 线下购物车修正，单独开启只影响最后提交订单数量
   - `fixOnline: true` 线上购物车修正，修正后会重新获取购物车信息
7. `noticeType = 3` 即 `Mac Sound` 提醒，仅对 OSX 系统有效。

## 主要更新：

- 2022年05月10日 新增`超重自动拆分`设置
- 2022年05月02日 新增`全平台购物车勾选状态同步`设置
- 2022年05月01日 新增`强制添加购物车`、`强制更新商店`设置，新增`下单金额限制`功能
- 2022年04月29日 添加`自动输入`设置
- 2022年04月28日 添加`永不停歇`设置
- 2022年04月26日 添加`优惠券`模块
- 2022年04月25日 `保供商品`筛选修改为黑白名单
- 2022年04月21日 `保供商品`可通过关键字排除
- 2022年04月20日 添加`保供自动抢单`模块
- 2022年04月19日 添加`运力暴力模式`模块
- 2022年04月18日 添加`Server酱`提醒模块
- 2022年04月16日 添加`小程序模拟`模式方式
- 2022年04月15日 选择地址后自动更新购物车
- 2022年04月14日 修正一些函数问题

## 效果

<img src="https://github.com/sari3l/sams_helper/blob/main/pics/pic_1.jpeg" width="50%"/><br/>
<img src="https://github.com/sari3l/sams_helper/blob/main/pics/pic_2.png" width="50%"/><br/>

## 声明

本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！
