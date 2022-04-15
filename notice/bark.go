package notice

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type BarkSet struct {
	Server  string `yaml:"barkServer"`
	Token   string `yaml:"barkToken"`
	Message string `yaml:"barkMessage"`
	Sound   string `yaml:"barkSound"`
}

func BarkPush(barkSet BarkSet) error {
	urlPath := fmt.Sprintf("%s/%s/%s?sound=%s", barkSet.Server, barkSet.Token, barkSet.Message, barkSet.Sound)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(urlPath)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("[%v] %s", resp.StatusCode, body))
	}
}
