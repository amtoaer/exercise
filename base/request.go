package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"iwsp/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (s *Session) get(url string) (content string, err error) {
	var resp *http.Response
	if resp, err = s.client.Get(url); err != nil {
		return
	}
	if content, err = readBody(resp); err != nil {
		return
	}
	return content, err
}

func (s *Session) buildMap() {
	s.countMap = make(map[string]int)
	content, err := s.get(s.infoURL)
	if err != nil {
		utils.Fatal(err)
	}
	r := regexp.MustCompile("var timeArr = (.+);")
	m := regexp.MustCompile(`(.+?)\(剩余(.+?)\)`)
	matchResult := r.FindStringSubmatch(content)
	if len(matchResult) < 2 {
		utils.Fatal("正则表达式匹配失败！")
	}
	timeArr := strings.ReplaceAll(matchResult[1], "'", "\"")
	var container map[string]interface{}
	err = json.Unmarshal([]byte(timeArr), &container)
	utils.Log("获取时段剩余人数...")
	for _, value := range container[time.Now().Format("2006-01-02")].([]interface{}) {
		item := value.(string)
		utils.Log(item)
		tmp := m.FindStringSubmatch(item)
		count, err := strconv.Atoi(tmp[2])
		if err != nil {
			utils.Fatal("时段剩余人数获取失败")
		}
		s.countMap[tmp[1]] = count
	}
	utils.Log("时段剩余人数获取成功")
}

// GetOrderList 得到历史预约列表
func (s *Session) GetOrderList() {
	utils.Log("开始请求并解析预约列表")
	content, err := s.get("http://book.neu.edu.cn/booking/page/orderList")
	if err != nil {
		utils.Fatal(err)
	}
	r := regexp.MustCompile("var orderList = (.+);")
	matchResult := r.FindStringSubmatch(content)
	if len(matchResult) < 2 {
		utils.Fatal("正则表达式匹配失败！")
	}
	orderList := strings.ReplaceAll(matchResult[1], "'", "\"")
	var container []map[string]interface{}
	err = json.Unmarshal([]byte(orderList), &container)
	if err != nil {
		utils.Log(err)
		utils.Fatal("预约列表解析失败！")
	}
	utils.Log("预约列表解析成功！")
	output := func(container []map[string]interface{}) {
		getStatus := func(status float64) string {
			switch status {
			case 0:
				return "已预约"
			case 1:
				return "已入场"
			case 2:
				return "已出场"
			case 3:
				return "已完成"
			case 4:
				return "已取消"
			case 5:
				return "已关闭"
			case 6:
				return "已过期"
			default:
				return "未知状态"
			}
		}
		for _, m := range container {
			var (
				position = m["ruleName"].(string)
				date     = m["bookPeriodStartTime"].(string)[0:10]
				duration = m["bookPeriodName"].(string)
				status   = getStatus(m["status"].(float64))
			)
			fmt.Printf("%14s%23s%23s%10s\n", position, date, duration, status)
		}
	}
	fmt.Printf("%14s%20s%20s%14s\n", "地点", "日期", "时段", "状态")
	output(container)
}

// Post post预约信息
func (s *Session) Post() {
	// 构建时段->人数的map
	s.buildMap()
	// 检测时段输入是否正确，人数是否还有剩余
	s.data.Check(s.countMap)
	data, err := json.Marshal(s.data)
	utils.Log("序列化预约信息...")
	if err != nil {
		utils.Fatal("预约信息序列化失败")
	}
	utils.Log("序列化结果:", string(data))
	utils.Log("开始进行post请求...")
	resp, err := s.client.Post(s.createURL, "application/json;charset=UTF-8", bytes.NewBuffer(data))
	if err != nil {
		utils.Fatal("post预约请求失败，请重试")
	}
	result, err := readBody(resp)
	utils.Log("返回结果：", result)
	reqResult := &returnData{}
	if err = json.Unmarshal([]byte(result), reqResult); err != nil {
		utils.Fatal("返回结果异常，预约失败")
	}
	fmt.Println(reqResult.Message)
}

func readBody(resp *http.Response) (string, error) {
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	return string(content), nil
}
