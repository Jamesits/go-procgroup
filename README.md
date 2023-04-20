# go-procgroup

Group your child processes and kill them at once.

## Usage

The object created by `group.NewCmd()` can be used as a drop-in of `exec.Cmd`.

```go
package main

import (
	"github.com/jamesits/go-procgroup"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// create a new process group
	g, err := procgroup.NewGroup()
	if err != nil {
		panic(err)
	}
	
	// create a new process in the group and launch it
	c, err := g.NewCmd()
	if err != nil {
		panic(err)
	}
	c.Path = filepath.Join(os.Getenv("WINDIR"), "system32", "cmd.exe")
	c.Args = []string{"cmd.exe", "/c", "notepad.exe"} // instruct CMD to create a subprocess
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Wait() // you still need to `Wait` for it so that Golang runtime does not leak memory

	// wait a while
	time.Sleep(3 * time.Second)
	
	// kill all the processes and their child processes in this group
	err = g.Terminate(0)
	if err != nil {
		panic(err)
	}
}
```

## Caveats

This library does not offer any level of security, and process grouping is not a security boundary. There is no
guarantee that a malicious process cannot break away from its process group. The library only offers a better
abstraction for managing long-running service processes and their own child processes.

To ensure a clean cleanup, the child processes you start will be killed automatically if the parent process dies.
