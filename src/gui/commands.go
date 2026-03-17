package main

import (
	"errors"
	keeper "github.com/Supraboy981322/keeper/golang"
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
		
		//show/hide the scrollback buffer
	  case "hide": { state.Scrollback.Hide = true }
	  case "show": { state.Scrollback.Hide = false }

		//toggle something
	  case "toggle": { return toggle(args) }

		//change a config setting
	  case "set": { return setting(args) }

		// TODO: more commands

	  default: {
			return errors.New("unknown command"), Event(ERR)
		}
	}

	//most commands are a NOP (at the moment)
	return nil, Event(NOP)
}

//helper to parse args
func parse_args(in string) []string {
	//state
	var res []string
	var esc, stringing bool 
	var skipping rune
	var mem []rune
	var cmd_start int

	drain := func() {
		keeper.Add(&res, string(keeper.Drain(&mem)))
	}
	for ; in[cmd_start] == ':'; cmd_start++ {}
	loop: for _, r := range in[cmd_start:] {
		if esc { keeper.Add(&mem, r) ; continue loop }

		if skipping != 0 && r != skipping {
			drain()
			skipping = 0
		} else if skipping != 0 {
			continue loop
		}

		sw: switch r {
			case '"': {
				if stringing { drain() }
				keeper.Flip(&stringing)
			}
		  case '\\': { keeper.Flip(&esc) }
			case ' ', '\t', '\n': if !stringing {
				skipping = r
				break sw
			}; fallthrough
		  default: { keeper.Add(&mem, r) }
		}
	}
	keeper.Add(&res, string(mem))
	keeper.Filter(&res, "")
	return res
}

// TODO: this
func toggle(args []string) (error, Event){
	if len(args) < 2 {
		return errors.New("not enough args. need something to toggle"), Event(ERR)
	}
	switch args[1] {
    case "scrollback", "history", "hist", "hide": { keeper.Flip(&state.Scrollback.Hide) }
	  default: {
			return errors.New("don't know how to toggle "+args[1]), Event(ERR)
		}
	}
	return nil, Event(NOP)
}


func setting(args []string) (error, Event) {
	if len(args) < 2 {
		return errors.New("not enough args, need something to change"), Event(ERR)
	}
	switch args[1] {
	  case "server": { /* TODO: this */ }
	}
	// TODO: this
	return nil,  Event(SETTING)
}
