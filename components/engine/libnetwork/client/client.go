package client

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	flag "github.com/docker/docker/pkg/mflag"
)

// CallFunc provides environment specific call utility to invoke backend functions from UI
type CallFunc func(string, string, interface{}, map[string][]string) (io.ReadCloser, int, error)

// NetworkCli is the UI object for network subcmds
type NetworkCli struct {
	out  io.Writer
	err  io.Writer
	call CallFunc
}

// NewNetworkCli is a conveninent function to create a NetworkCli object
func NewNetworkCli(out, err io.Writer, call CallFunc) *NetworkCli {
	return &NetworkCli{
		out:  out,
		err:  err,
		call: call,
	}
}

// getMethod is Borrowed from Docker UI which uses reflection to identify the UI Handler
func (cli *NetworkCli) getMethod(args ...string) (func(string, ...string) error, bool) {
	camelArgs := make([]string, len(args))
	for i, s := range args {
		if len(s) == 0 {
			return nil, false
		}
		camelArgs[i] = strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}
	methodName := "Cmd" + strings.Join(camelArgs, "")
	method := reflect.ValueOf(cli).MethodByName(methodName)
	if !method.IsValid() {
		return nil, false
	}
	return method.Interface().(func(string, ...string) error), true
}

// Cmd is borrowed from Docker UI and acts as the entry point for network UI commands.
// network UI commands are designed to be invoked from multiple parent chains
func (cli *NetworkCli) Cmd(chain string, args ...string) error {
	var errStr string
	if len(args) > 1 {
		method, exists := cli.getMethod(args[:2]...)
		if exists {
			return method(chain, args[2:]...)
		}
	}
	if len(args) > 0 {
		method, exists := cli.getMethod(args[0])
		if !exists {
			errStr = fmt.Sprintf("%s: '%s' is not a %s command. See '%s --help'.\n", chain, args[0], chain, chain)
			fmt.Fprintf(cli.err, errStr)
			return errors.New(errStr)
		}
		return method(chain, args[1:]...)
	}
	errStr = fmt.Sprintf("'%s' is not a valid command. See '%s --help'.\n", chain, chain)
	fmt.Fprintf(cli.err, errStr)
	return errors.New(errStr)
}

// Subcmd is borrowed from Docker UI and performs the same function of configuring the subCmds
func (cli *NetworkCli) Subcmd(chain, name, signature, description string, exitOnError bool) *flag.FlagSet {
	var errorHandling flag.ErrorHandling
	if exitOnError {
		errorHandling = flag.ExitOnError
	} else {
		errorHandling = flag.ContinueOnError
	}
	flags := flag.NewFlagSet(name, errorHandling)
	flags.Usage = func() {
		options := ""
		if signature != "" {
			signature = " " + signature
		}
		if flags.FlagCountUndeprecated() > 0 {
			options = " [OPTIONS]"
		}
		fmt.Fprintf(cli.out, "\nUsage: %s %s%s%s\n\n%s\n\n", chain, name, options, signature, description)
		flags.SetOutput(cli.out)
		flags.PrintDefaults()
		os.Exit(0)
	}
	return flags
}

func readBody(stream io.ReadCloser, statusCode int, err error) ([]byte, int, error) {
	if stream != nil {
		defer stream.Close()
	}
	if err != nil {
		return nil, statusCode, err
	}
	body, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, -1, err
	}
	return body, statusCode, nil
}
