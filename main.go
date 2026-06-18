package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sstallion/go-hid"

	"github.com/alesr/krgb/config"
	"github.com/alesr/krgb/keyboard"
	"github.com/alesr/krgb/tui"
	"github.com/alesr/krgb/xerrors"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "panic: %v\n%s\n", r, debug.Stack())
			os.Exit(int(xerrors.Internal))
		}
	}()

	if err := hid.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", xerrors.Wrap(xerrors.HIDInit, "hid init failed", err))
		os.Exit(int(xerrors.HIDInit))
	}

	conn, err := keyboard.Find()
	if err != nil {
		hid.Exit()

		if errors.Is(err, keyboard.ErrNotFound) {
			fmt.Println("Keychron K-series keyboard not found.")
			fmt.Println("Make sure your keyboard is connected via USB.")
			os.Exit(int(xerrors.CodeOf(err)))
		}
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(int(xerrors.CodeOf(err)))
	}

	settings := config.Load()

	p := tea.NewProgram(tui.New(conn.Device, conn.Info, settings))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(int(xerrors.Internal))
	}

	conn.Device.Close()
	hid.Exit()
}
