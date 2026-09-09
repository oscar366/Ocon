package functions

import (
    "fmt"
    "strings"
	"strconv"
    "oscarsgoofysite/OCON/state"
)

type Function struct {
    Value          string
    CommandsLength int
	Sections       map[string]int
}

var Functionslist = make(map[string]Function)

func AddToFunctions(str string, name string, commandsLength int, sections map[string]int) {
    funct := Function{
        Value:          str,
        CommandsLength: commandsLength,
		Sections: sections,
    }

    Functionslist[name] = funct
}

func replaceArgs(line string, args []string) string {
    for i, arg := range args {
        placeholder := "%" + strconv.Itoa(i+1)
        line = strings.ReplaceAll(line, placeholder, arg)
    }

    return line
}

func ExecuteFunction(name string/*name of the func*/, lines *[]string/*pointer to the lines of executeall*/, funcArgs []string/*function arugments*/) {
    val, ok := Functionslist[name] // gets the func
	
    if !ok {//if not ok
        fmt.Println("Function not found:", name)
        return
    }
	
    functionLines := strings.Split(val.Value, "\n")//split it into lines
	
    //replace %1, %2, %3, etc.
    for index, value := range functionLines {
        functionLines[index] = replaceArgs(value, funcArgs)
    }

    pointer := state.Pointer//set the pointer to a var
	
    // Insert the function lines at the current position.
    *lines = append(
        (*lines)[:pointer],
        append(
            functionLines,
            (*lines)[pointer+1:]...,
        )...,
    )//no idea what this does wrote it a cople of days ago think it edits the lines
	
	// Add the function's sections to the global section list
	for section, position := range val.Sections {
		state.SectionList[section] = pointer + position
	}
	
    // Move pointer back so the first function line gets executed.
    //state.Pointer--
}