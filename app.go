package main

import (
	"github.com/emacsist/alfred3/utils"
	"os"
	"log"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

type Bank struct {
	Name          string
	Huoqi         float64
	ThreeMon      float64
	SixMon        float64
	OneYear       float64
	TwoYear       float64
	ThreeYear     float64
	FiveYear      float64
	ZeroOneYear   float64
	ZeroThreeYear float64
	ZeroFiveYear  float64
}

func main() {
	query := utils.GetQuery()

	inputData := strings.Fields(query)

	banName := ""
	money := 0.0
	if len(inputData) == 2 {
		banName = strings.TrimSpace(inputData[0])
		money = stringToFloat(strings.TrimSpace(inputData[1]))
	} else {
		money = stringToFloat(strings.TrimSpace(query))
	}

	alfredResponse := utils.NewAlfredResponse()

	lines := Get()
	for i, line := range lines {
		if i == 0 || strings.HasPrefix(line, "#") {
			continue
		}


		data := strings.Fields(line)
		if len(data) < 11 {
			continue
		}
		if len(banName) > 0 {
			if strings.Contains(data[0], banName) {
				addItem(data, alfredResponse, money)
			}
		} else {
			addItem(data, alfredResponse, money)
		}

	}

	alfredResponse.WriteOutput()
}

func addItem(data []string, alfredResponse *utils.AlfredResponse, money float64) {
	var ban Bank
	ban.Name = data[0]
	ban.Huoqi = stringMoneyRateToFloat(data[1])
	ban.ThreeMon = stringMoneyRateToFloat(data[2])
	ban.SixMon = stringMoneyRateToFloat(data[3])
	ban.OneYear = stringMoneyRateToFloat(data[4])
	ban.TwoYear = stringMoneyRateToFloat(data[5])
	ban.ThreeYear = stringMoneyRateToFloat(data[6])
	ban.FiveYear = stringMoneyRateToFloat(data[7])
	ban.ZeroOneYear = stringMoneyRateToFloat(data[8])
	ban.ZeroThreeYear = stringMoneyRateToFloat(data[9])
	ban.ZeroFiveYear = stringMoneyRateToFloat(data[10])
	alfredResponse.AddDefaultItem(ban.Name + " 活期三个月 " + toString(money, ban.Huoqi/4, ban.Huoqi))
	alfredResponse.AddDefaultItem(ban.Name + " 活期六个月 " + toString(money, ban.Huoqi/2, ban.Huoqi))
	alfredResponse.AddDefaultItem(ban.Name + " 活期一年 " + toString(money, ban.Huoqi, ban.Huoqi))
	alfredResponse.AddDefaultItem(ban.Name + " 定存三个月 " + toString(money, ban.ThreeMon/4, ban.ThreeMon))
	alfredResponse.AddDefaultItem(ban.Name + " 定存六个月 " + toString(money, ban.SixMon/2, ban.SixMon))
	alfredResponse.AddDefaultItem(ban.Name + " 定存一年 " + toString(money, ban.OneYear, ban.OneYear))
	alfredResponse.AddDefaultItem(ban.Name + " 定存二年 " + toString(money, ban.TwoYear*2, ban.TwoYear))
	alfredResponse.AddDefaultItem(ban.Name + " 定存三年 " + toString(money, ban.ThreeYear*3, ban.ThreeYear))
	alfredResponse.AddDefaultItem(ban.Name + " 定存五年 " + toString(money, ban.FiveYear*5, ban.FiveYear))
	alfredResponse.AddDefaultItem(ban.Name + " 0定存一年 " + toString(money, ban.ZeroOneYear, ban.ZeroOneYear))
	alfredResponse.AddDefaultItem(ban.Name + " 0定存三年 " + toString(money, ban.ZeroThreeYear*3, ban.ZeroThreeYear))
	alfredResponse.AddDefaultItem(ban.Name + " 0定存五年 " + toString(money, ban.ZeroFiveYear*5, ban.ZeroFiveYear))
}

func toString(money float64, percent float64, yearRate float64) string {
	result := money * percent
	return fmt.Sprintf("可得 : %.3f 元(年利率 %.3f%%)", result, yearRate*100)
}

func stringMoneyRateToFloat(data string) float64  {
	return stringToFloat(data) / 100
}

func stringToFloat(data string) float64 {
	result, _ := strconv.ParseFloat(data, 64)
	//因为输入的是 百分比的数，所以，这里还要 /100
	return result
}

func Get() []string {
	var data []string
	file, err := os.Open("./money.rate.txt")
	if err != nil {
		data = append(data, err.Error())
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data = append(data, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		data = append(data, err.Error())
	}
	return data
}
