package executeExecute

/*
this just holds a single func used for executing in other spots without importing main i have to put it here in this
folder to stop compiler errors about this happing with imports: main -> commands -> main
*/

// state
import (
	"fmt"
	"oscarsgoofysite/OCON/state"
	//"oscarsgoofysite/OCON/commands"
)

func ExecuteCommand(words []string) error {
	if len(words) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := words[0]
	args := words[1:]

	if fn, ok := (*state.Commands)[cmd]; ok {
		return fn(args)
	}

	return fmt.Errorf("Unknown command: %s", cmd)
}