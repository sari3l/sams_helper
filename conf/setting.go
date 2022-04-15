package conf

import (
	"SAMS_buyer/notice"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Setting struct {
	AuthToken string            `yaml:"authToken"`
	NoticeSet notice.NoticerSet `yaml:"noticeSet"`
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
