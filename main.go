package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	logFilename       = ".wa.txt"
	userReadWritePerm = 0600
	timeFormat        = "2006-01-02 15:04:05 -0700"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	fmt.Println("WA: wake-up logger")

	t := time.Now()
	buf := make([]byte, 0, len(timeFormat)+1)
	buf = t.AppendFormat(buf, timeFormat)
	buf = append(buf, '\n')
	fmt.Printf("%s", buf)

	home := os.Getenv("HOME")
	if home == "" {
		return errors.New("cannot locate HOME directory")
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
