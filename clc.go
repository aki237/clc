package clc

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

// Commands is a dictionary of all commands to their correspnding Command structs.
type Commands map[string]Command

// App structure containing necessary fields for running a command line app.
type App struct {
	name     string
	info     string
	version  string
	commands Commands
}

// NewApp returns an App pointer with the passed in values as initializers.
func NewApp(name, info, version string) *App {
	if name == "" {
		name = os.Args[0]
	}
	if version == "" {
		version = "v0.00"
	}
	if info == "" {
		info = "Command-line tool"
	}
	return &App{name: name, info: info, version: version, commands: make(Commands, 0)}
}

// AddCommand method is used to add a new command and it's corresponding function.
func (a *App) AddCommand(command, help string, cmd Program) error {
	if reflect.ValueOf(cmd).Kind() != reflect.Ptr && reflect.ValueOf(cmd).IsValid() {
		return errors.New("Expected a pointer recieved a value")
	}
	a.commands[command] = Command{Args: cmd, help: help}
	return nil
}

// Usage prints the usage information of the application
func (a *App) Usage() {
	fmt.Fprint(os.Stderr, "NAME:\n\t")
	fmt.Fprint(os.Stderr, a.name+" - "+a.info+"\n\n")
	fmt.Fprint(os.Stderr, "USAGE:\n\t")
	fmt.Fprint(os.Stderr, a.name+" [global options] command [command options] [arguments...]\n\n")
	fmt.Fprint(os.Stderr, "VERSION:\n\t")
	fmt.Fprint(os.Stderr, a.version+"\n\n")
	fmt.Fprint(os.Stderr, "COMMANDS:")
	for command, sct := range a.commands {
		fmt.Fprint(os.Stderr, "\n\t"+command+"\t"+sct.help)
	}
	fmt.Fprint(os.Stderr, "\n\n")

}

// Run method is used to parse the commandline arguments and run the app.
func (a *App) Run() error {
	if len(os.Args) == 1 {
		a.Usage()
		return nil
	}
	command := os.Args[1]
	if command == "-h" || command == "-help" {
		a.Usage()
		return nil
	}
	c, ok := a.commands[command]
	if !ok {
		return errors.New(a.name + ": '" + command + "' is not a " + a.name + " command.")
	}
	if len(os.Args) < 2 {
		if c.Args != nil {
			return errors.New("Not enough arguments to the '" + command + "' command.")
		}
	}
	if c.Args == nil {
		return c.Exec()
	}
	if err := c.ParseArgs(os.Args[2:]); err != nil {
		return err
	}
	return c.Exec()
}
