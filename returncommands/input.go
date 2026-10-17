package returncommands

import (
    "fmt"
    "strconv"
)

func inputbool(args []string) []string {
    if len(args) != 1 {
        return []string{"only 1 arg needed for inputbool"}
    }

    if args[0] == "~true" {
        fmt.Print("Please input a bool (t/f or y/n): ")
    }

    var value string
    if _, err := fmt.Scan(&value); err != nil {
        return []string{"incorrect type"}
    }

    switch value {
    case "y", "t":
        return []string{"~true"}
    case "n", "f":
        return []string{"~false"}
    default:
        return []string{"incorrect type"}
    }
}

func inputnum(args []string) []string {
    if len(args) != 1 {
        return []string{"only 1 arg needed for inputnum"}
    }

    if args[0] == "~true" {
        fmt.Print("Please input a number: ")
    }

    var value int
    if _, err := fmt.Scan(&value); err != nil {
        return []string{"incorrect type"}
    }

    conv := strconv.Itoa(value)
    return []string{"'" + conv}
}

func inputchar(args []string) []string {
    if len(args) != 1 {
        return []string{"only 1 arg needed for inputchar"}
    }

    if args[0] == "~true" {
        fmt.Print("Please input a single letter: ")
    }

    var value string
    if _, err := fmt.Scan(&value); err != nil {
        return []string{"incorrect type"}
    }

    if len([]rune(value)) != 1 {
        return []string{"incorrect type"}
    }

    return []string{"\"" + value}
}

