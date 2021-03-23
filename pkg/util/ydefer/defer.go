package ydefer

//注册全局defer，中间需要defer 的操作可以注册，最后程序终止时运行defer clean
var (
	globalDefers = NewStack()
)

func Register(fns ...func() error) {
	globalDefers.Push(fns...)
}
func Clean() {
	globalDefers.Clean()
}
