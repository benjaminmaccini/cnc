package router

import (
	"cnc/pkg/utils"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// This works by overriding the default args (os.Args[1:]) with the supplied string
// Serves as the hook for the text-mode entry
// Capture the output in the designated writer
func ExecuteInternal(args []string, writer io.Writer, cmd *cobra.Command) {
	cmd.SetArgs(args)
	cmd.SetOut(writer)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Route a string to a command entry point
func RouteStringAsCommand(body string, cmd *cobra.Command) string {
	// Set an in-memory pipe for capturing the output
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// If 'cnc' prefixes the request then process it as a command, importantly
	// add the internal flag, which prevents certain commands from executing
	if strings.HasPrefix(body, "cnc") {
		command := utils.ParseStringAsCommand(body)

		// Strip the leading `cnc` to avoid errors
		ExecuteInternal(
			append(command[1:], "--internal"),
			w,
			cmd,
		)

		w.Close()
		out, err := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		if err != nil {
			log.Fatal(err)
		}

		return string(out)
	} else if strings.HasPrefix(body, "?") {
		ExecuteInternal(
			[]string{"chat", fmt.Sprintf("--query=\"%s\"", body[1:])},
			w,
			cmd,
		)
		w.Close()
		out, err := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		if err != nil {
			log.Fatal(err)
		}

		return string(out)
	} else {
		return "`cnc help` or ask a question by prefixing `?<question>`"
	}
}
