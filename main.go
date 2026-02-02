package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"uss-completion/completion"
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
		URI completion.UriAndFile `positional-arg-name:"<uri>" description:"uss URI string"`
	} `positional-args:"yes" required:"yes"`
}

var _ flags.Commander = &SubCommand{}

func (cmd *SubCommand) Execute(args []string) error {
	fmt.Println("Sub command called", args)
	fmt.Println("URI Value:", cmd.Args.URI)
	fmt.Printf("Base Options: %+v\n", rootCmd.BaseOptions)
	return nil
}

type VolumeCmd struct {
	Args struct {
		Files []completion.Filename `positional-arg-name:"<file>" description:"file name"`
	} `positional-args:"yes" required:"yes"`
}

func main() {
	normalMain()
	// simpleCompletion()
}

func simpleCompletion() {
	completion.WriteToLog("%s\n", completion.PrintArray(os.Args))
}

func normalMain() {
	completion.EnableLogging.Store(true)

	if os.Getenv("GO_FLAGS_COMPLETE_URI") != "" {
		completeURIFromArgs()
		return
	}

	parser := flags.NewParser(&rootCmd, flags.Default)
	parser.CompletionHandler = completion.PrintCompletionList

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

func completeURIFromArgs() {
	uri := ""
	if len(os.Args) >= 2 {
		uri = os.Args[1]
	}

	var empty completion.UriAndFile
	items := empty.Complete(uri)
	completion.PrintCompletionList(items)
}
