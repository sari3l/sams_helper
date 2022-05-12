package tools

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func GetCurrentDirectory() string {
	currentDir, _ := os.Executable()
	exPath := filepath.Dir(currentDir)
	return exPath
}

func GetFilePath(filename string) string {
	if CheckFileExists(filename) {
		return fmt.Sprintf("./%s", filename)
	}
	exPath := GetCurrentDirectory()
	return path.Join(exPath, filename)
}

func InitFile(filename string, fileContent string) error {
	configPath := GetFilePath(filename)
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fileContent)
	if err != nil {
		return err
	}
	return nil
}
