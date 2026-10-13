package commands

import "fmt"
import "strings"
import "oscarsgoofysite/OCON/helpers"

func EchoCmd(args []string) error {
    if len(args) == 0 {
        return helpers.NewError("no args")
    }

    text := args[0]
	text = strings.ReplaceAll(text, "_", " ")
    if len(text) > 0 && (text[0] == '"' || text[0] == '\'') {
        fmt.Println(text[1:])
    } else if len(text) > 0 && (text[0] == '~') {
		if text[1:] == "true" {
			fmt.Println("✅")
		} else {
			fmt.Println("❌")
		}
    } else {
		 return helpers.NewError("please put an accepted type into echo :). [" + args[0] + "]")
	}
	return nil
}