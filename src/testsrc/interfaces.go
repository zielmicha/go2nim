package testsrc

type Duck interface {
	quack() string
}

type DuckImpl struct  {
	name string
}

type SuperDuckImpl struct {
	level int
}

func (d *DuckImpl) quack() string {
	return "quack! " + d.name
}

func (d *SuperDuckImpl) quack() string {
	r := ""
	for i := 0; i < 10; i ++ {
		r += "QUACK!!! "
	}
	return r
}

func getLevel(d Duck) int {
	switch dk := d.(type) {
	case *DuckImpl:
		return -1
	case *SuperDuckImpl:
		return dk.level
	default:
		panic("duck not supported!")
	}
}

func main() {
	a := &DuckImpl{name: "Stefan"}
	b := a.(Duck)
}
