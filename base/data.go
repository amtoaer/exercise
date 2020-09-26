package base

import "iwsp/utils"

// InitData 通过预约地点得到用于提交的PostContent
func (s *Session) InitData(location string) {
	switch location {
	case "fycc":
		s.data = &fycc{}
		s.infoURL = "http://book.neu.edu.cn/booking/page/rule/13"
	default:
		utils.Fatal("预约地点输入错误，请检查。")
	}
}

// GetData 返回struct内的PostContent
func (s *Session) GetData() PostContent {
	return s.data
}
