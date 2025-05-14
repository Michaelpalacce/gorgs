package gorgs_test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/Michaelpalacce/gorgs"
)

var printerStorage []string = []string{}

func testPrinter(format string, a ...any) (n int, err error) {
	printerStorage = append(printerStorage, fmt.Sprintf(format, a...))
	return 0, nil
}

func TestGorgs_Parse(t *testing.T) {
	var printer gorgs.Printer = testPrinter

	tests := []struct {
		name      string // description of this test case
		arguments []string
		options   []gorgs.GorgsOptions
		opts      []gorgs.Opt
		// expected VarPtr values in order
		expected []any
		wantErr  bool
	}{
		{
			name:      "AddVar with shorthand and default, uses shorthand",
			arguments: []string{"-a=test"},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *string {
						a := ""
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  "default",
					Description:   "Some Description",
				},
			},
			expected: []any{"test"},
			wantErr:  false,
		},
		{
			name:      "AddVar with longhand and default, uses longhand",
			arguments: []string{"--append=test"},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *string {
						a := ""
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  "default",
					Description:   "Some Description",
				},
			},
			expected: []any{"test"},
			wantErr:  false,
		},
		{
			name:      "AddVar with default, uses default",
			arguments: []string{},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *string {
						a := ""
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  "default",
					Description:   "Some Description",
				},
			},
			expected: []any{"default"},
			wantErr:  false,
		},
		{
			name:      "AddVar with longhand and shorthand and default uses last",
			arguments: []string{"-a=short", "--append=long"},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *string {
						a := ""
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  "default",
					Description:   "Some Description",
				},
			},
			expected: []any{"long"},
			wantErr:  false,
		},
		{
			name:      "AddVar with int",
			arguments: []string{"-a=1"},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *int {
						a := 0
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  2,
					Description:   "Some Description",
				},
			},
			expected: []any{1},
			wantErr:  false,
		},
		{
			name:      "AddVar with int default",
			arguments: []string{},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *int {
						a := 0
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  2,
					Description:   "Some Description",
				},
			},
			expected: []any{2},
			wantErr:  false,
		},
		{
			name:      "AddVar with bool",
			arguments: []string{"-a"},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *bool {
						a := false
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  false,
					Description:   "Some Description",
				},
			},
			expected: []any{true},
			wantErr:  false,
		},
		{
			name:      "AddVar with bool default",
			arguments: []string{},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts: []gorgs.Opt{
				{
					VarPtr: func() *bool {
						a := false
						return &a
					}(),
					LonghandFlag:  "append",
					ShorthandFlag: "a",
					DefaultValue:  false,
					Description:   "Some Description",
				},
			},
			expected: []any{false},
			wantErr:  false,
		},
		{
			name:      "AddVar fails when arg provided but no opt",
			arguments: []string{"-a=short"},
			options:   []gorgs.GorgsOptions{gorgs.WithPrinter(printer), gorgs.WithFs(flag.NewFlagSet("", flag.ContinueOnError))},
			opts:      []gorgs.Opt{},
			expected:  []any{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := gorgs.NewGorgs(tt.arguments, tt.options...)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			for _, opt := range tt.opts {
				gotErr := a.AddOpt(opt)
				if gotErr != nil {
					if !tt.wantErr {
						t.Errorf("AddOpt() failed: %v", gotErr)
					}
					return
				}
			}

			err = a.Parse()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Parse() succeeded unexpectedly. err was %s", err.Error())
				}
				return
			}

			for i, opt := range tt.opts {
				if len(tt.expected) <= i {
					continue
				}

				switch v := opt.VarPtr.(type) {
				case *string:
					if ev, ok := tt.expected[i].(string); ok {
						if ev != *v {
							t.Errorf("%s, does not equal %s", ev, *v)
						}
					} else {
						t.Errorf("Variable at position %d expected to be a string but was not", i)
					}
				case *int:
					if ev, ok := tt.expected[i].(int); ok {
						if ev != *v {
							t.Errorf("%d, does not equal %d", ev, *v)
						}
					} else {
						t.Errorf("Variable at position %d expected to be a int but was not", i)
					}

				case *bool:
					if ev, ok := tt.expected[i].(bool); ok {
						if ev != *v {
							t.Errorf("%v, does not equal %v", ev, *v)
						}
					} else {
						t.Errorf("Variable at position %d expected to be a bool but was not", i)
					}
				default:
					t.Fatalf("Could not determine type of opt: %v", opt)
					return
				}
			}
		})
	}
}
