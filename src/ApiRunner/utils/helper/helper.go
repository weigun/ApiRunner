package helper

func Apply(times int, f func(t int)) {
	for i := 0; i < times; i++ {
		f(i)
	}
}
