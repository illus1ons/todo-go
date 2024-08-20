package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/illus1ons/todo-go/todo"
)

const (
	todoFile = "./todos.json"
)

func main() {
	add := flag.Bool("add", false, "할 일 추가")
	complete := flag.Int("complete", 0, "할 일 완료 표시")
	del := flag.Int("del", 0, "할 일 삭제")
	list := flag.Bool("list", false, "할 일 전체 목록")
  deleteAll := flag.Bool("deleteAll", false, "할 일 전체 삭제")
	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}
	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}

  case *deleteAll:
     err := todos.DeleteAll()
     if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
			os.Exit(1)
		}

	case *list:
		todos.Print()
	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil
}
