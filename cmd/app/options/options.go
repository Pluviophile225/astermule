package options

import "github.com/spf13/cobra"

type Options struct {
	Address    string
	Port       uint
	Target     string
	DagStr     string
	EntryParam string
}

func NewOptions() *Options {
	return &Options{}
}

// initFlags initializes flags by section name.
func (o *Options) Parse(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&o.Address, "address", "0.0.0.0", "The boot address of launching astermule.")
	cmd.PersistentFlags().UintVar(&o.Port, "port", 8080, "The boot port of launching astermule.")
	cmd.PersistentFlags().StringVar(&o.Target, "target", "/", "The target of launching astermule.")
	cmd.PersistentFlags().StringVar(&o.DagStr, "dag", "{}", "The dag of launching astermule.")
	cmd.PersistentFlags().StringVar(&o.EntryParam, "param", "", "The EntryParam of launching astermule.")
}
