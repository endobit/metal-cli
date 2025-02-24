package commands

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

const (
	attribute   = "attr"
	rack        = "rack"
	appliance   = "appliance"
	cluster     = "cluster"
	host        = "host"
	environment = "environment"
	model       = "model"
	zone        = "zone"
)
