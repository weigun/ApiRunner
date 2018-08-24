type randomService struct {
	accountPool chan int64 //账号随机池
	//numberPool  chan int64 //通用随机池
}

const (
	maxAccoun = 99999999999
)

var once sync.Once
var _randomService *randomService

func NewRandomService() *randomService {
	//随机服务为单例
	once.Do(func() {
		_randomService = &randomService{make(chan int64, 100)}
		go func() {
			for {
				_randomService.generateAccount()
			}
		}()
	})
	return _randomService
}

//func (this *randomService) generateRand(min, max int64) {
//	this.accountPool <- this.getRand(min, max)
//}

func (this *randomService) generateAccount() {
	//账号生成
	this.accountPool <- this.getRand(0, maxAccoun)
}

func (this *randomService) getRand(min, max int64) int64 {
	//获取一个[min,max]的随机数
	if min >= max {
		return max
	} else {
		rand.Seed(time.Now().UnixNano() + rand.Int63())
		return rand.Int63n(max-min) + min
	}
}

func (this *randomService) getAccount() string {
	//获取一个随机账号
	return strconv.FormatInt(<-this.accountPool, 10)
}