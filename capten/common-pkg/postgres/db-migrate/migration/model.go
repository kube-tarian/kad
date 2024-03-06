package migration

type Mode int

const (
	UP    Mode = Mode(1)
	DOWN  Mode = Mode(2)
	PURGE Mode = Mode(3)
)
