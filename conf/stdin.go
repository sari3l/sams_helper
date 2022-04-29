package conf

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func InputSelect(_len int) int {
	var index int
	for true {
		fmt.Println("\n[>] 请输入序号（0, 1, 2...)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("[!] 输入有误：%s!\n", err)
		} else if index > _len {
			fmt.Println("\n[!] 输入有误：超过最大序号！")
		} else {
			break
		}
	}
	return index
}

func InputIntList(_len int) []int {
	fmt.Println("\n[>] 请输入选择并用英文逗号隔开(,)，不选择直接回车即可：")
	var index []int
	var input string
	stdin := bufio.NewReader(os.Stdin)
	_, err := fmt.Fscanln(stdin, &input)
	if err != nil {
		//fmt.Printf("[!] 输入有误：%s!，将继续执行\n", err)
		return nil
	}
	inputs := strings.Split(input, ",")
	for _, v := range inputs {
		if len(v) == 0 {
			continue
		}
		value, err := strconv.Atoi(v)
		if err != nil {
			fmt.Printf("[!] 解析 %s 错误！\n", v)
			continue
		}
		if value >= _len {
			fmt.Printf("[!] 解析 %d 有误：超过最大序号！\n", value)
			continue
		}
		index = append(index, value)
	}
	return index
}

func OutputBytes(content []byte) {
	f := bufio.NewWriter(os.Stdout)
	defer func(f *bufio.Writer) {
		err := f.Flush()
		if err != nil {
			return
		}
	}(f)
	_, err := f.Write(content)
	if err != nil {
		return
	}
}
