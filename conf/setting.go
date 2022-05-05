package conf

import (
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
	AddForce       bool           `yaml:"addForce"`
	ParseSet       SupplyParseSet `yaml:"parseSet"`
	ShowCartAlways bool           `yaml:"showCartAlways"`
}

type SleepTimeSet struct {
	StepStoreSleep        int `yaml:"stepStoreSleep"`
	StepCartSleep         int `yaml:"stepCartSleep"`
	StepCartShowSleep     int `yaml:"stepCartShowSleep"`
	StepGoodsSleep        int `yaml:"stepGoodsSleep"`
	StepCapacitySleep     int `yaml:"stepCapacitySleep"`
	StepOrderSleep        int `yaml:"stepOrderSleep"`
	StepSupplySleep       int `yaml:"stepSupplySleep"`
	StepGoodsHotModeSleep int `yaml:"stepGoodsHotModeSleep"`
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

func InitSetting() (error, Setting) {
	setting := Setting{}
	err := tools.ReadFromYaml("config.yaml", &setting)
	if err != nil {
		return err, Setting{}
	}
	return nil, setting
}
