package cli

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

type HandleFunc func(set *flag.FlagSet) bool

// FlagHandle 命令行处理
func FlagHandle() bool {
	if len(os.Args) <= 1 {
		return true
	}
	logrus.Debugln(os.Args, os.Args[1])
	f := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	Allocation(f, os.Args[1])
	return false
}

var (
	reg = make(map[string]HandleFunc)
)

func Allocation(set *flag.FlagSet, command string) bool {
	handle, ok := reg[command]
	if !ok {
		set.ErrorHandling()
		return false
	}
	return handle(set)
}

// Register 注册命令行
func Register(name string, fn HandleFunc) {
	reg[name] = fn
}
