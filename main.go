// Copyright (C) 2020 Fabio Del Vigna
//
// This file is part of drbracket.
//
// drbracket is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// drbracket is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with drbracket.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"sync"

	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"github.com/yoskini/drbracket/lib/parser"
)

func walker(p string, files chan<- string) error {
	stat, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("Cannot stat file %s", p)
	}
	switch mode := stat.Mode(); {
	case mode.IsDir():
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Cannot explore path %s", path)
			}
			stat, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("Cannot stat file %s", path)
			}
			if stat.Mode().IsRegular() {
				files <- path
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Cannot walk filepath %s", p)
		}
	case mode.IsRegular():
		files <- p
	}
	return nil
}

func tester(c <-chan string) error {
	for f := range c {
		fh, err := os.Open(f)
		if err != nil {
			return fmt.Errorf("Cannot open file %s", f)
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		lineNum := 0
		parser := parser.NewBracketParser()
		if parser == nil {
			return fmt.Errorf("Cannot instantiate BracketParser")
		}
		for scanner.Scan() {
			lineNum++
			//TODO skip comments
			line := scanner.Text()
			err := parser.ParseLine(lineNum, line)
			if err != nil {
				fmt.Printf(err.Error())
				return err
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}
		if !parser.Empty() {
			b := parser.Top()
			fmt.Printf("Unclosed %v bracket at line: %v, col: %v\n", b.Kind, b.Line, b.Col)
		}
	}
	return nil
}

type Config struct {
	Version bool `short:"v" long:"version" description:"Print version"`
	Args    struct {
		Paths []string
	} `positional-args:"yes" required:"yes"`
}

var config = Config{
	Version: false,
}

func main() {
	var parser = flags.NewParser(&config, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	if config.Version {
		fmt.Printf("Version: 0.0.1")
	}

	fchan := make(chan string, 100)
	wgAll := sync.WaitGroup{}
	wgWalkers := sync.WaitGroup{}
	for _, path := range config.Args.Paths {
		wgAll.Add(1)
		wgWalkers.Add(1)
		go func(p string) {
			err := walker(p, fchan)
			if err != nil {
				logrus.Fatalf(err.Error())
			}
			wgWalkers.Done()
			wgAll.Done()
		}(path)
	}

	wgAll.Add(1)
	go func() {
		tester(fchan)
		wgAll.Done()
	}()

	wgWalkers.Wait()
	close(fchan)
	wgAll.Wait()
}
