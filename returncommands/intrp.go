package returncommands

import (
	"fmt"
	//"math/rand/v2"
	"strconv"
	"oscarsgoofysite/OCON/mathfuncs"
	"oscarsgoofysite/OCON/state"
	"log"
	"os/exec"
)

type Function func(args []string)  []string

var returncommands = map[string]Function{
		//sums
		"add": add,
		"mlt": mlt,
		"div": div,
		"sub": sub,
		
		//bools
		"iseven": iseven,
		"isodd": isodd,
		"isless": isless,
		"isgreater": isgreater,
		"isequal": isequal,
		"not": not,
		
		//more arithmetic
		"mod": mathfuncs.Modulo,//
		//
		"sin": mathfuncs.Sine,//#
		"cos": mathfuncs.CoSine,//#
		"tan": mathfuncs.Tangent,//#
		//
		"round": mathfuncs.Round,//#
		"abs": mathfuncs.AbsoluteValue,//#
		"log": mathfuncs.Logarithm,//#
}

func Intrp(args []string) []string {
	if (len(args) == 0) {
		fmt.Println("nothing in the return command")
		return []string{"\"Error: empty return command"}
	}
	val, ok := returncommands[args[0]]
	// If the key exists
	if !ok {
		fmt.Println("That return command does not exist: " + args[0])
		//slice := []string{s}
		//"\"Error: return command does not exist"
		return []string{"\"Error: return command does not exist"}
	}
	result := val(args[1:])
	return result
}

func AddToReturnCommands(command string, path string) {
	//add to return commands
	if state.DebugMode {
	fmt.Println("adding:" + path + " as Return:" + command)
	}
	var function Function = func(args []string) []string {
		out, err := exec.Command(path, args...).Output()

		
		if err != nil {
			log.Fatal(err)
		}
		return []string{string(out)}
	}
	
	//checker
	returncommands[command] = function
	_, ok := returncommands[command]
	//thnks stakoverflow: https://stackoverflow.com/questions/2050391/ddg#2050629
	// If the key exists
	if ok && state.DebugMode {
		// Do something
		fmt.Println("Key is found (in returncommands) :)")
	} else {
		fmt.Println("Key is missing from returncommands.")
	}
}

//all sums
func add(args []string) []string {
	//all functions downward are basicly a copy of this
	
	//if not enogh args are supliyed
	if len(args) < 2 {
		msg := "not enough arguments for add"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}
	//if there was an err converting into a number
	a, err := strconv.Atoi(args[0][1:])
	if err != nil {
		msg := "bad first number: " + args[0]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}
	// same for b
	b, err := strconv.Atoi(args[1][1:])
	if err != nil {
		msg := "bad second number: " + args[1]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}
	//if the numbers grow too big
	if (b > 0 && a > 9223372036854775807-b) ||
		(b < 0 && a < -9223372036854775808-b) {
		msg := "cannot add these numbers, they are too big!"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}
	//Finaly do the math
	return []string{"'" + strconv.Itoa(a+b)}
}

func mlt(args []string) []string {
	if len(args) < 2 {
		msg := "not enough arguments for multiply"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	a, err := strconv.Atoi(args[0][1:])
	if err != nil {
		msg := "bad first number: " + args[0]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	b, err := strconv.Atoi(args[1][1:])
	if err != nil {
		msg := "bad second number: " + args[1]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	result := a * b

	// detects multiplication overflow
	if a != 0 && result/a != b {
		msg := "cannot multiply these numbers, they are too big!"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	return []string{"'" + strconv.Itoa(result)}
}

func div(args []string) []string {
	if len(args) < 2 {
		msg := "not enough arguments for divide"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	a, err := strconv.Atoi(args[0][1:])
	if err != nil {
		msg := "bad first number: " + args[0]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	b, err := strconv.Atoi(args[1][1:])
	if err != nil {
		msg := "bad second number: " + args[1]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}
	//only this one has this check but its a good one
	if b == 0 {
		msg := "cannot divide by zero!"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	return []string{"'" + strconv.Itoa(a/b)}
}

func sub(args []string) []string {

	if len(args) < 2 {
		msg := "not enough arguments for subtract"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	a, err := strconv.Atoi(args[0][1:])
	if err != nil {
		msg := "bad first number: " + args[0]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	b, err := strconv.Atoi(args[1][1:])
	if err != nil {
		msg := "bad second number: " + args[1]
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	// detects subtraction overflow
	if (b > 0 && a < -9223372036854775808+b) ||
		(b < 0 && a > 9223372036854775807+b) {
		msg := "cannot subtract these numbers, they are too big!"
		fmt.Println(msg)
		return []string{"\"Error: " + msg}
	}

	return []string{"'" + strconv.Itoa(a-b)}
}



//boolines
/*
"iseven": iseven,
		"isodd": isodd,
		"isless": isless,
		"isgreater": isgreater,
		"isequal": isequal,
		"isnot": isnot,
*/
func iseven(args []string) []string {
    if len(args) < 1 {
        return []string{"Error: not enough args for iseven"}
    }

    n, err := strconv.Atoi(args[0][1:])
    if err != nil {
        return []string{"Error: invalid number for iseven"}
    }

    if n%2 == 0 {
        return []string{"~true"}
    }

    return []string{"~false"}
}

func isodd(args []string) []string {
    if len(args) < 1 {
        return []string{"Error: not enough args for isodd"}
    }

    n, err := strconv.Atoi(args[0][1:])
    if err != nil {
        return []string{"Error: invalid number for isodd"}
    }

    if n%2 != 0 {
        return []string{"~true"}
    }

    return []string{"~false"}
}

func isless(args []string) []string {
	a, _ := strconv.Atoi(args[0][1:])
	b, _ := strconv.Atoi(args[1][1:])

	if a < b {
		return []string{"~true"}
	} else {
		return []string{"~false"}
	}
}

func isgreater(args []string) []string {
	a, _ := strconv.Atoi(args[0][1:])
	b, _ := strconv.Atoi(args[1][1:])

	if a > b {
		return []string{"~true"}
	} else {
		return []string{"~false"}
	}
}

func isequal(args []string) []string {
	a, _ := strconv.Atoi(args[0][1:])
	b, _ := strconv.Atoi(args[1][1:])
	if a == b {
		return []string{"~true"}
	} else {
		return []string{"~false"}
	}
}

func not(args []string) []string {
	if (args[0] == "~true") {
		return []string{"~false"}
	} else {
		return []string{"~true"}
	}
}

