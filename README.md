# go-procgroup

Offers the ability to group the child processes and kill each group at once.

## Usage

```go
package main

import (
	"filepath"
	"github.com/jamesits/go-procgroup"
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
	c.Path = filepath.Join(os.Getenv("WINDIR"), "system32", "notepad.exe")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Wait()

	// wait a while
	time.Sleep(3 * time.Second)
	
	// kill all the processes and their child processes in this group
	err = g.Terminate(0)
	if err != nil {
		panic(err)
	}
}
```