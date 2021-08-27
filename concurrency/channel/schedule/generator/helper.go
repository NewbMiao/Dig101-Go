package generator

func Repeat(
	done <-chan struct{},
	values ...interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func TakeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{}) // 创建输出流
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ { // 只读取前num个元素
			select {
			case <-done:
				return
			case takeStream <- <-valueStream: // 从输入流中读取元素
			}
		}
	}()
	return takeStream
}

func AsStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	s := make(chan interface{}) // 创建一个unbuffered的channel
	go func() {                 // 启动一个goroutine，往s中塞数据
		defer close(s)             // 退出时关闭chan
		for _, v := range values { // 遍历数组
			select {
			case <-done:
				return
			case s <- v: // 将数组元素塞入到chan中
			}
		}
	}()
	return s
}

func GenerateChanStream(num int) <-chan <-chan interface{} {
	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)
		for i := 0; i < num; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}
