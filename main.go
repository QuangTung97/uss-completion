package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"learn-go-flags/completion"
)

type BaseOptions struct {
	Name string `short:"n" long:"name" description:"a name"`
	Age  int    `long:"age" description:"age value"`
}

type RootCommand struct {
	BaseOptions

	Sub    SubCommand `command:"sub" description:"sub command"`
	Volume VolumeCmd  `command:"vol" description:"volume command"`
}

var rootCmd RootCommand

type SubCommand struct {
	Args struct {
		URI completion.UriValue `positional-arg-name:"<uri>" description:"uss URI string"`
	} `positional-args:"yes" required:"yes"`
}

var _ flags.Commander = &SubCommand{}

func (cmd *SubCommand) Execute(args []string) error {
	fmt.Println("Sub command called", args)
	fmt.Println("URI Value:", cmd.Args.URI)
	fmt.Printf("Base Options: %+v\n", rootCmd.BaseOptions)
	return nil
}

type VolumeCmd struct{}

func printCompletionList(items []flags.Completion) {
	for _, v := range items {
		if strings.HasSuffix(v.Item, completion.NoSpace) {
			v.Item = strings.TrimSuffix(v.Item, completion.NoSpace)
			fmt.Println(v.Item)
		} else {
			fmt.Println(v.Item + " ") // add space to the end
		}
	}
}

func main() {
	// normalMain()
	simpleCompletion()
}

func simpleCompletion() {
	completion.WriteToLog("%+v\n", os.Args[1])
}

func normalMain() {
	parser := flags.NewParser(&rootCmd, flags.Default)
	parser.CompletionHandler = printCompletionList

	_, err := parser.Parse()
	if err != nil {
		var parseErr *flags.Error
		if errors.As(err, &parseErr) {
			if errors.Is(parseErr.Type, flags.ErrHelp) {
				return
			}
		}
	}
}
