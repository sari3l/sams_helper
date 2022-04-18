package notice

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FTQQSet struct {
	Server  string `yaml:"ftqqServer"`
	SendKey string `yaml:"ftqqSendKey"`
	Title   string `yaml:"ftqqTitle"`
	Desp    string `yaml:"ftqqDesp"`
}

// FTQQPush Server酱，具体查看：https://sct.ftqq.com/
func FTQQPush(ftqqSet FTQQSet) error {
	urlPath := fmt.Sprintf("%s/%s.send", ftqqSet.Server, ftqqSet.SendKey)
	data := url.Values{
		"text": []string{ftqqSet.Title},
		"desp": []string{ftqqSet.Desp},
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(
		urlPath,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
