package main

import (
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	
	"oscarsgoofysite/OCON/commands"
	"oscarsgoofysite/OCON/state"
	"oscarsgoofysite/OCON/returncommands"
	"oscarsgoofysite/OCON/functions"
	"oscarsgoofysite/OCON/updating"
	
	"os/exec"//for fancy imports
	
	"strconv"
	//"io"
	"runtime"
	"syscall"
	"unsafe" //:-0 
	"slices"
	//"time"
	//"errors"
	//"regexp"
)


//set this to that tag vershion thingy 
const Version = "v3.7.1"
/*
added errors
made better if statements
new easteregg
added wait
added start and run
added allvars
userinput

added tonum
*/
//major, minor, patch
const versionName = "Fih"


func main() {
	
	fmt.Println("🐟")
	//setup the commands in state
	state.Commands = &Commmands
		
	
	//test if there needs to be an update here
	updating.IfNecesitoActualizar(Version)
	
	if len(os.Args) < 2 {
		var i string
		fmt.Println("This is OCON a language created by Oscar! (see more at: https://github.com/oscar366/Ocon or https://ocon.oscarsgoofy.site)")
		fmt.Println(`To execute an .ocon file do: "ocon execute {path}"`)
		fmt.Println("this is a command line tool use it in cmd.exe")
		fmt.Println("You have to add ocon to path yourself")
		fmt.Println("")
		fmt.Println("type something and click enter to leave")
		fmt.Scan(&i)
		fmt.Println(i)
		return
	}
	//test if ocon is in path else place it there
	
	switch os.Args[1] {
	case "execute":
		
		if len(os.Args) < 3 {
			fmt.Println("missing file path")
			os.Exit(2)
			return
		}
		readFile(os.Args[2])
	case "🐟":
	
		fish()
		os.Exit(0)
	case "help":
		fmt.Println("This is OCON a language created by Oscar! To see more about it see the github (https://github.com/oscar366/Ocon) or the website (https://ocon.oscarsgoofy.site)")
		fmt.Println(`To execute an .ocon file do: "ocon execute {path}"`)
		fmt.Println("")
		fmt.Println("`ocon 🐟` ??")
		fmt.Println("`ocon execute {path}` execute ocon file")
		fmt.Println("`ocon update` get new ocon updates (see more about it with `ocon update help`)")
		fmt.Println("`ocon version` to see what version you are using")
		os.Exit(0)
	case "version":
		fmt.Println("V" + Version + " - " + versionName)
	case "update":
		if len(os.Args) == 3 {
			if os.Args[2] == "help" {
				fmt.Println("Running `ocon update` it downloads the next exe file for ocon in the same folder (does not deleate the prev update)")
			}
		} else if len(os.Args) < 3 && len(os.Args) > 1 {
			fmt.Println("downloading now")
			updating.ReinstallOcon()
		}
	default:
		fmt.Println("Put in an a real input see `ocon help`")
		os.Exit(2)
	}
}


