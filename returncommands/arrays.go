package returncommands

import "strings"
import "strconv"
import "fmt"


func length(args []string) []string {//expected args = _'1,'2,'3,'4
	//to little
	if len(args) < 1 {
		fmt.Println("not enogh args for length command")
		return []string{"\"Error: " + "not enogh args for length command"}
	}
	//to meany
	if len(args) > 1 {
		fmt.Println("to meny args for length command")
		return []string{"\"Error: " + "to meny args for length command"}
	}
	//juuuust right
	
	//split the list
	list := strings.Split(args[0][1:], ",")
	langth := strconv.Itoa(len(list))
	return []string{langth}
}