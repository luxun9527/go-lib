package inject

import (
	"flag"
	"os"
)

//面对的情况，当设置了某个flag，就执行某个函数。
var	defaultFlag = &Flag{
	FlagSet: flag.CommandLine,
	m:       make(map[string]func(val flag.Value), 5),
}
type Flag struct {
	*flag.FlagSet
	m map[string]func(val flag.Value)
}

func Register(name string, handler func(val flag.Value)) {
	defaultFlag.register(name, handler)
}
func Parse() {
	defaultFlag.parseFlag()
}

func (f *Flag) register(name string, handler func(val flag.Value)) {
	f.m[name] = handler
}
func (f *Flag) parseFlag() {
	f.Parse(os.Args[1:])
	f.Visit(func(fl *flag.Flag) {
		f1, ok := f.m[fl.Name]
		if !ok {
			return
		}
		f1(fl.Value)
	})
}
