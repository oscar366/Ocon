package commands

import (
	//"fmt"
	"oscarsgoofysite/OCON/helpers"
	"oscarsgoofysite/OCON/state"
	"strconv"
)

//this contains all var bars
func VarCmd(args []string) error {

	if len(args) == 0 { //you know what this means
		return helpers.NewError("No args for var command:0")
	}

	//empty var
	if len(args) == 1 {
		state.VarStorage[args[0][1:]] = ""
		return nil
	}

	if len(args) == 2 {
		//state.VarStorage[0]
		//state.VarStorage[1]
		//append(
		//len

		//remember to do :1 otherwise to init the var its $"erer not $erer
		state.VarStorage[args[0][1:]] = args[1]

		return nil
	}

	return helpers.NewError("too many args for var command")
}

func UnLinkVar(args []string) error {
	
	if len(args) == 0 { //you know what this means
		return helpers.NewError("No args for unlink command")
	} else if len(args) == 2 {
		return helpers.NewError("To meny args for unlink command")
	}
	
	delete(state.VarStorage, args[0])
	return nil
}


func IncrementCmd(args []string) error {
	if len(args) < 2 {
		return helpers.NewError("not enogh args for increment")
	}

	if args[1][0] != '\'' {
		return helpers.NewError("incorecct type for incrment")
	}

	varname := args[0][1:]
	val, ok := state.VarStorage[varname]
	if !ok {
		return helpers.NewError("varname: " + varname + " was not found.")
	}

	val = val[1:]

	//debug line
	//fmt.Printf("val = %q\n", val)

	//check if they can be turned into numbers
	current, err := strconv.Atoi(val)
	if err != nil {
		return helpers.NewError("variable is not a number")
	}

	inc, err := strconv.Atoi(args[1][1:])
	if err != nil {
		return helpers.NewError("increment is not a number")
	}

	current = current + inc // converts the string into a number were 1: (as strings are arrys of chars) uis everytihng past the "   and then adds it to the vars value
	state.VarStorage[varname] = "\"" + strconv.Itoa(current) // return to string
	return nil
} //also rember to add the ' back in otherwise it all crashes



func DecremntCmd(args []string) error {
	if len(args) < 2 {
		return helpers.NewError("not enogh args for Decrement")
	}

	if args[1][0] != '\'' {
		return helpers.NewError("incorecct type for Decrment")
	}

	varname := args[0][1:]
	val, ok := state.VarStorage[varname]
	if !ok {
		return helpers.NewError("varname: " + varname + " was not found.")
	}

	val = val[1:]

	//debug line
	//fmt.Printf("val = %q\n", val)

	//check if they can be turned into numbers
	current, err := strconv.Atoi(val)
	if err != nil {
		return helpers.NewError("variable is not a number")
	}

	inc, err := strconv.Atoi(args[1][1:])
	if err != nil {
		return helpers.NewError("decrement is not a number")
	}

	current = current - inc // converts the string into a number were 1: (as strings are arrys of chars) uis everytihng past the "   and then adds it to the vars value
	state.VarStorage[varname] = "\"" + strconv.Itoa(current) // return to string
	return nil
}