package main
import "fmt"

func main(arg int) (ret1 int) {
    fmt.Println("hello world")
	arg = 5
}

func (s *SectionReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case SeekStart:
		offset += s.base
	}
}

func Seek1() (r1 int64, r2 int64) {
	if true {
		r1 := 5
	}
}
