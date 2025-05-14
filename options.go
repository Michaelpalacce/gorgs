package gorgs

import "flag"

// WithUsage can be used to add extra usage information to the current command
func WithUsage(usage string) GorgsOptions {
	return func(a *Gorgs) error {
		a.usage = usage
		return nil
	}
}

// WithExamples Will add examples to the current comamnd
func WithExamples(examples string) GorgsOptions {
	return func(a *Gorgs) error {
		a.examples = examples
		return nil
	}
}

// WithFs can be used to pass your own FlagSet to mainly define your own ErrorHandling.
// By default if you don't pass a Flagset then the default flag.CommandLine will be used
func WithFs(fs *flag.FlagSet) GorgsOptions {
	return func(a *Gorgs) error {
		a.fs = fs
		fs.Usage = a.GetUsage
		return nil
	}
}

// Modify is used to Modify the Gorgs with different stuff like examples or usage
func (s *Gorgs) Modify(options ...GorgsOptions) error {
	for _, option := range options {
		if err := option(s); err != nil {
			return err
		}
	}

	return nil
}
