package news

type Executor struct {
	Fn func(news chan<- News)
	Ch chan News
}

func (f *Executor) Run() {
	go f.Fn(f.Ch)
}
