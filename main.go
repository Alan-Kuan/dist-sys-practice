package main

import (
	"fmt"
	"os"

	"alan-kuan/dist-sys-practice/node"
)

func handleErr(err error) {
    fmt.Fprintln(os.Stderr, "An error occurred:", err)
    os.Exit(1)
}

func main() {
    n, err := node.NewNode()
    if err != nil {
        handleErr(err)
    }

    if err := n.Run(); err != nil {
        handleErr(err)
    }
}
