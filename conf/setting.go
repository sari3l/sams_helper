package conf

import (
	"SAMS_buyer/notice"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ProxySet struct {
	IsEnabled bool   `yaml:"isEnabled"`
	ProxyUrl  string `yaml:"proxyUrl"`
}

type Setting struct {
	AuthToken     string            `yaml:"authToken"`
	DeviceType    int64             `yaml:"deviceType"`
	DeliveryType  int64             `yaml:"deliveryType"`
	FloorId       int64             `yaml:"floorId"`
	IgnoreInvalid bool              `yaml:"ignoreInvalid"`
	ProxySet      ProxySet          `yaml:"proxySet"`
	NoticeSet     notice.NoticerSet `yaml:"noticeSet"`
}

func InitSetting() (error, Setting) {
	setting := Setting{}
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return err, setting
	}

	err = yaml.Unmarshal(file, &setting)
	if err != nil {
		return err, setting
	}
	return nil, setting
}
