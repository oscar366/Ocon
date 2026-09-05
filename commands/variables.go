package commands

import (
	"fmt"
	"oscarsgoofysite/OCON/state"
	"strconv"
)


//this contains all var bars
func VarCmd(args []string) {
	
	if len(args) == 0 {//you know what this means
		fmt.Println("No args for var command:0")
	}
	
	//empty var
	if len(args) == 1 {
		state.VarStorage[args[0][1:]] = ""
		//fmt.Println("epmty vars have not yet been implemnted")
	}
	
	if len(args) == 2 {
		//state.VarStorage[0]
		//state.VarStorage[1]
		//append(
		//len
		
		
		//remember to do :1 otherwise to init the var its $"erer not $erer
		state.VarStorage[args[0][1:]] = args[1]
		
		return
	}
	
}

func IncrementCmd(args []string) {
	if len(args) < 2 {
		fmt.Println("not enogh args for increment")
		return
	}
	
	if args[1][0] != '\'' {
		fmt.Println("incorecct type for incrment")
	}
	
	varname := args[0][1:]
	val, ok := state.VarStorage[varname]
	if !ok {
		fmt.Println("varname: " + varname + " was not found.")
		return
	}
	val = val[1:]
	
	//debug line
	//fmt.Printf("val = %q\n", val)
	
	//check if they can be turned into numbers
	current, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("variable is not a number")
		return
	}

	inc, err := strconv.Atoi(args[1][1:])
	if err != nil {
		fmt.Println("increment is not a number")
		return
	}
	
	current = current + inc // converts the string into a number were 1: (as strings are arrys of chars) uis everytihng past the "   and then adds it to the vars value
	state.VarStorage[varname] = "\"" + strconv.Itoa(current) // return to string
}	//also rember to add the ' back in otherwise it all crashes
func DecremntCmd(args []string) {
	if len(args) < 2 {
		fmt.Println("not enogh args for Decrement")
		return
	}
	
	if args[1][0] != '\'' {
		fmt.Println("incorecct type for Decrment")
	}
	
	varname := args[0][1:]
	val, ok := state.VarStorage[varname]
	if !ok {
		fmt.Println("varname: " + varname + " was not found.")
		return
	}
	val = val[1:]
	
	//debug line
	//fmt.Printf("val = %q\n", val)
	
	//check if they can be turned into numbers
	current, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("variable is not a number")
		return
	}

	inc, err := strconv.Atoi(args[1][1:])
	if err != nil {
		fmt.Println("decrement is not a number")
		return
	}
	
	current = current - inc // converts the string into a number were 1: (as strings are arrys of chars) uis everytihng past the "   and then adds it to the vars value
	state.VarStorage[varname] = "\"" + strconv.Itoa(current) // return to string
}