package main

import (
	"errors"
)

func cmd() (error, Event) {
	//reset the buffer on return 
	defer func() {
		state.CmdBuf = []rune{':'}
	}()

	//parse the args
	args := parse_args(string(state.CmdBuf))
	//NOP on empty input
	if len(args) < 1 { return nil , Event(NOP) }
	//switch on command name
	switch args[0] {
		//quit (ignores any remaining args)
		case "q": { state.Exit = true }
	  case "hide": { state.Scrollback.Hide = true }
	  case "show": { state.Scrollback.Hide = false }

	  case "toggle": { return toggle(args) }

		// TODO: more commands

	  default: {
			return errors.New("unknown command"), Event(ERR)
		}
	}

	//most commands are a NOP (at the moment)
	return nil, Event(NOP)
}

func parse_args(in string) []string {
	var res []string
	var esc, stringing bool 
	var skipping rune
	var mem []rune
	var cmd_start int
	for ; in[cmd_start] == ':'; cmd_start++ {}
	loop: for _, r := range in[cmd_start:] {
		if esc { add(&mem, r) ; continue loop }

		if skipping != 0 && r != skipping {
			add(&res, string(drain(&mem)))
			skipping = 0
		} else if skipping != 0 {
			continue loop
		}

		//basic common chars that need to be handled
		//  NOTE: appostrophee currently isn't supported for string
		//  TODO: ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
		switch r {
			case '"': { flip(&stringing) }
			case ' ', '\t', '\n': { skipping = r }
		  case '\\': { flip(&esc) }
		  default: { add(&mem, r) }
		}
	}
	add(&res, string(mem))
	//l := len(res)
	//for _, a := range res {
	//	if len(a) != 0 { add(&res, a) }
	//}
	filter(&res, "")
	return res
}

// TODO: this
func toggle(args []string) (error, Event){
	if len(args) < 2 {
		return errors.New("not enough args. need something to toggle"), Event(ERR)
	}
	switch args[1] {
    case "scrollback", "history", "hist", "hide": { flip(&state.Scrollback.Hide) }
	  default: {
			return errors.New("don't know how to toggle "+args[1]), Event(ERR)
		}
	}
	return nil, Event(NOP)
}
