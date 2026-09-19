package returncommands

import "strings"
import "strconv"
import "fmt"
//import "slices"


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

func appennd(args []string) []string {
	if len(args) < 1 {
		fmt.Println("not enough args for append")
		return []string{"\"Error: not enough args for append"}
	}

	result := ""
	//loop for args and add them all up
	for _, arg := range args {
		result += "," + arg[1:]
	}

	return []string{"*" + result[1:]}
}


func prepend(args []string) []string {
	if len(args) < 2 {
		fmt.Println("mpt enough args for prepend")
		return []string{"\"Error: not enough args for prepend"}
	}
	result := ""
	//reverce for loop from: https://stackoverflow.com/questions/13190836/is-there-a-way-to-iterate-over-a-slice-in-reverse-in-go
	for i := len(args)-1; i >= 0; i-- {
		result += "," + args[i][1:]
	}
	return []string{"*" + result[1:]}
}

func remove(args []string) []string {
	//args *'1,'2,'3 '2 -> *'1,'3
	//check
	if len(args) < 2 {
		fmt.Println("not enogh args for remove returncommand")
		return []string{"\"Error: not enogh args for remove returncommand"}
	} else if len(args) > 2 {
		fmt.Println("Too meany args in returncommand remove")
		return []string{"\"Error: Too meany args in returncommand remove"} 
	}
	//turn first arg into a slice
	array := strings.Split(args[0][1:], ",")
//	fmt.Println(array)
	//get the number
	removeIndex, err := strconv.Atoi(args[1][1:])
	//fmt.Println(args[1][1:])
	if err != nil {
		fmt.Println("Error in converting remove index into a number")
		return []string{"\"Error: converting remove index into a number"}
	}
	array = append(array[:removeIndex], array[removeIndex+1:]...)
	//convert array back into a string
	arraystr := ""
	for _,item := range array {
		arraystr += ","+item
	}
	
	return []string{"*"+arraystr[1:]}
}

func contains(args []string) []string {
	if len(args) < 2 {
		fmt.Println("Too little args for contains")
		return []string{"\"Error: Too little args for contains"}
	}

	if len(args) > 2 {
		fmt.Println("Too many args for contains")
		return []string{"\"Error: Too many args for contains"}
	}

	array := strings.Split(args[0][1:], ",")

	target := args[1]

	for _, v := range array {
		if v == target {
			return []string{"~true"}
		}
	}

	return []string{"~false"}
}

func indexOf(args []string) []string {
		if len(args) < 2 {
		fmt.Println("Too little args for indexOf")
		return []string{"\"Error: Too little args for indexOf"}
	}

	if len(args) > 2 {
		fmt.Println("Too many args for indexOf")
		return []string{"\"Error: Too many args for indexOf"}
	}

	array := strings.Split(args[0][1:], ",")

	target := args[1]

	for i, v := range array {
		if v == target {
			index := strconv.Itoa(i)
			return []string{"'" + index}
		}
	}
	
	//if the for loop does not termante
	fmt.Println(args[1] + "Is not conataned in " + args[0])
	return []string{"\"Error: " + args[1] + "Is not conataned in " + args[0]}
}

func join(args []string) []string {
	if len(args) < 2 {
		fmt.Println("not enogh args in join")
		return []string{"\"Error: Not enogh args in join"}
	}
	return []string{args[0] + "," + args[1][1:]}
}
