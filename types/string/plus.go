package main

/*
# literal string concat in compile
$ go tool compile -m  plus.go             
plus.go:21:6: can inline plus
plus.go:23:11: s1 + "yxz" does not escape
plus.go:24:28: s1 + "y" + s1 + "z" + s1 does not escape
plus.go:25:33: s1 + "y" + s1 + "z" + s1 + "z" does not escape

# 2-5: concatstring2-concatstring5  
#  >5: concatstrings
$ go tool compile -S  plus.go|grep concat
        0x0068 00104 (plus.go:20)       CALL    runtime.concatstring2(SB)
        0x00eb 00235 (plus.go:21)       CALL    runtime.concatstring5(SB)
        0x01e1 00481 (plus.go:22)       CALL    runtime.concatstrings(SB)
        rel 105+4 t=8 runtime.concatstring2+0
        rel 236+4 t=8 runtime.concatstring5+0
        rel 482+4 t=8 runtime.concatstrings+0
*/
func plus() {
	s1 := "x"
	s2 := s1 + "y" + "x" + "z"
	s3 := s1 + "y" + s1 + "z" + s1
	s4 := s1 + "y" + s1 + "z" + s1 + "z"
	println(s2, s3, s4)
}
