package app

type Error int

const (
	InputError Error = iota + 1
	DBError
)
