# sams_helper

山姆下单助手

修改自： https://github.com/robGoods/sams v1.0

## 使用方式

### 0x01 配置

主要修改 `config.yaml` 文件中以下内容：

```yaml
authToken: "<authToken>"	# 山姆会员token，从 Header 中获取
runMode: 1                      # 1-> 山姆抢购，2->保供抢购
bruteCapacity: false            # 运力爆破模式，开启则会尝试覆盖所有时间
deliveryType: 2                 # 1->极速达，2->全城配
ignoreInvalid: true             # 是否忽略无效商品
noticeType: 0			# 0->不通知, 1->Bark, 2->Server酱，3->OSX 系统语音
```

其他配置（尤其是通知相关配置）请查看文件内注释自行修改

### 0x02 运行方式

#### I. 源代码

命令行运行

```shell
go run main.go
```

#### II. Release

1. 从 [Releases](https://github.com/sari3l/sams_helper/releases) 页面下载对应系统版本。
2. 将 [config.yaml](https://github.com/sari3l/sams_helper/blob/main/config.yaml) 单独保存至与执行文件同级目录下并修改。
   ```plain
   .
   ├── LICENSE
   ├── README.md
   ├── config.yaml
   └── sams_helper
   ```
3. 命令行运行可执行文件。

## 注意

1. `优惠券`模块受限于没有很多优惠券进行测试，可能存在很多BUG，欢迎在issue中提供脱敏的使用优惠券的下单数据包
2. `runMode: 2` 即`保供自动抢单`模式，可以配置黑白名单进行筛选，具体内容在`config.yaml - supplySet`自行配置
3. `bruteCapacity: true` 开启运力暴力模式 
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


### 已知情况

1. 部分地区不释放运力，具体表现为`20:59:00`左右运力列表并未更新，且唯一显示运力为隔天不可用
2. 部分地区保供列表套餐有数量，添加购物车但无法结算，目前认为运力不可达

## 更新：

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
