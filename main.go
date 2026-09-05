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
	
	"os/exec"//for fancy imports
	"net/http"//for updates
	"encoding/json"//for github api
	"strconv"
	"io"
)
//set PATH=%PATH%;C:\path\to\your\install\directory

//set this to that tag vershion thingy 
const version = "v1.2.1"//ver 1, 2 new fetures, 1 debug fix from prev vir
//major, minor, patch
const versionName = "semi-stable"

var necesitoactualizar = false


func main() {
	fmt.Println("🐟")
	
	
	//test if there needs to be an update here
	ifNecesitoActualizar()
	
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
		if len(os.Args) == 3 {
			if os.Args[2] == "help" {
				fmt.Println("Running `ocon update` it downloads the next exe file for ocon in the same folder (does not deleate the prev update)")
			}
		} else if len(os.Args) < 3 && len(os.Args) > 1 {
			fmt.Println("downloading now")
			reinstallOcon()
		}
	default:
		fmt.Println("Put in an a real input see `ocon help`")
	}
}

func fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
//reinstall the file for updates
func reinstallOcon() {
	if !necesitoactualizar {
		fmt.Println("tú no necesitas actualizar (you don't need an update) 🐟")
		return
	}

	// Get the latest release
	data, err := fetch("https://api.github.com/repos/oscar366/ocon/releases/latest")
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}

	// Get assets_url
	var release struct {
		AssetsURL string `json:"assets_url"`
	}

	if err := json.Unmarshal(data, &release); err != nil {
		fmt.Println("Error parsing release:", err)
		return
	}

	// Get the release assets
	data2, err := fetch(release.AssetsURL)
	if err != nil {
		fmt.Println("Error getting release assets:", err)
		return
	}

	// Parse assets
	var assets []ReleaseAsset

	if err := json.Unmarshal(data2, &assets); err != nil {
		fmt.Println("Error parsing assets:", err)
		return
	}

	if len(assets) == 0 {
		fmt.Println("No release assets found.")
		return
	}

	// Display assets
	fmt.Println("Which one do I install? (type number):")

	for i, asset := range assets {
		fmt.Printf("%d: %s\n", i, asset.Name)
	}

	var userInput string
	fmt.Scan(&userInput)

	num, err := strconv.Atoi(userInput)
	if err != nil {
		fmt.Println("Invalid number:", err)
		return
	}

	if num < 0 || num >= len(assets) {
		fmt.Println("Out of bounds")
		return
	}

	// Download selected asset
	data, err = fetch(assets[num].BrowserDownloadURL)
	if err != nil {
		fmt.Println("Download failed:", err)
		return
	}

	err = os.WriteFile("ocon", data, 0755)
	if err != nil {
		fmt.Println("Failed to write file:", err)
		return
	}

	fmt.Println("==================")
	fmt.Println("Steps After Update")
	fmt.Println("==================")
	fmt.Println("Ocon has downloaded the new version to the same location as the older version.")
	fmt.Println("1. Find the newly downloaded file.")
	fmt.Println("2. Add the appropriate file extension for your operating system.")
	fmt.Println("   - Windows: .exe")
	fmt.Println("   - Linux: no extension required")
	fmt.Println("3. Delete the older version of Ocon.")
	fmt.Println("4. Run the new version.")
}


type Update struct {
	TagName string `json:"tag_name"`
	Name string `json:"name"`
	Body string `json:body`
}
//if needs update
func ifNecesitoActualizar() {
	url := "https://api.github.com/repos/oscar366/Ocon/releases/latest"

	data, err := fetch(url)
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}

	var u Update

	if err := json.Unmarshal(data, &u); err != nil {
		fmt.Println("Error parsing update JSON:", err)
		return
	}

	if u.TagName != version {
		fmt.Println("\033[31mYou need an update.\033[0m")
		fmt.Println("Current:", version)
		fmt.Println("Latest:", u.TagName)
		necesitoactualizar = true
	}
}

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
	
	if lines[0] == "!debugOFF" || lines[1] == "!debugOFF" {
		state.DebugMode = false;
	} else {
		fmt.Println("executing:" + path)//same as using `if state.DebugMode`
	}
	
	if lines[0] == "!windowOFF" || lines[1] == "!windowOFF" {
		fmt.Println("Not currently built")
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
	"func": emptycommand,
	"end": emptycommand,
	
	"import": emptycommand,
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
			panic("not enogh args for 'func' command at roughly at:" + string(i))
		}
		
		var n int
		//parts [1][1:] is the func name
		for lines[n] != "end @" + parts[1][1:] && n < (len(lines) + 1) {
			n++
		}
		endlinenum := n //i know its at + 1 but i removed it cuz we dont need the end anyway
		//i is being wird here
		codeslice := lines[i+1:endlinenum]
		codestring := strings.Join(codeslice, "\n")
		//fmt.Printf("%q\n", strings.Join(codeslice, "\n")) //debug line
		length := endlinenum - (i + 1)
		//here we remove the lines so it does not execute
		lines = append(lines[:i],lines[endlinenum+1:]...)
		functions.AddToFunctions(codestring, parts[1][1:], length)
	}
	//for loops
	if len(parts) >= 2 && parts[0] == "for" {
		fmt.Println("not currently on oscar is working on it please be paishont")
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
		
		if(len(element) > 1 && element[0] == '$') {
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