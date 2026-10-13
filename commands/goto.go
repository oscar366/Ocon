package commands

import (
    "oscarsgoofysite/OCON/state"
	"oscarsgoofysite/OCON/helpers"
	"strconv"
)

func GotoCmd(args []string) error {
    //fmt.Println("going to the line at: " + args[0])

    if len(args) == 0 {//if the goto has NO args eg (goto)
		return helpers.NewError("no Args for goto")
    }

    wereto := args[0]
    if wereto[0] == '\'' {//gotos a line number
		n, err := strconv.Atoi(wereto[1:]) // converts the string into a number were 1: (as strings are arrys of chars) uis everytihng past the "  
		if err != nil {//if it aculy is NOT a number then leave
			return helpers.NewError("invalid number:" + wereto[1:])
		}
		state.Pointer = n - 1 //-1 becuse after this the poiter adds 1
	} else if wereto[0] == '"' { //goses to a sec, section, or, § (all the same thing but i digress)
		userSec := wereto[1:]//this is everything past the "
		list := state.SectionList
		sec, ok := list[userSec]//this is an int (map[string]int) of were the pointer needs to go
		if !ok {
			return helpers.NewError("section " + userSec + " not found. WARNING ocon is contuning to next line")
		}
		state.Pointer = sec - 1//-1 becuse after this the poiter adds 1
	}
	return nil//end of program
}