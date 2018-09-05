package report

import (
	utils "ApiRunner/utils"
	_ "encoding/json"
	_ "fmt"
	"log"
	_ "log"
	_ "os"
	_ "path/filepath"
	"strings"
	"sync"
	"time"
)

type CaseNum struct {
	TotalCases int64
	Successes  int64
	Failures   int64
	Errors     int64
}

func (this CaseNum) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type Summary struct {
	Title     string
	StartTime int64
	Duration  int64
	CaseNum
}

func (this Summary) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type Validator struct {
	Check      bool
	Comparator string
	Expect     string
	Actual     string
}

func (this Validator) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type ExecuteDetail struct {
	RequestData  map[string]string
	ResponseData map[string]string
	Validators   []Validator
}

func (this ExecuteDetail) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type Record struct {
	//一个执行结果
	//	CaseSetName string
	Status    bool
	Api       string
	Elapsed   int64 //ms
	Detail    ExecuteDetail
	Traceback string // 捕获到的异常
}

func (this Record) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type Info struct {
	CaseSetName string
	Host        string
}

func (this Info) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type RecordSet struct {
	//执行结果的集合，即一个用例集，都是以用例集为单位的
	Info
	List []Record
	CaseNum
}

func (this RecordSet) Add2Cache(uid uint32) {
	iCachePtr.add(uid, this)
}

type Footer bool

func (this Footer) Add2Cache(uid uint32) {
	iCachePtr.doneChan <- uid
}

type TestReport struct {
	Summary   Summary
	RecordSet RecordSet //用例集记录列表
}

func (this *TestReport) dump() {
	log.Println("=====================Exported========================")
	log.Printf("title:%s\n", this.Summary.Title)
	log.Printf("StartTime:%s\n", time.Unix(this.Summary.StartTime, 0).Format(`20060102_150405`))
	log.Printf("Duration:%s\n", this.Summary.Duration)
	log.Printf("Successes:%s\n", this.Summary.Successes)
	log.Printf("Failures:%s\n", this.Summary.Failures)
	log.Printf("Errors:%s\n", this.Summary.Errors)
	log.Println("--------------------------------------")
	log.Printf("CaseSetName:%s\n", this.RecordSet.CaseSetName)
	log.Printf("Host:%s\n", this.RecordSet.Host)
	for i, r := range this.RecordSet.List {
		log.Printf("Api:%s\n", r.Api)
		log.Printf("Status:%s\n", r.Status)
		log.Printf("Elapsed:%s\n", r.Elapsed)
		for k, v := range r.Detail.RequestData {
			log.Printf("%s => %s\n", k, v)
		}
		for k, v := range r.Detail.ResponseData {
			log.Printf("%s => %s\n", k, v)
		}
		for _, vali := range r.Detail.Validators {
			log.Printf("\tcheck:%s\n", vali.Check)
			log.Printf("\tComparator:%s\n", vali.Comparator)
			log.Printf("\tExpect:%s\n", vali.Expect)
			log.Printf("\tActual:%s\n", vali.Actual)
		}
		log.Printf("Traceback:%s\n", r.Traceback)
		log.Printf("record %d end-----------------------\n", i)
	}
	log.Println("=====================Exported========================")
}

func (this *TestReport) export() {
	//输出html报表
	//只是简单的做模板渲染
	this.dump()

}

func newTestReport() *TestReport {
	return &TestReport{}
}

//type ReportPool struct {
//	reportChan chan TestReport
//}

//func (this *ReportPool) push(report TestReport) {
//	this.reportChan <- report
//}

func export2Html(uid uint32) {
	//依据uid来组装报表
	repo := newTestReport()
	s := iCachePtr.GetSummary(uid)
	log.Println("cn3:", s)
	s.Duration = time.Now().Unix() - s.StartTime
	repo.Summary = s
	//	details := iCachePtr.GetExecuteDetails(uid)
	//		detail.Validators = iCachePtr.getValidatorList(uid)
	//组装每项用例记录
	info := iCachePtr.GetInfo(uid)
	//	recoSet := iCachePtr.GetRecordSet(uid)
	recoSet := RecordSet{}
	recoSet.TotalCases = s.TotalCases
	recoSet.Successes = s.Successes
	recoSet.Failures = s.Failures
	recoSet.Errors = s.Errors
	recoSet.CaseSetName = info.CaseSetName
	recoSet.Host = info.Host
	for _, reco := range iCachePtr.GetRecords(uid) {
		//reco.Detail = details[i]
		recoSet.List = append(recoSet.List, reco)
	}
	repo.RecordSet = recoSet
	log.Println("repo:", repo)
	repo.export()

	//	}
}

type PIcacheInterfance interface {
	Add2Cache(uint32)
}

type trouble struct {
	uid  uint32
	item interface{}
}
type itemCache struct {
	summaryCache       map[uint32]Summary
	validatorCache     map[uint32][]Validator
	executeDetailCache map[uint32][]ExecuteDetail
	recordSetCache     map[uint32]RecordSet
	recordCache        map[uint32][]Record
	infoCache          map[uint32]Info
	doneChan           chan uint32
	itemChan           chan trouble
	lastUid            uint32
	counter            int64
	reportNum          int64
}

