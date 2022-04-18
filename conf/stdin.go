package conf

import (
	"bufio"
	"fmt"
	"os"
)

func InputSelect(_len int) int {
	var index int
	for true {
		fmt.Println("\n[>] 请输入序号（0, 1, 2...)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("[!] 输入有误：%s!\n", err)
		} else if index >= _len {
			fmt.Println("\n[!] 输入有误：超过最大序号！")
		} else {
			break
		}
	}
	return index
}
