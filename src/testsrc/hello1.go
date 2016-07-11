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
