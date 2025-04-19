package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ncruces/zenity"
	"go.bug.st/serial"
	"golang.org/x/term"
)

var commonBaudrates = []string{"4800", "9600", "19200", "38400", "57600", "115200", "230400", "460800", "921600", "Custom"}

func main() {

	var ports = []string{}
	var err error

	for {

		if ports, err = serial.GetPortsList(); err != nil {
			zenity.Error(err.Error())
			os.Exit(1)
		}

		if len(ports) == 0 {
			err = zenity.Warning("No active COM Ports found!", zenity.ExtraButton("Rescan"), zenity.OKLabel("Cancel"))

			if err == zenity.ErrExtraButton {
				continue
			} else {
				os.Exit(1)
			}

		} else {
			break
		}

	}

	port, err := zenity.List("Select COM Port", ports)

	if err != nil {
		os.Exit(0)
	}

	baudrate, err := zenity.List("Select Baudrate", commonBaudrates, zenity.DefaultItems("9600"))

	if err != nil {
		os.Exit(0)
	}

	if baudrate == "Custom" {

		for {
			baudrate, err = zenity.Entry("Set Baudrate", zenity.EntryText("9600"))

			if err != nil {
				os.Exit(0)
			}

			if _, err := strconv.Atoi(baudrate); err != nil {
				zenity.Error(fmt.Sprintf("'%s' isn't a numeric value!", baudrate))
			} else {
				break
			}

		}

	}

	nBaudrate, _ := strconv.Atoi(baudrate)

	comPort, err := serial.Open(port, &serial.Mode{BaudRate: nBaudrate})

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if _, err := io.Copy(os.Stdout, comPort); err != nil {
			log.Fatal(err)
		}
		log.Println("channel closed.")
	}()

	term.MakeRaw(int(os.Stdin.Fd()))

	if _, err := io.Copy(comPort, os.Stdin); err != nil {
		log.Fatal(err)
	}
	log.Println("channel closed.")

}
