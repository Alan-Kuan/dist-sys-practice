package main

import (
	"fmt"
	"os"

	"alan-kuan/dist-sys-practice/message"
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

    n.On("echo", func (msg message.Message) error {
        err := n.Reply(msg, message.MessageBody{
            Type: "echo_ok",
            Echo: msg.Body.Echo,
        })
        if err != nil {
            return err
        }
        return nil
    })

    if err := n.Run(); err != nil {
        handleErr(err)
    }
}
