package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type (
	boolFlag struct {
		value *bool
	}

	stringFlag struct {
		value *string
	}

	Appliance   struct{ stringFlag }
	Arch        struct{ stringFlag }
	Cluster     struct{ stringFlag }
	Model       struct{ stringFlag }
	Rack        struct{ stringFlag }
	Environment struct{ stringFlag }
	Host        struct{ stringFlag }
	JSON        struct{ boolFlag }
	Make        struct{ stringFlag }
	Rename      struct{ stringFlag }
	Template    struct{ stringFlag }
	TimeZone    struct{ stringFlag }
	Value       struct{ stringFlag }
	Zone        struct{ stringFlag }
)

func (b boolFlag) Val() bool {
	if b.value == nil {
		return false
	}

	return *b.value
}

func (b boolFlag) Ptr() *bool {
	return b.value
}

func (s stringFlag) Val() string {
	if s.value == nil {
		return ""
	}

	return *s.value
}

func (s stringFlag) Ptr() *string {
	return s.value
}

func (b *boolFlag) Add(flags *pflag.FlagSet, object string) {
	b.value = flags.Bool("json", false, "output "+object+" as JSON")
}

func (a *Appliance) Add(flags *pflag.FlagSet, object string) {
	a.value = flags.String("appliance", "", "appliance for the "+object)
}

func (a *Arch) Add(flags *pflag.FlagSet, object string) {
	a.value = flags.String("arch", "", "architecture for the "+object)
}

func (c *Cluster) Add(flags *pflag.FlagSet, object string) {
	c.value = flags.String("cluster", "", "cluster for the "+object)
}

func (e *Environment) Add(flags *pflag.FlagSet, object string) {
	e.value = flags.String("environment", "", "environment for the "+object)
}

func (h *Host) Add(flags *pflag.FlagSet, object string) {
	h.value = flags.String("host", "", "host for the "+object)
}

func (m *Make) Add(flags *pflag.FlagSet, object string) {
	m.value = flags.String("make", "", "make for the "+object)
}

func (m *Model) Add(flags *pflag.FlagSet, object string) {
	m.value = flags.String("model", "", "model for the "+object)
}

func (r *Rack) Add(flags *pflag.FlagSet, object string) {
	r.value = flags.String("rack", "", "rack for the "+object)
}

func (r *Rename) Add(flags *pflag.FlagSet, object string) {
	r.value = flags.String("model", "", "rename the "+object)
}

func (t *Template) Add(flags *pflag.FlagSet, object string) {
	t.value = flags.String("template", "", "template for the "+object)
}

func (t *TimeZone) Add(flags *pflag.FlagSet, object string) {
	t.value = flags.String("timezone", "", "time zone for the "+object)
}

func (v *Value) Add(flags *pflag.FlagSet, object string) {
	v.value = flags.String("value", "", "value of the "+object)
}

func (z *Zone) Add(flags *pflag.FlagSet, object string) {
	z.value = flags.String("zone", "", "zone for the "+object)
}

func (a *Appliance) Required(flags *pflag.FlagSet) {
	required(flags, "appliance")
}

func (c *Cluster) Required(flags *pflag.FlagSet) {
	required(flags, "cluster")
}

func (e *Environment) Required(flags *pflag.FlagSet) {
	required(flags, "environment")
}

func (m *Make) Required(flags *pflag.FlagSet) {
	required(flags, "make")
}

func (m *Model) Required(flags *pflag.FlagSet) {
	required(flags, "model")
}

func (r *Rack) Required(flags *pflag.FlagSet) {
	required(flags, "rack")
}

func (z *Zone) Required(flags *pflag.FlagSet) {
	required(flags, "zone")
}

func required(flags *pflag.FlagSet, name string) {
	if err := cobra.MarkFlagRequired(flags, name); err != nil {
		panic(err)
	}
}
