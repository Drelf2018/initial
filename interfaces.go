package initial

type BeforeInitial1 interface {
	BeforeInitial()
}

type BeforeInitial2 interface {
	BeforeInitial() error
}

type AfterInitial1 interface {
	AfterInitial()
}

type AfterInitial2 interface {
	AfterInitial() error
}
