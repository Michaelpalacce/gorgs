package gorgs

import (
	"flag"
	"fmt"
	"strings"
)

type opt struct {
	defaultValue  any
	description   string
	shorthandFlag string
	longhandFlag  string
}

type Gorgs struct {
	usage     string
	examples  string
	fs        *flag.FlagSet
	arguments []string

	opts []opt
}

type GorgsOptions func(*Gorgs) error

// NewGorgs is used to create a new instance of Gorgs
// The arguments are generally os.Args[1:]
// options can be used to modify the Gorgs behavior. If they are not passed during creation, you can do so later by calling g.Modify()
func NewGorgs(arguments []string, options ...GorgsOptions) (*Gorgs, error) {
	args := &Gorgs{
		arguments: arguments,
	}

	if err := args.Modify(options...); err != nil {
		return nil, err
	}

	return args, nil
}

// AddVar adds variables to be parsed later when calling Gorgs.Parse
func (a *Gorgs) AddVar(varPtr any, longhand string, shorthand string, defaultValue any, message string) {
	fs := a.getCorrectFs()

	a.opts = append(a.opts, opt{
		defaultValue:  defaultValue,
		shorthandFlag: shorthand,
		longhandFlag:  longhand,
		description:   message,
	})

	switch v := varPtr.(type) {
	case *string:
		if dv, ok := defaultValue.(string); ok {
			if shorthand != "" {
				fs.StringVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				fs.StringVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been a string", defaultValue))
		}
	case *bool:
		if dv, ok := defaultValue.(bool); ok {
			if shorthand != "" {
				fs.BoolVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				fs.BoolVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been a bool", defaultValue))
		}
	case *int:
		if dv, ok := defaultValue.(int); ok {
			if shorthand != "" {
				fs.IntVar(v, shorthand, dv, message)
			}
			if longhand != "" {
				fs.IntVar(v, longhand, dv, message)
			}
		} else {
			panic(fmt.Sprintf("%v should have been an int", defaultValue))
		}
	default:
		panic("Var must be a pointer of string, bool or int")
	}
}

// GetUsage fetches the usage information based on the information given to all the variables that were added
func (a *Gorgs) GetUsage() {
	maxShort := 0
	maxLong := 0

	for _, o := range a.opts {
		if o.shorthandFlag != "" && len(o.shorthandFlag) > maxShort {
			maxShort = len(o.shorthandFlag)
		}
		if o.longhandFlag != "" && len(o.longhandFlag) > maxLong {
			if len(o.longhandFlag) > maxLong {
				maxLong = len(o.longhandFlag)
			}
		}
	}

	fmt.Println(a.usage)

	if len(a.opts) > 0 {
		fmt.Println("Options:")
	}

	for _, o := range a.opts {
		short := ""
		if o.shorthandFlag != "" {
			short = fmt.Sprintf("-%-*s", maxShort, o.shorthandFlag)
		} else {
			short = strings.Repeat(" ", maxShort+1)
		}

		long := ""
		if o.longhandFlag != "" {
			long = fmt.Sprintf("--%-*s", maxLong, o.longhandFlag)
		} else {
			long = strings.Repeat(" ", maxLong+2)
		}

		desc := o.description

		if o.defaultValue != nil && o.defaultValue != "" {
			desc += fmt.Sprintf(" (default: %v)", o.defaultValue)
		}

		fmt.Printf("    %s    %s    %s\n", short, long, desc)
	}
	if a.examples != "" {
		fmt.Println(a.examples)
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
