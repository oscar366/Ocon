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
	"os/exec"
	"strconv"
)
//set PATH=%PATH%;C:\path\to\your\install\directory
const version = "0.0.0"
//major, minor, patch
const versionName = "semi-stable"

var necesitoactualizar = false

func main() {
	fmt.Println("🐟")
	
	
	//test if there needs to be an update here
	
	
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
			return
		}
		readFile(os.Args[2])
		
	case "🐟":
		fmt.Println("🐟")
		for i := 0; i < 9999; i++ {
			fmt.Println(string(i) + ": 🐟")
		}
	case "help":
		fmt.Println("This is OCON a language created by Oscar! To see more about it see the github (https://github.com/oscar366/Ocon) or the website (https://ocon.oscarsgoofy.site)")
		fmt.Println(`To execute an .ocon file do: "ocon execute {path}"`)
		fmt.Println("")
		fmt.Println("`ocon 🐟` ??")
		fmt.Println("`ocon execute {path}` execute ocon file")
		fmt.Println("`ocon update` get new ocon updates (see more about it with `ocon update help`)")
		fmt.Println("`ocon version` to see what version you are using")
	case "version":
		fmt.Println("V" + version + " - " + versionName)
	case "update":
		//necesitoactualizar
		fmt.Println("curently inactive oscar is working on it :)")
		/*if os.Args[2] == "help" {
			fmt.Println("Running `ocon update` it downloads the next exe file for ocon in the same folder (does not deleate the prev update)")
		} else if len(os.Args) < 3 && len(os.Args) > 1 {
			fmt.Println("downloading now")
			
		}*/
	default:
		fmt.Println("Put in an a real input see `ocon help`")
	}
}

func readFile(path string) {
	fmt.Println("executing:" + path)
	
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
	
	
	executeall(lines)
}

type Command func([]string)
var commmands = map[string]Command{//typo i cant fix commmands
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
	
	
	"import": emptycommand,
}


func executeall(lines []string) {
	//before runing the program get the positons of all the sections
fmt.Println(" ")
fmt.Println("doing setup")
fmt.Println("===============================")
fmt.Println(" ")

//more complacted for loop :-0
for i, line := range lines {
    parts := strings.Fields(line) // sperate by " "
    if len(parts) >= 2 && (parts[0] == "§" || parts[0] == "sec") { //if there is more then 2 inputs and its a section then run
        state.SectionList[parts[1][1:]] = i //add to the map the name and line number of the section
    }
	
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
					fmt.Println("")
					fmt.Println("debug fancy imports:")
					fmt.Println("____________________")
					addToCommands(parts[5][1:], parts[3][1:])
				} else if parts[2] == "returncommand" {
					returncommands.AddToReturnCommands(parts[5][1:], parts[3][1:])
				} else {
					panic("Error: not a command or return command in fancy import. I think its at: " + strconv.Itoa(i))
				}
			}
		}
}



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

//real execution above is setup
	for state.Pointer < len(lines) {
		//fmt.Printf("%d: %s\n", state.Pointer, lines[state.Pointer]) debug
	
		//rember that 0 is equal to start of a list
		execute(lines[state.Pointer])
		state.Pointer += 1// dude to this after goto is init the pointer adds one so that is why its n-1
	}
}

func addToCommands(command string, path string) {
	//add to commands
	fmt.Println("adding:" + path + " as:" + command)
	var function Command = func(args []string) {
		cmd := exec.Command(path, args...)

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	
	//checker
	commmands[command] = function
	_, ok := commmands[command]
	//thnks stakoverflow: https://stackoverflow.com/questions/2050391/ddg#2050629
	// If the key exists
	if ok {
		// Do something
		fmt.Println("Key is found :)")
	} else {
		fmt.Println("Key is missing from commands.")
	}
}

func execute(prgmstring string) {
	if strings.TrimSpace(prgmstring) == "" {
		return
	}
	//fmt.Printf("Executing: %q\n", prgmstring) debug
	
	words := strings.Split(prgmstring, " ")
	
//these are vars for in between
inreturn := false
//unset will set later. This is the pos of after the "["
firstpos := 0
	
	//replaces vars with var content and does return functions
	for i, element := range words {
		//this is for the rutrn command
		if inreturn && (element == "]") {
			//get the return of the return func interp
			returncommand := words[firstpos:i]
			value := returncommands.Intrp(returncommand)
			//fmt.Println("debug. Value gotten from interp: ", value[0])
			//rebuild the slaice with value inbtween and remove the ret command
			words = append(
				words[:firstpos - 1],
				append(value, words[i+1:]...)...,
			)
			inreturn = false
		}
		//if there is a "]" part but no starting counterpart
		
		if(element == "[") {
			inreturn = true
			firstpos = i + 1
		}
		
		if(element[0] == '$') {
			//if it is a var then replaces
			varname := element[1:]
			
			val, ok := state.VarStorage[varname]
			//if var does NOT exsit
			if (!ok) {
				fmt.Println("Var does not exist. Or other bug.")
				return
			}
			//replace with var value
			words[i] = val
		}
	}
	
	
	cmd := words[0]
	args := words[1:]
	//get the function(fn) if it exits in the commands list(ok)
	if fn, ok := commmands[cmd]; ok {
			fn(args)
	} else {
			fmt.Println("Unknown command:", cmd)
	}
}

func emptycommand(args []string) {
	return //the sections dont do anything but we requrie this so it does not say "unknown command"
}