//ifNecesitoActualizar
func readFile(path string) {
	
	
	
	lines := []string{}
	
	file, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	
	scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
		lines = append(lines, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
	
	if lines[0] == "!debugOFF" || (len(lines) >= 2 && lines[1] == "!debugOFF") {
		state.DebugMode = false;
	} else {
		fmt.Println("executing:" + path)//same as using `if state.DebugMode`
	}
	
	if lines[0] == "!windowOFF" || (len(lines) >= 2 && lines[1] == "!windowOFF") {
		fmt.Println("Not currently built")
	}
	state.DocumentData = lines //to save to whole file so other packages can use it
	
	executeall(lines)
}

/*
cmds:
wait - wait for a sertan amount of time#
🥚 - easter egg #
start - start a program and contuine #
run - run a program and wait until it stops #
userinput - put user input into a var 
unlink - remove var from memory  #
allvars - print all vars#

*/
type Command = func([]string) error
var Commmands = map[string]Command{//typo i cant fix commmands
	//other
    "echo":      commands.EchoCmd,
	"#":	emptycommand,//comments
	
	//vars
    "var":       commands.VarCmd,
    "increment": commands.IncrementCmd,
	"decrement": commands.DecremntCmd, 
	//"program": commands.ProgramCmd,
	
	//sections
	"§": emptycommand, //we aculy dont need these commands but to not confuse the inteprted we keep them
	"sec": emptycommand,
	"goto": commands.GotoCmd,
	
	//conditionals
	"if": commands.If,
	"func": emptycommand,
	"end": emptycommand,
	
	"import": emptycommand,
	
	//for loops
	"for": emptycommand,//forr cuz for exists
	"endfor": commands.Endfor,
	"try": commands.Try,
	"unlink": commands.UnLinkVar,
	//eastereggs
	"🥚": easteregg,
	"🐟": easteregg,
	"🫃": easteregg,
	
	//other0
	"wait": commands.Wait,
	"start": commands.Start,
	"run": commands.Run,
	"allvars": commands.Allvars,
	
	//userinput
	/*"ifkeypressed": commands.Ifkeyispressed,
	"waituntilkeypressed": commands.Waituntilkeypressed,*/
	
}




func executeall(lines []string) {
//before runing the program get the positons of all the sections
if state.DebugMode {
fmt.Println(" ")
fmt.Println("doing setup")
fmt.Println("===============================")
fmt.Println(" ")
}
//more complacted for loop :-0
//this is the preprossesing things
//for exaple puting the sections into the array and adding the import statments
for i, line := range lines {
    parts := strings.Fields(line) // sperate by " "
	//sectionals
    if len(parts) >= 2 && (parts[0] == "§" || parts[0] == "sec") { //if there is more then 2 inputs and its a section then run
        state.SectionList[parts[1][1:]] = i //add to the map the name and line number of the section
    }
	//if its an import statment
	if len(parts) >= 2 && parts[0] == "import" {
		if parts[1] != "f" {
			importfile := parts[1][1:]

			dat, err := os.ReadFile(importfile)
			if err != nil {
				panic(fmt.Sprintf("%v. at: %d", err, i))
			}

			newlines := strings.Split(string(dat), "\n")

			// Replace the import line with the imported lines thank you stakoverflow
			lines = append(
				lines[:i],
				append(newlines, lines[i+1:]...)...,
			)

			// Move past the newly inserted lines
			i += len(newlines) - 1
		} else {
			//were doing fancy imports in here if were doing them agean 
				//check as
				//import f command {path} as x
				//0      1  2       3     4  5
				if parts[4] != "as" {
					panic("Error: 'as' not found in fancy import. possibly at: " + strconv.Itoa(i))
				}
				
				if parts[2] == "command" {
					//asname, path
					if state.DebugMode {
					fmt.Println("")
					fmt.Println("debug fancy imports:")
					fmt.Println("____________________")
					}
					addToCommands(parts[5][1:], parts[3][1:])
				} else if parts[2] == "returncommand" {
					returncommands.AddToReturnCommands(parts[5][1:], parts[3][1:])
				} else {
					panic("Error: not a command or return command in fancy import. I think its at: " + strconv.Itoa(i))
				}
			}
		}
	//the len(parts) is importent but i have no idea why
	if len(parts) >= 2 && parts[0] == "func" {
		if len(parts) < 2 {
			panic("not enogh args for 'func' command roughly at:" + strconv.Itoa(i))
		}


		var n int
		//parts [1][1:] is the func name
		for n < len(lines) && lines[n] != "end @"+parts[1][1:] {
			n++
		}
		endlinenum := n //i know its at + 1 but i removed it cuz we dont need the end anyway
		//i is being wird here
		/*fmt.Printf("FUNC: %s i=%d n=%d len=%d\n",
			parts[1][1:], i, n, len(lines))*/

		codeslice := lines[i+1:endlinenum]
		codestring := strings.Join(codeslice, "\n")
		//get sections:
		sections := make(map[string]int)

		for j, line := range codeslice {
			parts := strings.Fields(line)

			if len(parts) >= 2 && (parts[0] == "sec" || parts[0] == "§") {
				sections[parts[1][1:]] = j
			}
		}
		
		//fmt.Printf("%q\n", strings.Join(codeslice, "\n")) //debug line
		length := endlinenum - (i + 1)
		//here we remove the lines so it does not execute
		//check if out of bounds
		if endlinenum >= len(lines) {
			fmt.Println("error in finding end out of bounds")
			return
		}
		lines = append(lines[:i],lines[endlinenum+1:]...)
		functions.AddToFunctions(codestring, parts[1][1:], length, sections)
	}
	
	
}


if state.DebugMode {
fmt.Println("")
fmt.Println("debug S:")
for key, value := range state.SectionList {
    fmt.Println("Key:", key, "Value:", value)
}
fmt.Println("____________________________")
fmt.Println(" ")
fmt.Println("setup complete")
fmt.Println("====================================")
fmt.Println(" ")
}



//real execution above is setup
	for state.Pointer < len(lines) {
		
		if lines[state.Pointer] == "" {
			state.Pointer += 1//move on
			continue
		}
		if lines[state.Pointer][0] == '@' /*if its a function*/ {
			//name string, lines *[]string
			//& is an address (like in mem) and * is a pointer to an address
			funcCommandWords := strings.Split(lines[state.Pointer], " ")
			functions.ExecuteFunction(funcCommandWords[0][1:], &lines, funcCommandWords[1:])
			continue
		}
		
		
		//rember that 0 is equal to start of a list
		execute(lines[state.Pointer])
		state.Pointer += 1// it does this after goto is executed that is why its n+1
	}
	
	
}









func addToCommands(command string, path string) {
	//add to commands
	if state.DebugMode {
	fmt.Println("adding:" + path + " as:" + command)
	}
	var function Command = func(args []string) error {
		cmd := exec.Command(path, args...)

		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	
	//checker
	Commmands[command] = function
	_, ok := Commmands[command]//check if a key exists in a map
	// If the key exists
	if ok && state.DebugMode {
		// Do something
		fmt.Println("Key is found :)")
	} else {
		fmt.Println("Key is missing from commands.")
	}
}







func execute(prgmstring string) {
	if strings.TrimSpace(prgmstring) == "" {//check if empty or broken
		return
	}
	if prgmstring[0] == '!' {
		return
	}
	//fmt.Printf("Executing: %q\n", prgmstring) debug
	
	words := strings.Split(prgmstring, " ")
	

	
	//replaces vars with var content and does return functions and array bull
	//also called execute preprosses
	words = executePreprocessing(words)
	

	cmd := words[0]
	args := words[1:]
	
	
	//get the function(fn) if it exits in the commands list(ok)
	if fn, ok := Commmands[cmd]; ok {
				err := fn(args)//normal command
				if err != nil {
					fmt.Println("Program was killed. Pointer:" + strconv.Itoa(state.Pointer))
					os.Exit(1)
				}
	} else {
			fmt.Println("Unknown command:", cmd)
	}
}
/*
todo make "_a_ automaticly turn into object: "\" a "
*/
func executePreprocessing(words []string) []string {
		//these are vars for in between
		inreturn := false
		//unset will set later. This is the pos of after the "["
		firstpos := 0
	
		//for return commands with value change
		valueChangeChar := ""
	
		returndepth := 0
		mto := false//to see if there is more then one return command
	
		for i, element := range words {
		if inreturn {
			if element == "[" {
				returndepth++
			}
			
			if element == "]" {
				returndepth--
				if returndepth == 0 {
				//get the return of the return func interp
				returncommand := words[firstpos:i]
				var value = []string{}
				//preprosses inbedded other returns:
				/*fmt.Println("mto =", mto)
				fmt.Println("returncommand =", returncommand)*/

				if mto {
					//fmt.Println("🔥 CALLING preProcessReturn")
					value = preProcessReturn(returncommand)
				} else {
					//fmt.Println("calling normal return")
					value = returncommands.Intrp(returncommand)
				}

				//rebuild the slaice with value inbtween and remove the ret command
				if valueChangeChar != "" {//replace the type
					value[0] = valueChangeChar + value[0] //replace the type
				}
				words = append(
					words[:firstpos - 1],
					append(value, words[i+1:]...)...,
				)
				//inreturn = false
				valueChangeChar = "" // reset it
				}
				
			}
			
			if returndepth > 1 {
				mto = true
			}
			
		}
		
		//Start a new returncommand
		if len(element) > 1 && element[1] == '[' && returndepth == 0 {//for return commands with value change
			inreturn = true
			firstpos = i + 1
			returndepth = 1
			mto = false
			valueChangeChar = string(element[0]) //set to the value before the [ eg "[ would be the "
			//fmt.Printf("i=%d element=%q inreturn=%v firstpos=%d valueChangeChar=%q\n",i, element, inreturn, firstpos, valueChangeChar)
		} else if element == "[" && returndepth == 0 {
			inreturn = true
			returndepth = 1
			firstpos = i + 1
			mto = false
		}




		//vars
		if len(element) > 1 && element[0] == '$' {
			//if it is a var then replaces
			varname := element[1:]
			
			
			/*
			varname := $banana|2
			value := '1,'2,'3,'4,'5
			*/
			isAnArrayWithSelector := false
			pos := -1 //-1 for unknow pois 
			if strings.Contains(varname, "|") {
				pos = strings.Index(element, "|")
				varname = varname[:pos-1]//replace var name -1 to disinclude the |
				isAnArrayWithSelector = true
			}
			//fmt.Println("Varname " + varname)
			
			val, ok := state.VarStorage[varname]
			//if var does NOT exsit
			if (!ok) {
				fmt.Println("Var does not exist. VarName: " + varname)
				return []string{}
			}
			//fmt.Println("Val " + val)
			//get the array value 
			if isAnArrayWithSelector {
				arrayItems := strings.Split(val, ",")//get all the items in an array
				
				//cant use varname here cuz its bean modfyed above
				arrayIndex, err := strconv.Atoi(element[1:][pos:])//get the index and truen it into an interger
				if err != nil {
					fmt.Println("Error in turning array index into a int for a var array. Index: " + varname[pos+1:])
					return []string{}
				}
				val = arrayItems[arrayIndex]//set the value to the item at index
			}
			
			
			
			//replace with var value
			words[i] = val
		}
		if len(element) > 1 && element[0] == '*' && strings.Contains(element, "|") {//eg _'1,'2,'3|2 for an array
			pos := strings.Index(element, "|") + 1 // +1 so | isent counted
			if pos == -1 {
				fmt.Println("there was an error in finding the index of this array :(")
				return []string{}
			}//Itoa
			itemNumber, err := strconv.Atoi(element[pos:])//evryting past the selecter "|"
			if err != nil {
				fmt.Println("Error in converting index into a string. Index = " + element[pos:])
				return []string{}
			}
			list := strings.Split(element[1:pos-1], ",") // all the elements or items in an array
			//fmt.Printf("%v", list)
			words[i] = list[itemNumber]
		}
		
		//for adding strings and arrays
		if i+2 <= len(words)/*check if i+2 is useable*/ && words[i+1] == "+"/*if next "word" is plus */ && ( element[0] == '*' || element[0] == '"' ) /* and im a string or an array */ {
			replaceText := ""
			//get type
			switch element[0] {
				case '"':
					//type check
					if words[i+2][0] != '"' {
						fmt.Println("Error concating strings: type mismatch")
						return []string{}
					}
					replaceText = "\"" + element[1:] + words[i+2][1:]
				case '*':
					//type check
					if words[i+2][0] != '*' {
						fmt.Println("Error concating arrays: type mismatch")
						return []string{}
					}
					myarray := element[1:]
					otherarray := words[i+2][1:]
					replaceText = "*" + myarray + "," + otherarray
				default:
					fmt.Println("Error concating: incorect types")
					return []string{}
			}
			
			// for exaple:
			// echo "banana + "_apple
			// 0      1     2     3
			//        i     i+1 i+2
			//    |
			//    v
			// echo "banana_apple
			// 0      1
			end := i + 3
			if end > len(words) {
				end = len(words)
			}

			words = append(
				words[:i],
				append([]string{replaceText}, words[end:]...)...,
			)
		}
		
		
	}
	
	
	return words
}

func preProcessReturn(str []string) []string {
	
	/*fmt.Println("PREPROCESS RETURN CALLED")
    fmt.Println("str =", str)*/
	
	//create the list of return commands i need to do and there args
	//operations := []string{}
	/*
		name  cmd&args
		0      [a µ1]
		1      [b µ2 '2]
		2      [c]
	*/
	depth := 0
	data := map[int][]string{}
	for _, item := range str {
		if item == "[" {
			data[depth] = append(data[depth], fmt.Sprintf("µ%d", depth+1))
			depth++
			data[depth] = []string{}
			continue
		}
		if item == "]" {
			depth--
			continue
		}
		data[depth] = append(data[depth], item)
	}
	//fmt.Println(data)
	//create a array of all the commands returns
	ans := map[int][]string{}
	
	//return []string{"Igotnothin"}
	for d := len(data) - 1; d >= 0; d-- {
		values := data[d]
		//fmt.Println("processing:", d, values)
		if slices.Contains(values,"µ" + strconv.Itoa(d + 1)) {
			//replace a func
			pos := slices.Index(values,"µ" + strconv.Itoa(d + 1)) //find the positon of object needing replacing
			//replace with ans at [d+1] returns array so we combine them with this
			result := ans[d+1]
			values = append(
				values[:pos],
				append(result, values[pos+1:]...)...,
			) 
		}
		ans[d] = returncommands.Intrp(values)
	}
	return ans[0]
}


func emptycommand(args []string) error {
	return nil //the sections dont do anything but we requrie this so it does not say "unknown command"
}



/*
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⣤⣤⣤⣤⣤⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡟⠁⠀⠀⠙⢿⣿⣿⣿⡿⠋⠀⠀⠀⠀⠙⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⡀⠀⠀⠀⠀⠈⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢨⠁⢠⣾⣶⣦⠀⢸⣿⣿⢠⣾⣿⣶⡀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⠀⢸⣿⣿⣿⠤⠘⠀⠘⠼⣿⣿⣿⡇⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣧⡀⢹⠟⠁⠀⠀⠀⠀⠀⠈⠙⢟⣁⠀⢀⣼⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡟⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠻⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡆⠣⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⡤⠖⠀⠀⣠⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣦⡘⠢⠤⠤⠤⠤⠤⠒⠉⠁⠀⢀⣠⣴⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⠟⠉⠢⣄⣢⠐⣄⠠⣄⢢⣼⠞⠉⠀⠈⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⣿⣿⡟⠀⠀⠀⠀⠉⠙⠚⠓⠊⠉⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣿⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀ ⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀
⠀⠀⠀⡰⠉⠈⠑⠠⢀⢸⣿⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀
⠀⠀⠀⡇⠀⠀⠀⠀⠀⠉⠙⠛⠛⠛⠿⠿⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠉⠀⠀⠀⠘⢿⣿⠀⠀⠀
⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⠟⠉⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⡇⠀⠀⠀
⠀⠀⢰⠃⠀⠀⢠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⠘⣿⣿⣿⣿⣿⡟⠁⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⡇⠀⠀⠀
⠀⡠⠊⠀⢀⠐⡀⠈⠄⠂⡐⠀⢂⠐⠈⠠⢀⠀⠀⢻⡤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣷⠀⢀⠛⡛⢫⠑⡄⢃⡐⢈⠐⡀⠄⠐⠀⠀⠀⠀⠈⢿⡇⠂⠀⠀
⢠⢁⠀⠄⢂⡐⠠⡁⠌⡐⠠⢁⠂⠌⢠⢁⠂⡐⠀⠘⣿⣳⢤⡀⡀⢀⠀⡀⢀⠀⡀⠠⡀⠤⣁⢿⠀⠄⣂⠑⡂⠥⡘⢠⠐⢂⠰⠀⠌⡐⢈⠐⡀⠀⠀⠀⠑⢄⠀⠀
⠈⢧⡘⡐⢂⠤⠑⡠⢁⠆⡁⠆⠌⣂⠁⡂⠌⡐⠀⠀⢹⣿⣷⣧⣝⣢⠱⡰⣈⢆⢡⢃⠴⡱⣌⣾⠈⡐⢠⠘⡠⢁⠆⡡⢘⠠⡁⠎⡐⡈⢄⠢⢀⠡⢀⠈⠂⠀⠑⡀
⠀⠀⠙⢵⣊⠴⡁⢆⠡⢂⠅⡊⠔⡠⠘⢄⠒⡀⢁⠀⠀⢻⣿⣿⣿⣿⣿⣷⣷⣾⣶⣿⣾⣿⣿⣿⠀⠐⡄⠢⢁⠆⡘⢄⠡⢂⠱⢠⠑⡨⢄⠢⣁⠒⡄⢊⠄⣂⢀⡡
⠀⠀⠀⠀⠉⠲⣍⢢⠱⡈⢆⠱⡈⠔⡉⢄⠒⠄⢂⠀⢀⠀⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠀⠂⢄⠣⠌⣂⠱⡈⢆⠡⢊⠄⢣⠐⢢⠑⡄⢣⡘⢆⡳⣬⠞⠁
⠀⠀⠀ ⠀⠀⠈⢣⡞⡰⢈⠆⡱⢈⠔⡨⢘⡈⠆⢌⠀⡐⠨⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠠⢉⠄⢢⠑⡄⢣⠘⡄⠣⢌⢊⡔⡉⢦⠩⡜⣡⢞⡷⠋⠁⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠙⢶⡉⢆⡱⢈⠆⡑⠢⢌⡘⢄⠣⡐⣡⠏⠉⠉⠉⠉⠉⠉⠉⠉⠍⠉⠉⢳⢁⠊⡜⢠⠃⡜⢠⠃⣌⠱⣈⠦⢰⡉⢆⡳⣼⠟⠁⠀⠀ ⠀⠀
⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠹⣖⡰⢃⡜⢄⠳⣠⠚⣌⠖⣥⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣎⠴⡈⢆⠱⡈⢆⠱⡠⢃⠖⣌⠣⣜⢣⠟⠁⠀⠀⠀⠀ ⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠓⢯⡼⣬⣓⣦⣟⣼⡿⠚⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣶⣉⢆⠳⡌⡜⢢⠱⡩⢜⣤⢻⡼⠋⠀⠀ ⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠉⠙⠛⠛⠋⠉⠀ ⠀⠀ ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⢾⣳⣼⣜⣧⣳⡽⣞⠞⠋⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ⠀⠀⠀⠀⠈⠉⣉⢉⣉⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ FROM https://emojicombos.com/linux-ascii-art I DON'T KNOW WHO CREATED IT 
*/
const tux = "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⣤⣤⣤⣤⣤⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡟⠁⠀⠀⠙⢿⣿⣿⣿⡿⠋⠀⠀⠀⠀⠙⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⡀⠀⠀⠀⠀⠈⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢨⠁⢠⣾⣶⣦⠀⢸⣿⣿⢠⣾⣿⣶⡀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⠀⢸⣿⣿⣿⠤⠘⠀⠘⠼⣿⣿⣿⡇⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣧⡀⢹⠟⠁⠀⠀⠀⠀⠀⠈⠙⢟⣁⠀⢀⣼⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡟⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠻⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡆⠣⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⡤⠖⠀⠀⣠⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣦⡘⠢⠤⠤⠤⠤⠤⠒⠉⠁⠀⢀⣠⣴⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⠟⠉⠢⣄⣢⠐⣄⠠⣄⢢⣼⠞⠉⠀⠈⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⣿⣿⡟⠀⠀⠀⠀⠉⠙⠚⠓⠊⠉⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣿⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠇⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡀⠀⠀⠀⠀\n⠀⠀⠀⢀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀\n⠀⠀⠀⡰⠉⠈⠑⠠⢀⢸⣿⣿⣿⣿⣿⣿⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄⠀⠀⠀\n⠀⠀⠀⡇⠀⠀⠀⠀⠀⠉⠙⠛⠛⠛⠿⠿⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⠉⠀⠀⠀⠘⢿⣿⣿⡇⠀⠀⠀\n⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⣿⠟⠉⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⡇⠀⠀⠀\n⠀⠀⢰⠃⠀⠀⢠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⠘⣿⣿⣿⣿⣿⡟⠁⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⡇⠀⠀⠀\n⠀⡠⠊⠀⢀⠐⡀⠈⠄⠂⡐⠀⢂⠐⠈⠠⢀⠀⠀⢻⡤⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣷⠀⢀⠛⡛⢫⠑⡄⢃⡐⢈⠐⡀⠄⠐⠀⠀⠀⠀⠈⢿⡇⠂⠀⠀\n⢠⢁⠀⠄⢂⡐⠠⡁⠌⡐⠠⢁⠂⠌⢠⢁⠂⡐⠀⠘⣿⣳⢤⡀⡀⢀⠀⡀⢀⠀⡀⠠⡀⠤⣁⢿⠀⠄⣂⠑⡂⠥⡘⢠⠐⢂⠰⠀⠌⡐⢈⠐⡀⠀⠀⠀⠑⢄⠀⠀\n⠈⢧⡘⡐⢂⠤⠑⡠⢁⠆⡁⠆⠌⣂⠁⡂⠌⡐⠀⠀⢹⣿⣷⣧⣝⣢⠱⡰⣈⢆⢡⢃⠴⡱⣌⣾⠈⡐⢠⠘⡠⢁⠆⡡⢘⠠⡁⠎⡐⡈⢄⠢⢀⠡⢀⠈⠂⠀⠑⡀\n⠀⠀⠙⢵⣊⠴⡁⢆⠡⢂⠅⡊⠔⡠⠘⢄⠒⡀⢁⠀⠀⢻⣿⣿⣿⣿⣿⣷⣷⣾⣶⣿⣾⣿⣿⣿⠀⠐⡄⠢⢁⠆⡘⢄⠡⢂⠱⢠⠑⡨⢄⠢⣁⠒⡄⢊⠄⣂⢀⡡\n⠀⠀⠀⠀⠉⠲⣍⢢⠱⡈⢆⠱⡈⠔⡉⢄⠒⠄⢂⠀⢀⠀⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠀⠂⢄⠣⠌⣂⠱⡈⢆⠡⢊⠄⢣⠐⢢⠑⡄⢣⡘⢆⡳⣬⠞⠁\n⠀⠀⠀⢀⠀⠀⠈⢣⡞⡰⢈⠆⡱⢈⠔⡨⢘⡈⠆⢌⠀⡐⠨⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠠⢉⠄⢢⠑⡄⢣⠘⡄⠣⢌⢊⡔⡉⢦⠩⡜⣡⢞⡷⠋⠁⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠙⢶⡉⢆⡱⢈⠆⡑⠢⢌⡘⢄⠣⡐⣡⠏⠉⠉⠉⠉⠉⠉⠉⠉⠍⠉⠉⢳⢁⠊⡜⢠⠃⡜⢠⠃⣌⠱⣈⠦⢰⡉⢆⡳⣼⠟⠁⠀⠀⡆⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠄⠀⠀⠹⣖⡰⢃⡜⢄⠳⣠⠚⣌⠖⣥⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣎⠴⡈⢆⠱⡈⢆⠱⡠⢃⠖⣌⠣⣜⢣⠟⠁⠀⠀⠀⠀⠃⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠓⢯⡼⣬⣓⣦⣟⣼⡿⠚⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⣶⣉⢆⠳⡌⡜⢢⠱⡩⢜⣤⢻⡼⠋⠀⠀⢸⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠀⠀⠀⠉⠙⠛⠛⠋⠉⠀⡀⠀⠀⠐⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⢾⣳⣼⣜⣧⣳⡽⣞⠞⠋⠀⠀⠀⠀⠒⠀⠀⠀⠀⠀⠀\n⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠀⠀⠀⠈⠉⣉⢉⣉⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀"
func fish() {
	if runtime.GOOS == "windows" {
		//fmt.Println(tux)
			//kernel32DLL = syscall.NewLazyDLL("kernel32.dll")   
			user32 := syscall.NewLazyDLL("user32.dll") //import needed dll
			messageBox := user32.NewProc("MessageBoxW")//get MessageBoxW func
			//set text and title
			text := syscall.StringToUTF16("🐟 Has Escaped :-0")
			title := syscall.StringToUTF16("Error")
			//execute func
			messageBox.Call(
				0,//no value here cuz we dont need user input
				uintptr(unsafe.Pointer(&text[0])),//set text
				uintptr(unsafe.Pointer(&title[0])),//set title
				0x10, // set the icon to error
			)
	} else {
		fmt.Println(tux)
	}
}

func easteregg(args []string) error {
	if len(args) > 0 {
		return nil
	}
	
	
	if runtime.GOOS == "windows" {
		user32 := syscall.NewLazyDLL("user32.dll")
		createWindow := user32.NewProc("CreateWindowExA")
	
		className, _ := syscall.BytePtrFromString("STATIC")//class STATIC should be a custom class but i dont want to make one
		title, _ := syscall.BytePtrFromString("OCONS window of awesomeness")

		hwnd, _, err := createWindow.Call(//add them all together
			0,
			uintptr(unsafe.Pointer(className)),//becuse this is static the window title will aso have the text inside the window
			uintptr(unsafe.Pointer(title)),
			0x00CF0000,
			100,//x
			100,//y
			800,//width
			600,//hight
			0,
			0,
			0,
			0,
		)
		//print errors
		fmt.Println("HWND:", hwnd)
		fmt.Println("Error:", err)
		
		//set text inside windo	
		text, _ := syscall.BytePtrFromString("OCON is the best! 🐟")
			/*
			so we also call createWindow for text and buttons becuse they are also windows?
			Mainwindo
			|- Textwindo
		*/
		textHwnd, _, err := createWindow.Call(
			0,
			uintptr(unsafe.Pointer(className)),
			uintptr(unsafe.Pointer(text)),
			0x50000000, // WS_CHILD | WS_VISIBLE
			0,//x
			20,//y
			150,
			150,
			hwnd,
			0,
			0, //main window however best practres to do something else?
			0,
		)
		fmt.Println("Text HEND:", textHwnd)
		fmt.Println("Error:", err)
		
		showWindow := user32.NewProc("ShowWindow")
		showWindow.Call (
			hwnd,
			5,//SW_SHOW
		)
		//time.Sleep(10 * time.Second)
		//window call vars
		getMessage := user32.NewProc("GetMessageA")
		translateMessage := user32.NewProc("TranslateMessage")
		dispatchMessage := user32.NewProc("DispatchMessageA")
		/*postQuitMessage := user32.NewProc("PostQuitMessage")
		defWindowProc := user32.NewProc("DefWindowProcA")*/
		type POINT struct {//needed in MSG
			X int32
			Y int32
		}

		type MSG struct {//the thing we send to the program to see if its working and stuff
			Hwnd     uintptr
			Message  uint32
			WParam   uintptr
			LParam   uintptr
			Time     uint32
			Pt       POINT
			LPrivate uint32
		}

		var msg MSG
		stop := false
		//for loop to call window
		for stop == false{
			result, _, err := getMessage.Call(
				uintptr(unsafe.Pointer(&msg)),
				0, 0, 0,
			)
	
			if int32(result) == -1 {
				fmt.Println("GetMessage failed:", err)
				break
			}
	
			if result == 0 {
				break
			}
	
			translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	} else {
		fmt.Println(tux)
	}
	/*
@echo OFF
color 0a
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⡀⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⣾⣿⣿⣷⣶⣦⣀⣀⠀⠀⢀⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣻⣷⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⡿⠋⠙⠛⠿⣿⣿⣿⣿⣾⣿⣿⣿⣿⣶⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⢶⣄⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠃⠀⠀⠀⠀⠀⠈⠉⠉⠉⠉⠁⢹⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠿⣮⡿⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣾⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣦⣀⡀⠀⠀⠀⠀⠀⠀⠀⢠⣴⣶⡄
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⡿⢿⣿⣿⣦⠀⠀⠀⠀⠀⠀⠈⠛⠛⠁
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠋⠉⢸⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣷⢠⣾⣿⣿⠀⠀⠀⠀⠀⠀⢠⣶⣶⡆
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣶⣾⣿⣧⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣼⣿⣿⣿⣿⡿⠋⠀⠀⠀⠀⠀⠀⠈⠛⠛⠃
echo ⠀⠀⠀⣀⣀⣀⣀⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣉⣛⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⡛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⣶⢿⣷⠀
echo ⢀⣴⣿⣿⣿⣿⣿⣿⣿⣷⣄⠀⠀⠀⠀⠀⣀⣴⣶⣿⣿⣿⣿⣿⣿⣿⠿⠿⠿⠿⠟⠛⠛⠉⠉⠉⠉⠉⠙⠛⠻⠿⣿⣿⣶⣄⠀⠀⠀⠀⠀⠀⠀⠀⢉⣉⠁⠀
echo ⠘⣿⣿⣿⡄⠀⠀⠀⢙⣿⣿⣆⠀⢀⣴⣿⣿⣿⠿⠛⠋⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⢿⣿⣷⡄⠀⠀⠀⠀⠀⢸⡏⢿⡇⠀
echo ⠀⠀⢻⣿⣿⡄⠀⠀⠀⢻⣿⣿⢰⣿⣿⡿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⣦⡀⠀⠀⠀⠹⣿⣿⣆⠀⠀⠀⠀⠀⠉⠉⠀⠀
echo ⠀⠀⠀⢿⣿⣷⠀⠀⠀⠘⣿⣿⣿⣿⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠛⠁⠀⠀⠀⠀⢹⣿⣿⡄⠀⠀⢰⣶⣶⡄⠀⠀
echo ⠀⠀⠀⠸⣿⣿⡆⠀⠀⠀⠘⠿⠿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣇⣀⣀⠸⠷⠟⠃⠀⠀
echo ⠀⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣶⣄⠀⠀⠀⠀⠀⠀⠀⣀⣀⣠⣤⣴⣴⣾⣿⣿⣿⣟⣿⡇⠀⠀⠀⠀⠀
echo ⠀⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⢿⣿⣿⣶⣶⣶⣾⣿⣿⣿⣿⣿⡿⠿⠟⠛⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⠀⠀⢰⣿⣿⡇⠀⠀⠀⠀⢠⣿⣿⣷⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠛⠛⠛⠛⠛⠉⠉⠁⠀⠀⠀⠀⢀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⠀⢀⣾⣿⣿⠁⠀⠀⠀⠀⣸⣿⣿⢿⣿⣿⣷⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣀⣠⣴⣾⣿⡿⠃⠀⠀⠀⠀⠀⠀⠀⠀
echo ⢰⣾⣿⣿⡿⠃⠀⠀⠀⠀⠀⣽⣿⣿⠀⠀⠙⠻⢿⣿⣿⣶⣶⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣤⣶⣶⣶⣶⣶⣶⣾⣿⣿⣿⣿⡿⠟⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠘⣿⣿⡏⠀⠀⠀⠀⠀⠀⣸⣿⣿⡏⠀⠀⠀⠀⠀⠉⠛⠛⠻⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠿⠛⠛⠛⠛⠛⠛⠋⠛⠉⠀⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⢻⣿⣿⠀⠀⠀⢀⣠⣶⣿⣿⠏⠀⠀⠀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⢸⣿⣿⣷⣾⣿⣿⣿⡿⠛⠁⢀⣴⣿⣿⣿⣿⣿⣿⣿⣿⣦⢠⣶⣾⣶⣷⡄⠀⠀⠀⢀⣤⣤⣤⣤⣄⡀⠀⠀⣰⣶⣦⣤⣤⣤⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⠀⠉⠛⠉⠙⠙⠉⠁⠀⠀⠀⣿⣿⣿⠏⠉⠀⠉⠙⣿⣿⣿⣿⣿⣿⠛⠛⠁⠀⠀⣼⣿⣿⡿⠿⢿⣿⣿⣆⡀⢹⣿⣿⣿⣿⣿⣿⣦⡄⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠻⣿⣿⣶⣤⣴⣶⣾⣿⣿⠟⣿⣿⡏⠀⠀⠀⠀⠀⣿⣿⣟⠀⠀⠀⢹⣿⣿⡇⢸⣿⣿⡏⠀⠘⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠻⠿⠿⠿⠿⠟⠛⠁⠀⣿⣿⣿⣦⣤⣶⣶⡄⠙⣿⣿⣶⣶⣾⣿⣿⡿⠃⠈⣿⣿⣷⠀⠀⠸⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀
echo ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠻⠿⠿⠿⠿⠃⠀⠙⠛⠛⠛⠛⠛⠉⠀⠀⠀⠙⠛⠁⠀⠀⠀⠙⠛⠋⠀⠀⠀⠀⠀⠀⠀
pause


HWND CreateWindowA(
  [in, optional] LPCSTR    lpClassName,
  [in, optional] LPCSTR    lpWindowName,
  [in]           DWORD     dwStyle,
  [in]           int       x,
  [in]           int       y,
  [in]           int       nWidth,
  [in]           int       nHeight,
  [in, optional] HWND      hWndParent,
  [in, optional] HMENU     hMenu,
  [in, optional] HINSTANCE hInstance,
  [in, optional] LPVOID    lpParam
);
	*/
return nil
}