func (this *itemCache) __add(tr trouble) {
	//TODO 执行相同用例时，需要将对应的缓存删除掉，防止重复的数据
	this.counter++
	//	uid := this.lastUid
	uid := tr.uid
	item := tr.item
	log.Printf("%d >>>>>>>>>recv:%T,counter:%d", uid, item, this.counter)
	switch item.(type) {
	case Summary:
		if uid != 0 {
			this.summaryCache[uid] = item.(Summary)
			this.reportNum++
		}
	case Validator:
		//		this.validatorCache[uid] = item.(validator)
		if uid != 0 {
			this.validatorCache[uid] = append(this.validatorCache[uid], item.(Validator))
		}
	case ExecuteDetail:
		if uid != 0 {
			this.executeDetailCache[uid] = append(this.executeDetailCache[uid], item.(ExecuteDetail))
		}
	case Record:
		if uid != 0 {
			this.recordCache[uid] = append(this.recordCache[uid], item.(Record))
			//			item := item.(Record)
			//			csName := item.CaseSetName
			//			sc := this.recordCache[uid]
			//			sc[csName] = append(sc[csName], item)
		}
	case RecordSet:
		if uid != 0 {
			this.recordSetCache[uid] = item.(RecordSet)
		}
	case CaseNum:
		if uid != 0 {
			item := item.(CaseNum)
			log.Println("cn1:", item, uid)
			s := this.GetSummary(uid)
			log.Printf("cn1.5:%T\n", s)
			//更新summary
			s.TotalCases += item.TotalCases
			s.Errors += item.Errors
			s.Failures += item.Failures
			s.Successes += item.Successes
			//map的value不可寻址，所以需要先编辑局部变量，然后再重新赋值
			this.summaryCache[uid] = s
			log.Println("cn2:", s, this.GetSummary(uid), uid)
			//TODO 更新对应的recordSet?
		}
	case Info:
		if uid != 0 {
			this.infoCache[uid] = item.(Info)
		}
	case uint32:
		if this.counter%2 == 0 {
			// 做一个校验，防止数据串行有问题
			log.Fatalf("uid counter is not odd!!!counter is %d", this.counter)
		}
		this.lastUid = item.(uint32)
	default:
		log.Printf("not support type for report compment %T\n", item)
	}
}

func (this *itemCache) add(uid uint32, item interface{}) {
	//	log.Println("++++++this *itemCache:", this)
	// TODO 不能依赖管道的顺序
	log.Println("++++++add:", uid, item)
	this.itemChan <- trouble{uid, item}
	//	this.itemChan <- uid
	//	this.itemChan <- item
}

func (this *itemCache) GetSummary(uid uint32) Summary {
	return this.summaryCache[uid]
}

func (this *itemCache) GetValidatorList(uid uint32) []Validator {
	return this.validatorCache[uid]

}

func (this *itemCache) GetExecuteDetails(uid uint32) []ExecuteDetail {
	return this.executeDetailCache[uid]
}

func (this *itemCache) GetRecords(uid uint32) []Record {
	return this.recordCache[uid]
}

func (this *itemCache) GetRecordSet(uid uint32) RecordSet {
	return this.recordSetCache[uid]
}

func (this *itemCache) GetInfo(uid uint32) Info {
	return this.infoCache[uid]
}

func (this *itemCache) removeCache(uid uint32) {
	delete(this.infoCache, uid)
	delete(this.summaryCache, uid)
	delete(this.validatorCache, uid)
	delete(this.executeDetailCache, uid)
	delete(this.recordCache, uid)
	delete(this.recordSetCache, uid)
}

var iCachePtr *itemCache

//var repPoolPtr *ReportPool
var once sync.Once

func InitItemCache() {
	//TODO 这样模式的代码太多，可以重构
	log.Println("InitItemCache....")
	once.Do(func() {
		iCachePtr = &itemCache{make(map[uint32]Summary), make(map[uint32][]Validator), make(map[uint32][]ExecuteDetail), make(map[uint32]RecordSet), make(map[uint32][]Record), make(map[uint32]Info), make(chan uint32, 32), make(chan trouble, 64), 0, 0, 0}
		go func() {
			// 串行化获取报表组件
			for {
				select {
				//如果两个chan都能读，则会随机读取一个，因为是带缓存的chan，应该问题不大
				case uid := <-iCachePtr.doneChan:
					log.Println("ready to export2Html.......")
					iCachePtr.reportNum--
					export2Html(uid)
					iCachePtr.removeCache(uid)
					if iCachePtr.reportNum <= 0 {
						utils.SendSignal()
					}
				case it := <-iCachePtr.itemChan:
					iCachePtr.__add(it)
				}
			}
		}()
		log.Println("InitItemCache....done!", iCachePtr)
	})
	log.Println("after init", iCachePtr)
}

/*
私有函数
*/

func generateFileName(caseSetName string) string {
	//生成报告文件名,格式20180831_123000_CaseSetName.html
	return strings.Join([]string{utils.GetDateTime(), caseSetName}, "_") + `.html`
}
