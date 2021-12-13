package app

import (
	"fmt"
	"os"

	"github.com/fairyhunter13/envcompact/internal/customlog"
	"github.com/fairyhunter13/envcompact/pkg/compacter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Option struct {
	Verbose      bool
	Silent       bool
	PrintVersion bool
	Input        string
}

func (o *Option) Assign(opts []FuncOption) *Option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type Instance struct {
	option     *Option
	fileSource *os.File
}

type FuncOption func(*Option)

func WithVerbosity(verbose, silent bool) FuncOption {
	return func(o *Option) {
		o.Verbose = verbose
		o.Silent = silent
	}
}

func WithInputPath(input string) FuncOption {
	return func(o *Option) {
		o.Input = input
	}
}

func New(opts ...FuncOption) *Instance {
	opt := new(Option).Assign(opts)

	return &Instance{
		option: opt,
	}
}

func (app *Instance) Init() (err error) {
	if app.option.Input == "" {
		customlog.Get().Info("Reading from os.Stdin.")
		app.fileSource = os.Stdin
	} else {
		customlog.Get().Info(
			"Reading from input file.",
			zap.String("path", app.option.Input),
		)
		app.fileSource, err = os.Open(app.option.Input)
		err = errors.WithStack(err)
	}
	return
}

func (app *Instance) Close() {
	if app.fileSource != nil {
		app.fileSource.Close()
	}
}

func (app *Instance) Run() (err error) {
	var (
		compactEngine = compacter.NewParser(app.fileSource)
		line          string
		lineNum       = uint64(1)
		first         = true
	)
	for compactEngine.Scan() {
		if first {
			first = false
		} else {
			fmt.Println()
		}
		line = compactEngine.Text()
		customlog.Get().Debug(
			"Reading from the file source.",
			zap.String("line", line),
			zap.Uint64("lineNum", lineNum),
		)
		fmt.Print(line)
		lineNum++
	}
	return
}
