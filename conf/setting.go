package conf

import (
	"errors"
	"fmt"
	"sams_helper/notice"
	"sams_helper/tools"
)

type ProxySet struct {
	IsEnabled bool   `yaml:"isEnabled"`
	ProxyUrl  string `yaml:"proxyUrl"`
}

type AutoFixPurchaseLimitSet struct {
	IsEnabled  bool `yaml:"isEnabled"`
	FixOffline bool `yaml:"fixOffline"`
	FixOnline  bool `yaml:"fixOnline"`
}

type SupplyParseSet struct {
	IsEnabled bool     `yaml:"isEnabled"`
	Mode      int      `yaml:"mode"`
	KeyWords  []string `yaml:"keyWords"`
}

type SupplySet struct {
	OnlySupply     bool           `yaml:"onlySupply"`
	AddForce       bool           `yaml:"addForce"`
	ShowCartAlways bool           `yaml:"showCartAlways"`
	ParseSet       SupplyParseSet `yaml:"parseSet"`
}

type SleepTimeSet struct {
	StepStoreSleep            int `yaml:"stepStoreSleep"`
	StepCartSleep             int `yaml:"stepCartSleep"`
	StepCartShowSleep         int `yaml:"stepCartShowSleep"`
	StepGoodsSleep            int `yaml:"stepGoodsSleep"`
	StepCapacitySleep         int `yaml:"stepCapacitySleep"`
	StepOrderSleep            int `yaml:"stepOrderSleep"`
	StepSupplySleep           int `yaml:"stepSupplySleep"`
	StepGoodsHotModeSleep     int `yaml:"stepGoodsHotModeSleep"`
	StepUpdateStoreForceSleep int `yaml:"stepUpdateStoreForceSleep"`
}

type AutoInputSet struct {
	IsEnabled       bool  `yaml:"isEnabled"`
	InputPayMethod  int   `yaml:"inputPayMethod"`
	InputAddress    int   `yaml:"inputAddress"`
	InputCouponList []int `yaml:"inputCouponList"`
}

type MoneySet struct {
	AmountMin  int64 `yaml:"amountMin"`
	AmountMax  int64 `yaml:"amountMax"`
	TotalLimit int64 `yaml:"totalLimit"`
	TotalCalc  int64
}

type AddGoodsFromFileSet struct {
	IsEnabled     bool `yaml:"isEnabled"`
	Mode          int  `yaml:"mode"`
	ShowGoodsInfo bool `yaml:"showGoodsInfo"`
}

type Setting struct {
	AuthToken               string                  `yaml:"authToken"`
	RunMode                 int                     `yaml:"runMode"`
	SupplySet               SupplySet               `yaml:"supplySet"`
	BruteCapacity           bool                    `yaml:"bruteCapacity"`
	UpdateStoreForce        bool                    `yaml:"updateStoreForce"`
	SleepTimeSet            SleepTimeSet            `yaml:"sleepTimeSet"`
	DeviceType              int64                   `yaml:"deviceType"`
	DeliveryType            int64                   `yaml:"deliveryType"`
	StoreType               int64                   `yaml:"storeType"`
	FloorId                 int64                   `yaml:"floorId"`
	IgnoreInvalid           bool                    `yaml:"ignoreInvalid"`
	AutoFixPurchaseLimitSet AutoFixPurchaseLimitSet `yaml:"autoFixPurchaseLimit"`
	PerDateLen              int                     `yaml:"perDateLen"`
	SassId                  string                  `yaml:"sassId"`
	ProxySet                ProxySet                `yaml:"proxySet"`
	NoticeSet               notice.NoticerSet       `yaml:"noticeSet"`
	RunUnlimited            bool                    `yaml:"runUnlimited"`
	AutoInputSet            AutoInputSet            `yaml:"autoInputSet"`
	MoneySet                MoneySet                `yaml:"moneySet"`
	AddGoodsFromFileSet     AddGoodsFromFileSet     `yaml:"addGoodsFromFileSet"`
	CartSelectedStateSync   bool                    `yaml:"cartSelectedStateSync"`
	AutoShardingForOrder    bool                    `yaml:"autoShardingForOrder"`
}

