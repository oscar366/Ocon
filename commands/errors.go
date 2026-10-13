package commands

import (
	"oscarsgoofysite/OCON/helpers/executeExecute"
	"oscarsgoofysite/OCON/helpers"
	"oscarsgoofysite/OCON/state"
)
//This errors.go is for dealing with errors (there are 2 cuz i am dumb)
func Try(args []string) error {
	/*
	args to int table
	var "a COMMAND COMMANDARGS
	0   1    2       3...
	
	COMMAND COMMANDARGS
	0          1...
	*/
	isVar := false
	if args[0] == "var" {
		isVar = true
	}
	var err error
	if isVar {
		err = executeExecute.ExecuteCommand(args[2:])
	} else {
		err = executeExecute.ExecuteCommand(args[0:])
	}
	
	if err != nil {
		// Try caught it.
		// If we need to set a var do so
		if isVar {
			if args[0][0] == '"' {
				state.VarStorage[args[1][1:]] = "\"" + err.Error()
			} else {
				return helpers.NewError("Try did not have the propper type for var")
			}
		}
	}

	return nil
}