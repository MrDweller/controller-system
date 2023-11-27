package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	controllersystem "github.com/MrDweller/controller-system/controller-system"
	"github.com/MrDweller/orchestrator-connection/models"
)

type Cli struct {
	controllerSystem controllersystem.ControllerSystem
	running          bool
}

func StartCli(controllerSystem controllersystem.ControllerSystem) {

	var output io.Writer = os.Stdout
	var input *os.File = os.Stdin

	fmt.Fprintln(output, "Starting controller system cli...")

	cli := Cli{
		controllerSystem: controllerSystem,
		running:          true,
	}

	for {
		if cli.running == false {
			fmt.Fprintln(output, "Stopping the controller system!")

			err := controllerSystem.StopControllerSystem()
			if err != nil {
				log.Panic(err)
			}
			break
		}

		fmt.Fprint(output, "enter command: ")

		reader := bufio.NewReader(input)
		input, _ := reader.ReadString('\n')

		commands := strings.Fields(input)
		cli.handleCommand(output, commands)
	}
}

func (cli *Cli) Stop() {
	cli.running = false
}

func (cli *Cli) handleCommand(output io.Writer, commands []string) {
	numArgs := len(commands)
	if numArgs <= 0 {
		fmt.Fprintln(output, errors.New("no command found"))
		return
	}

	command := strings.ToLower(commands[0])

	switch command {
	case "lamp":
		if numArgs == 2 {

			lampStatusNumber, err := strconv.Atoi(commands[1])
			if err != nil {
				fmt.Fprintln(output, err)
				return
			}
			lampStatus := false
			if lampStatusNumber == 1 {
				lampStatus = true
			}

			var controll map[string]any = map[string]any{}
			controll["lampOn"] = lampStatus

			err = cli.controllerSystem.SendControll(models.ServiceDefinition{
				ServiceDefinition: "lamp",
			}, controll)

			if err != nil {
				fmt.Fprintln(output, err)

			}
		}

	case "help":
		fmt.Fprintln(output, helpText)

	case "clear":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

	case "exit":
		cli.Stop()

	default:
		fmt.Fprintln(output, errors.New("no command found"))
	}

}

var helpText = `
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	[ CONTROLLER APPLICATION SYSTEM COMMAND LINE INTERFACE ]
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

COMMANDS:
	command [command options] [args...]

VERSION:
	v1.0
	
COMMANDS:
	lamp <arg>			Toggle the lamp, turn on if arg is '1'
	help				Output this help prompt
	clear				Clear the terminal
	exit				Stop the controller system
`