const configContent = "###############################\n# 必选 基础设置\n###############################\n# 山姆会员 Token\nauthToken: \"<authToken>\"\n# 运行模式\n# 1->山姆抢购 | 2->保供抢购\nrunMode: 1\n# 配送方式\n# 1->极速达 | 2->全城配\ndeliveryType: 1\n\n###############################\n# 可选 基础设置\n###############################\n# 运力暴力模式（仅全城配有效）\nbruteCapacity: false\n# 强制更新商店信息，每隔 stepUpdateStoreForceSleep 循环执行\n# 注：早中时间抢购建议开启\nupdateStoreForce: false\n# 全平台购物车勾选状态同步\n# 注意：1. 同一件商品，极速达全城配相互切换时，若商品有货则默认勾选\n#      2. 其他情景未知\ncartSelectedStateSync: false\n\n###############################\n# 下单金额限制，非必须不需修改\n###############################\nmoneySet:\n  # 单次订单金额下限（元）\n  amountMin: 0\n  # 单次订单金额上限（元）\n  amountMax: 10000\n  # 累计金额上限（元）\n  totalLimit: 10000\n\n###############################\n# 默认配置，一般不需修改\n###############################\n# 设备类型模拟\n# 1->IOS | 2->小程序\ndeviceType: 1\n# 是否忽略无效商品\nignoreInvalid: true\n# 商店类型（仅占位暂无实际作用）\n# 2->生鲜直达 | 4->极速达 | 8->全球购保税\nstoreType: 2\n# 商品类型\n# 1->普通商品 | 2->全球购商品\nfloorId: 1\n# 运力搜索时间跨度（日）\nperDateLen : 7\n# SaasId\nsaasId: \"1818\"\n\n###############################\n# 通知功能\n###############################\n# 通知设置\nnoticeSet:\n  # 通知类型\n  # 0->不通知 | 1->Bark | 2->Server酱 | 3->OSX 系统语音\n  noticeType: 0\n\n  # Bark 通知配置\n  bark:\n    # Bark 端点地址\n    barkServer: \"https://api.day.app\"\n    # Bark Token\n    barkToken: \"<barkToken>\"\n    # 信息内容\n    barkMessage: \"【山姆助手】抢购成功，请及时付款！\"\n    # 提示音\n    barkSound: \"telegraph\"\n\n  # Server酱 通知配置\n  ftqq:\n    # Server酱 端点地址\n    ftqqServer: \"https://sctapi.ftqq.com\"\n    # 用户 Token，获取方式：https://sct.ftqq.com/sendkey\n    ftqqSendKey: \"<ftqqSendKey>\"\n    # 消息通道：最多支持两通道如：\"9|66\"，具体通道请看 https://sct.ftqq.com/sendkey\n    # 默认：方糖服务号\n    ftqqChannel: \"9\"\n    # 信息标题\n    ftqqTitle: \"【山姆助手】提醒\"\n    # 信息内容\n    ftqqDesp: \"抢购成功，请及时付款！\"\n\n  # OSX 系统语音配置\n  sound:\n    # 语音信息\n    soundMessage: \"山姆抢到啦 快去付款\"\n    # 重复播报次数\n    soundTimes: 2\n    # 播报声音\n    soundVoice: \"Ting-ting\"\n\n###############################\n# 高级功能\n###############################\n# 永不停歇（下单成功不退出）\nrunUnlimited: false\n\n# 超重订单自动分批下单（仅极速达有效）\nautoShardingForOrder: false\n\n# （精准）自动搜索、强制添加指定商品\n# 注：通过此功能无货商品可提前加入购物车\naddGoodsFromFileSet:\n  # 是否开启\n  isEnabled: false\n  # 加载方式，冷加载只执行一次，热加载每隔 stepGoodsHotModeSleep 循环执行\n  # 1->冷加载 | 2->热加载\n  mode: 1\n  # 显示商品细节\n  showGoodsInfo: true\n\n# 保供抢购关键字黑白名单筛选\nsupplySet:\n  # 只关注保供套餐\n  onlySupply: false\n  # 不管是否有库存数量，强制添加进购物车\n  addForce: false\n  # 总是输出购物车信息（不建议开启）\n  showCartAlways: false\n  # 黑白名单筛选\n  parseSet:\n    # 是否开启\n    isEnabled: false\n    # 筛选方式\n    # 1->白名单 | 2->黑名单\n    mode: 1\n    # 检测`套餐名`、`套餐内容`有以下关键字，例：[\"蔬菜\", \"保供套餐C\", \"瑞士卷\"]\n    keyWords: [\"<keyWord>\"]\n\n# 代理设置\nproxySet:\n  # 是否开启代理\n  isEnabled: false\n  # 代理 URL\n  proxyUrl: \"http://127.0.0.1:8080\"\n\n# 自动修正限购商品数量\nautoFixPurchaseLimit:\n  # 是否开启\n  isEnabled: true\n  # 是否开启线下修正，只影响最后提交订单内容\n  fixOffline: true\n  # 是否开启线上修正，会影响线上购物车商品的数量\n  fixOnline: false\n\n# 自动输入设置\nautoInputSet:\n  # 是否开启\n  isEnabled: false\n  # 支付方式序号\n  # 0->微信 | 1->支付宝 | 2->银联 | 3->沃尔玛礼品卡\n  inputPayMethod: 0\n  # 地址序号\n  inputAddress: 0\n  # 优惠券序号列表，例：[0,1,2]\n  inputCouponList: []\n\n# 各步骤休眠时间（毫秒）\nsleepTimeSet:\n  # 请求商店信息\n  stepStoreSleep: 1000\n  # 请求购物车信息\n  stepCartSleep: 1000\n  # 购物车信息展示\n  stepCartShowSleep: 1000\n  # 检查商品有效性\n  stepGoodsSleep: 1000\n  # 请求运力信息\n  stepCapacitySleep: 1000\n  # 下单间隔\n  stepOrderSleep: 100\n  # 请求保供商品信息\n  stepSupplySleep: 1000\n  # 热加载睡眠时间\n  stepGoodsHotModeSleep: 10000\n  # 强制更新商店间隔\n  stepUpdateStoreForceSleep: 5000\n"
const goodsListContent = "# 商品名称: 数量\n小胡鸭 柠檬酸辣去骨凤爪 580g: 2\nMM 进口全脂纯牛奶 200ml*24: 1"

func InitSetting() (error, Setting) {
	currentDir := tools.GetCurrentDirectory()
	filesCache := map[string][]string{
		"config.yaml":    {"配置文件", configContent},
		"goodsList.yaml": {"待购商品文件", goodsListContent},
	}
	for file, content := range filesCache {
		if !tools.CheckFileExists(tools.GetFilePath(file)) {
			fmt.Printf("[!] 未检查到 %s %s，准备释放于 %s\n", content[0], file, currentDir)
			if err := tools.InitFile(file, content[1]); err != nil {
				return errors.New(fmt.Sprintf("[!] 释放 config.yaml 异常，请手动复制 https://github.com/sari3l/sams_helper/blob/main/%s\n", file)), Setting{}
			}
		}
	}

	setting := Setting{}
	filePath := tools.GetFilePath("config.yaml")
	err := tools.ReadFromYaml(filePath, &setting)
	if err != nil {
		return err, Setting{}
	}
	if len(setting.AuthToken) < 64 {
		return AuthTokenErr, Setting{}
	}
	return nil, setting
}
