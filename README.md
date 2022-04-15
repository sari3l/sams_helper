# sams_helper

山姆下单助手

修改自： https://github.com/robGoods/sams

## 使用方式

修改 `config.yaml` 文件中以下内容

```yaml
authToken: "<authToken>"	# 山姆会员token，从 Header 中获取
deviceType: 1			# 1->移动端模拟，1->小程序模拟
noticeType: 1			# 0->不通知, 1->bark, 2->mac sound
barkToken: "<barkToken>"	# 若 noticeType 设为 1，需要同时将此参数设置为 Bark 通知 token
```


## 更新：

2022年04月16日 添加`小程序模拟`模式，修复逻辑问题
2022年04月15日 选择地址后自动更新购物车
2022年04月14日 修正一些函数问题

## 声明

本项目仅供学习交流，严禁用作商业行为，特别禁止黄牛加价代抢等！

因违法违规等不当使用导致的后果与本人无关，如有任何问题可联系本人删除！
