package tools

import (
	"fmt"
	"strconv"
)

func StringToInt64(str string) int64 {
	result, _ := strconv.ParseInt(str, 10, 64)
	return result
}

func Int64ToString(num int64) string {
	result := strconv.FormatInt(num, 10)
	return result
}

func SPrintMoneyStr(money string) string {
	moneyNum := StringToInt64(money)
	return SPrintMoney(moneyNum)
}

func SPrintMoney(money int64) string {
	return fmt.Sprintf("%d.%d", money/100, money%100)
}
