package cli

//go:generate gotip tool "github.com/dmarkham/enumer" -type Verb -transform lower -text

type Verb int

const (
	Add Verb = iota
	Dump
	List
	Load
	Remove
	Report
	Set
)
