package udwConsole

import (
	"fmt"
	"github.com/tachyon-protocol/udw/udwJson"
	"os"
	"sort"
	"strings"
	"sync"
)

var VERSION = ""

type Command struct {
	Name       string
	Desc       string
	Runner     func()
	FuncV2     interface{}
	Hidden     bool
	CompleteFn CompletionFn
}

func (action *Command) runAction() {
	if action.Runner != nil {
		action.Runner()
	} else if action.FuncV2 != nil {
		MustRunCommandLineFromFuncV2(action.FuncV2)
	} else {
		panic("[udwConsole.Main] action not runable")
	}
}

type CommandGroup struct {
	commandMap map[string]Command
}

func NewCommandGroup() *CommandGroup {
	return &CommandGroup{commandMap: map[string]Command{}}
}

func (g *CommandGroup) Main() {
	actionName := ""
	if len(os.Args) >= 2 {
		actionName = os.Args[1]
	}
	lowerActionName := strings.ToLower(actionName)

	action, exist := g.commandMap[lowerActionName]
	if !exist {
		fmt.Println("command [" + actionName + "] not found.(case insensitive)")
		g.Help()
		if lowerActionName != "" {
			os.Exit(-1)
		}
		return
	}

	os.Args = os.Args[1:]

	action.runAction()
}

func (g *CommandGroup) AddCommand(action Command) *CommandGroup {
	name := strings.ToLower(action.Name)
	_, exist := g.commandMap[name]
	if exist {
		panic("command " + action.Name + " already defined.(case insensitive)")
	}
	g.commandMap[name] = action
	return g
}

func (g *CommandGroup) defaultHelp() {
	os.Stdout.WriteString("Usage: ")
	nameList := make([]string, len(g.commandMap))
	i := 0
	for name := range g.commandMap {
		nameList[i] = name
		i++
	}
	sort.Strings(nameList)
	for i := 0; i < len(nameList); i++ {
		lowerName := nameList[i]
		action := g.commandMap[lowerName]
		if action.Hidden {
			continue
		}
		os.Stdout.WriteString("\t" + action.Name + "\n")
	}
}

func (g *CommandGroup) subCommand() {
	for name := range g.commandMap {
		fmt.Printf("%v ", name)
	}
}

func (g *CommandGroup) Help() {
	action, exist := g.commandMap["help"]
	if exist {
		action.runAction()
		return
	} else {
		g.defaultHelp()
	}
}

func (g *CommandGroup) HasCommand(name string) bool {
	name = strings.ToLower(name)
	_, exist := g.commandMap[name]
	return exist
}

func (g *CommandGroup) AddCommandWithName(name string, f interface{}) *CommandGroup {
	simpleFunc, ok := f.(func())
	if ok {
		return g.AddCommand(Command{
			Name:   name,
			Runner: simpleFunc,
		})
	}
	mustEnsureValidFuncV2(f)
	return g.AddCommand(Command{
		Name:   name,
		FuncV2: f,
		CompleteFn: CreateCompletion(CreateCompletionReq{
			ReflectArgOfFunc: f,
		}),
	})

}

var gDefaultCommandGroup *CommandGroup
var gDefaultCommandGroupInit sync.Once

func GetDefaultCommandGroup() *CommandGroup {
	gDefaultCommandGroupInit.Do(func() {
		gDefaultCommandGroup = NewCommandGroup()
	})
	return gDefaultCommandGroup
}

func (g *CommandGroup) CompleteCmd(args []string) (waitSelectList []string) {
	udwJson.MustWriteFileIndent(`/tmp/1.json`, args)
	if len(args) < 2 {
		return nil
	}
	if len(args) == 2 {
		for _, cmd := range g.commandMap {
			if cmd.Hidden {
				continue
			}
			waitSelectList = append(waitSelectList, cmd.Name)
		}
		sort.Strings(waitSelectList)
	} else {
		subCmdName := strings.ToLower(args[1])
		subCmd, ok := g.commandMap[subCmdName]
		if ok && subCmd.CompleteFn != nil {
			_, waitSelectList = subCmd.CompleteFn(args)
		}
	}
	return waitSelectList
}

func Main() {
	group := GetDefaultCommandGroup()
	if VERSION != "" && group.HasCommand("version") == false {
		group.AddCommand(Command{
			Name:   "version",
			Runner: version,
		})
	}
	if group.HasCommand("help") == false {
		group.AddCommand(Command{
			Name:   "help",
			Runner: group.defaultHelp,
		})
	}
	if group.HasCommand("subcommand") == false {
		group.AddCommand(Command{
			Name:   "subCommand",
			Runner: group.subCommand,
		})
	}
	group.Main()
}

func AddCommand(action Command) *CommandGroup {
	return GetDefaultCommandGroup().AddCommand(action)
}

func AddCommandWithName(name string, runner interface{}) *CommandGroup {
	return GetDefaultCommandGroup().AddCommandWithName(name, runner)
}

func version() {
	fmt.Println(VERSION)
}
