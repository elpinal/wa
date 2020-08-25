package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	logFilename       = ".wa.txt"
	userReadWritePerm = 0600
	timeFormat        = "2006-01-02 15:04:05 -0700"
)

var showFlag = flag.Bool("show", false, "show log and exit")

func main() {
	flag.Parse()

	if err := run(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("given %d argument(s), but expected no argument", len(args))
	}

	if *showFlag {
		return show()
	}

	fmt.Println("WA: wake-up logger")

	t := time.Now()
	buf := make([]byte, 0, len(timeFormat)+1)
	buf = t.AppendFormat(buf, timeFormat)
	buf = append(buf, '\n')
	fmt.Printf("%s", buf)

	home, err := getHome()
	if err != nil {
		return err
	}

	file := filepath.Join(home, logFilename)
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, userReadWritePerm)
	if err != nil {
		return fmt.Errorf("cannot open log: %w", err)
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to update log: %w", err)
	}

	fmt.Println("Updated", file)
	return nil
}

func show() error {
	home, err := getHome()
	if err != nil {
		return err
	}

	file := filepath.Join(home, logFilename)

	cmd := exec.Command("less", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getHome() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("cannot locate HOME directory")
	}
	return home, nil
}
