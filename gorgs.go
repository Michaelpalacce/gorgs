package gorgs

import (
	"flag"
	"fmt"
	"strings"
)

type Opt struct {
	VarPtr        any
	DefaultValue  any
	Description   string
	ShorthandFlag string
	LonghandFlag  string
}

type Printer func(format string, a ...any) (n int, err error)

type Gorgs struct {
	usage     string
	examples  string
	fs        *flag.FlagSet
	arguments []string
	printer   Printer

	opts []Opt
}

// NewGorgs is used to create a new instance of Gorgs
// The arguments are generally os.Args[1:]
// options can be used to modify the Gorgs behavior. If they are not passed during creation, you can do so later by calling g.Modify()
// Default printer will be set to fmt.Printf
func NewGorgs(arguments []string, options ...GorgsOptions) (*Gorgs, error) {
	args := &Gorgs{
		arguments: arguments,
		printer:   fmt.Printf,
	}

	if err := args.Modify(options...); err != nil {
		return nil, err
	}

	return args, nil
}

// AddOpt can be used instead of AddVar if you want to build the Opt yourself.
func (a *Gorgs) AddOpt(opt Opt) error {
	return a.AddVar(opt.VarPtr, opt.LonghandFlag, opt.ShorthandFlag, opt.DefaultValue, opt.Description)
}

// AddVar adds variables to be parsed later when calling Gorgs.Parse
func (a *Gorgs) AddVar(varPtr any, longhand string, shorthand string, defaultValue any, description string) error {
	fs := a.getCorrectFs()

	a.opts = append(a.opts, Opt{
		VarPtr:        varPtr,
		DefaultValue:  defaultValue,
		ShorthandFlag: shorthand,
		LonghandFlag:  longhand,
		Description:   description,
	})

	switch v := varPtr.(type) {
	case *string:
		if dv, ok := defaultValue.(string); ok {
			if shorthand != "" {
				fs.StringVar(v, shorthand, dv, description)
			}
			if longhand != "" {
				fs.StringVar(v, longhand, dv, description)
			}
		} else {
			return fmt.Errorf("%v should have been a string", defaultValue)
		}
	case *bool:
		if dv, ok := defaultValue.(bool); ok {
			if shorthand != "" {
				fs.BoolVar(v, shorthand, dv, description)
			}
			if longhand != "" {
				fs.BoolVar(v, longhand, dv, description)
			}
		} else {
			return fmt.Errorf("%v should have been a bool", defaultValue)
		}
	case *int:
		if dv, ok := defaultValue.(int); ok {
			if shorthand != "" {
				fs.IntVar(v, shorthand, dv, description)
			}
			if longhand != "" {
				fs.IntVar(v, longhand, dv, description)
			}
		} else {
			return fmt.Errorf("defaultValue %v should have been an int", defaultValue)
		}
	default:
		return fmt.Errorf("var must be a pointer of string, bool or int")
	}

	return nil
}

// GetUsage fetches the usage information based on the information given to all the variables that were added
func (a *Gorgs) GetUsage() {
	maxShort := 0
	maxLong := 0

	for _, o := range a.opts {
		if o.ShorthandFlag != "" && len(o.ShorthandFlag) > maxShort {
			maxShort = len(o.ShorthandFlag)
		}
		if o.LonghandFlag != "" && len(o.LonghandFlag) > maxLong {
			if len(o.LonghandFlag) > maxLong {
				maxLong = len(o.LonghandFlag)
			}
		}
	}

	_, _ = a.printer("%s\n", a.usage)

	if len(a.opts) > 0 {
		_, _ = a.printer("Options:\n")
	}

	for _, o := range a.opts {
		short := ""
		if o.ShorthandFlag != "" {
			short = fmt.Sprintf("-%-*s", maxShort, o.ShorthandFlag)
		} else {
			short = strings.Repeat(" ", maxShort+1)
		}

		long := ""
		if o.LonghandFlag != "" {
			long = fmt.Sprintf("--%-*s", maxLong, o.LonghandFlag)
		} else {
			long = strings.Repeat(" ", maxLong+2)
		}

		desc := o.Description

		if o.DefaultValue != nil && o.DefaultValue != "" {
			desc += fmt.Sprintf(" (default: %v)", o.DefaultValue)
		}

		_, _ = a.printer("    %s    %s    %s\n", short, long, desc)
	}
	if a.examples != "" {
		_, _ = a.printer("%s\n", a.examples)
	}
}

// Parse will parse the given cmdline arguments
func (a *Gorgs) Parse() error {
	fs := a.getCorrectFs()

	return fs.Parse(a.arguments)
}

func (a *Gorgs) getCorrectFs() *flag.FlagSet {
	if a.fs != nil {
		return a.fs
	} else {
		return flag.CommandLine
	}
}
