package main

import (
	"bytes"
	"sync"
)

const maxSize = 64 << 10 // 64kb

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	if buf.Cap() > maxSize { // 避免内存泄露
		println("ignore it")
		return
	}
	buf.Reset() // reset 复用
	buffers.Put(buf)
}

func main() {
	PutBuffer(bytes.NewBuffer([]byte{1, 2, 3}))
	if tmp, ok := buffers.Get().(*bytes.Buffer); ok {
		println("Get buffer from pool: ", tmp.Bytes())
		tmp.WriteString("12345678")
		println("After buffer updated: ", tmp.Bytes())
		PutBuffer(tmp)
		println("Refetch buffer from pool: ", buffers.Get().(*bytes.Buffer).Bytes())
		return
	}
	print("Not get buffer from pool")
}
