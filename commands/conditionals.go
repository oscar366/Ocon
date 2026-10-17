package commands

import (
	"oscarsgoofysite/OCON/state"
	"oscarsgoofysite/OCON/helpers"
	"oscarsgoofysite/OCON/helpers/executeExecute"
	//"strconv"
	//"strings"
	//"time"
	"os"
	"slices"
)



//However all ideas and methouds of doing are my ideas

func Endfor(args []string) error {
	//we not having the condishon in for
	if args[0] == "~true" {
		//go back to for loop
		
		//get the lines
		
		n := 1
		for state.Pointer-n >= 0 && state.DocumentData[state.Pointer-n] != "for" {
			n++
		}
		
		if state.Pointer-n < 0 {
			// "for" wasn't found
			return helpers.NewError("Error: no acompnying for statement for currently executeing endfor")
		} else {
			// found it at state.Pointer - n
			to := state.Pointer - n
			state.Pointer = to - 1 //-1 cuz after this it increse sate.Pointer by 1
		}
	} else if args[0] == "~false" {
		//we dont do anythin
		return nil
	} else {
		//Error
		return helpers.NewError("Non bool put in Endfor?")
	}
	return nil
}
func If(args []string) error {
	//check if it has else and loc
	elsePos := slices.Index(args, "else")//find pos of else
	
	
	if args[0] == "~true" {
		if elsePos == -1 {
			return ifCommands(args[1:])
		} else {
			return ifCommands(args[1:elsePos])
		}
	} else if args[0] == "~false" {
		if elsePos == -1 {//if no else
			return nil
		} else {
			/* elsePos+1
			   if cmd args else cmd args
								 ^
								 |
								 
				[(elsePos+1):]
				
				if cmd args else cmd args
								  [    ]
			*/
			return ifCommands(args[(elsePos+1):])
		}
	} else {
		return helpers.NewError("Got nether '~true' nor '~false' in if bool")
	}
}
func ifCommands(args []string) error {
	if len(args) == 0 {
		return helpers.NewError("if statement has no command")
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {

	case "execute":
		return executeExecute.ExecuteCommand(cmdArgs)

	case "goto":
		return GotoCmd(cmdArgs)

	case "continue":
		if len(cmdArgs) != 0 {
			return helpers.NewError("In if statement continue args was not = 0")
		}
		return nil

	case "setvar":
		return VarCmd(cmdArgs)

	case "throwError":
		if len(cmdArgs) != 1 {
			return helpers.NewError("In if statement throwError args was not = 1")
		}
		return helpers.NewError(cmdArgs[0][1:])

	case "echo":
		return EchoCmd(cmdArgs)

	case "crash":
		if len(cmdArgs) != 0 {
			return helpers.NewError("In if statement crash args was not = 0")
		}
		os.Exit(0)
		return nil

	default:
		return helpers.NewError(
			"Not a command in IF [" + cmd + "] maybe you should ask oscar to add it?",
		)
	}
}