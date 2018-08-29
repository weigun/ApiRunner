package runner

type PIresponseInterface interface {
	GetCode() int
	GetContent() string
	GetErrMsg() string
	GetElapsed() int64
}
