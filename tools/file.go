package tools

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

func ReadFromYaml(path string, target any) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, target)
	if err != nil {
		return err
	}
	return nil
}

func FileMd5Calc(path string) (error, string) {
	f, err := os.Open(path)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err = io.Copy(md5hash, f); err != nil {
		return err, ""
	}

	return nil, fmt.Sprintf("%x", md5hash.Sum(nil))
}
