package temp

type ITemp interface {
	Get() error
}

type Temp struct {
}

func (t Temp) Get() error {
	return nil
}

func NewTemp() *Temp {

	return &Temp{}
}

func aaa(t ITemp) {
	t.Get()

}

func main() {

	t := NewTemp()
	aaa(t)
}
