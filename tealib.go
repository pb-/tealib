package tealib

import (
	"os"
	"strings"
	"time"
)

// Message describes an event that happened.
type Message interface{}

// State holds the entire program's state.
type State interface{}

// Command describes a side effect to be executed by the runtime.
type Command interface{}

// None means no commands.
var None = []Command{}

// InitFunc returns the initial state and initial commands.
type InitFunc func(args []string, env map[string]string) (State, []Command)

// UpdateFunc receives the current state, a message, and the current time and returns the new state and commands to run.
type UpdateFunc func(State, Message, time.Time) (State, []Command)

// RenderFunc renders state at given time.
type RenderFunc func(State, time.Time)

// NoRender is an empty render function.
func NoRender(State, time.Time) {}

// ModuleFunc implements the entry point of a module.
type ModuleFunc func(chan<- Message) chan<- Command

// OptionFunc implements a library run-time option.
type OptionFunc func(*tealib)

// Module adds the supplied module to the runtime.
func Module(module ModuleFunc) OptionFunc {
	return func(t *tealib) {
		t.modules = append(t.modules, module)
	}
}

type tealib struct {
	modules []ModuleFunc
}

// Run starts the runtime using the provided functions.
func Run(init InitFunc, update UpdateFunc, render RenderFunc, options ...OptionFunc) {
	t := &tealib{
		modules: []ModuleFunc{},
	}

	for _, option := range options {
		option(t)
	}

	messages := make(chan Message)
	commands := make(chan Command)
	sinks := []chan<- Command{}

	for _, module := range t.modules {
		sinks = append(sinks, module(messages))
	}

	go commandBroadcaster(commands, sinks)
	loop(init, update, render, messages, commands)
}

func loop(init InitFunc, update UpdateFunc, render RenderFunc, messages <-chan Message, commands chan<- Command) {
	state, cmds := init(os.Args, env())
	for _, cmd := range cmds {
		commands <- cmd
	}

	render(state, time.Now())

	for {
		msg := <-messages
		state, cmds = update(state, msg, time.Now())
		for _, cmd := range cmds {
			commands <- cmd
		}

		render(state, time.Now())
	}
}

func commandBroadcaster(commands <-chan Command, sinks []chan<- Command) {
	for {
		cmd := <-commands
		for _, sink := range sinks {
			sink <- cmd
		}
	}
}

func env() map[string]string {
	vars := map[string]string{}

	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		vars[pair[0]] = pair[1]
	}

	return vars
}
