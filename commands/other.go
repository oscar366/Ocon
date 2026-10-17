package commands

import (
    "os/exec"
	"oscarsgoofysite/OCON/helpers"
	"strconv"
	"time"
	"fmt"
	"bufio"
	"os"
	"strings"
)
//wait a few secs
func Wait(args []string) error {
	if len(args) != 1 {
		return helpers.NewError("Incorrect args for wait")
	}
	value, err := strconv.Atoi(args[0][1:])
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(value)/*custom type*/ * time.Millisecond)
	return nil
}


//execute other programs
func Start(args []string) error {
    if len(args) == 0 {
        return helpers.NewError("start: no program specified")
    }

    cmd := exec.Command(args[0][1:], args[1:]...)
    return cmd.Start()
}

func Run(args []string) error {
    if len(args) == 0 {
        return helpers.NewError("run: no program specified")
    }

    cmd := exec.Command(args[0][1:], args[1:]...)
    return cmd.Run()
}




//userinput
/*
func Waituntilkeypressed(args []string) error {
    reader := bufio.NewReader(os.Stdin)

    for {
        input, err := reader.ReadString('\n')
        if err != nil {
            return err
        }

        // No key specified: accept any input
        if len(args) == 0 || args[0] == "" {
            return nil
        }

        // Stop when specified key is found
        if strings.Contains(input, args[0][1:]) {
            return nil
        }

        fmt.Print(input)
    }
}

func Ifkeyispressed(args []string) error {
    if len(args) < 3 {
        return helpers.NewError("ifkeyispressed needs more then the amount of args given")
    }

    reader := bufio.NewReader(os.Stdin)

    input, err := reader.ReadString('\n')
    if err != nil {
        return err
    }

    // Stop when specified key is found
    if strings.Contains(input, args[0][1:]) {
        //KEY IS FOUND
        return If(append([]string{"~true"}, args[1:]...))
    }

    fmt.Print(input)
    return nil
}*/