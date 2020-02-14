# debugger

> See: ['Advanced Debugging Techniques of Go Code' by Andrii Soldatenko](https://www.youtube.com/watch?v=2kjmLQY8RJk)

## dlv

use in docker

```go
docker build -t dlvapp debugger/dlv

docker run --rm -it --security-opt="apparmor=unconfined" --cap-add=SYS_PTRACE -v "$(pwd):/go/src/app" dlvapp bash
```


## gdb
gdb 8.3.1
solve mac issue

- create gdb-cert in System keychain
follow [gdb-PermissionsDarwin](https://sourceware.org/gdb/wiki/PermissionsDarwin)

- let gdb understand go's dwarf
gdb pretty print
`go build -ldflags=-compressdwarf=false -gcflags=all="-N -l" -o hello test.go`

> inside gdb will use `$GOROOT/src/runtime/runtime-gdb.py` for loading Go Runtime support.
> which u can see:  `strings hello |grep gdb`
> xxx/src/runtime/runtime-gdb.py

- libsystem_darwin.dylib error
Just ignore it, see [GDB giving weird errors](https://stackoverflow.com/questions/58125727/gdb-giving-weird-errors)
