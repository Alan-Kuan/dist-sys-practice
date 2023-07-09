package main

import (
	"fmt"
	"os"

	"alan-kuan/dist-sys-practice/node"
)

func main() {
    n := node.NewNode()

    if err := n.Run(); err != nil {
        fmt.Fprintln(os.Stderr, "An error occurred:", err)
        os.Exit(1)
    }
}
