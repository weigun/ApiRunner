// report_manager.go
package services

import (
	"ApiRunner/dao"
	"ApiRunner/models"
	"ApiRunner/utils"

	"encoding/json"
	"fmt"
)

/*
报表管理器
测试报表的生成与管理
*/
type ReportManager struct {
	cache dao.Cache
}

var ReportMgr *ReportManager

const DefaultReportKey = `report`

func (rm *ReportManager) Add(r *models.Report) string {
	jsReport := r.Json()
	rid := utils.MD5(jsReport)
	key := fmt.Sprintf(`%s:%s`, DefaultReportKey, rid)
	rm.cache.Put(key, jsReport, 3600*24*30)
	return rid
}

func (rm *ReportManager) GetReport(rid string) *models.Report {
	key := fmt.Sprintf(`%s:%s`, DefaultReportKey, rid)
	jsReport := rm.cache.Get(key)
	if jsReport != `{}` {
		var report models.Report
		if err := json.Unmarshal([]byte(jsReport), &report); err != nil {
			panic(err)
		}
		return &report
	}
	return nil
}

func (rm *ReportManager) Get(rid string) string {
	key := fmt.Sprintf(`%s:%s`, DefaultReportKey, rid)
	jsReport := rm.cache.Get(key)
	return jsReport
}

func (rm *ReportManager) Remove(rid string) error {
	key := fmt.Sprintf(`%s:%s`, DefaultReportKey, rid)
	return rm.cache.Delete(key)
}

func init() {
	ReportMgr = &ReportManager{dao.GetCache()}
}
