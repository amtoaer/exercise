package base

import "iwsp/utils"

// InitData 通过预约地点得到用于提交的data interface
func (s *Session) InitData(location string) {
	utils.Log("您输入的地点", location)
	switch location {
	case "fycc":
		s.data = new(fycc)
	default:
		utils.Fatal("预约地点输入错误，请检查。")
	}
}

// GetData 返回struct内的data interface
func (s *Session) GetData() interface{} {
	return s.data
}