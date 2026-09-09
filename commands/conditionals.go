package commands

import (
	"fmt"
	//"oscarsgoofysite/OCON/state"
	//"strconv"
)


//90% of the coding is mine the rest is ai :(. 
//However all ideas and methouds of doing are my ideas
func If(args []string) {
	// args does NOT contain "if"
	//
	// Example:
	// if ~true goto "main else continue
	//
	// args:
	// 0 = ~true
	// 1 = goto
	// 2 = "main
	// 3 = else
	// 4 = continue
	//
	// Another possible form:
	//
	// if ~false continue else goto "main
	//
	// args:
	// 0 = ~false
	// 1 = continue
	// 2 = else
	// 3 = goto
	// 4 = "main

	if len(args) == 0 {
		fmt.Println("Error: not enough args for if")
		return
	}

	// The first argument should be ~true or ~false.
	boolValue := args[0][1:]

	// If the boolean is true, execute the command
	// immediately after it.
	if boolValue == "true" {

		// true + goto
		if args[1] == "goto" {
			GotoCmd([]string{args[2]})
			return
		}

		// true + continue
		if args[1] == "continue" {
			return
		}
	}

	// If the boolean is false, skip the first command
	// and execute whatever comes after "else".
	if boolValue == "false" {

		// false + goto/continue + else
		if args[2] == "else" {

			// false goto "x else ...
			if args[3] == "goto" {
				GotoCmd([]string{args[4]})
				return
			}

			// false continue else ...
			if args[3] == "continue" {
				return
			}
		}

		// false + continue + else
		//
		// This is different because "else" is at args[2].
		//
		// 0 = ~false
		// 1 = continue
		// 2 = else
		// 3 = goto/continue
		// 4 = "label
		if args[1] == "continue" && args[2] == "else" {

			if args[3] == "goto" {
				GotoCmd([]string{args[4]})
				return
			}

			if args[3] == "continue" {
				return
			}
		}

		return
	}

	// The boolean wasn't true or false.
	fmt.Println("Error: neither true nor false for bool was: " + args[0])
}
