package main

import (
	"fmt"
	"os"

	"alan-kuan/dist-sys-practice/pkg/node"
	"alan-kuan/dist-sys-practice/pkg/utils"
)

func main() {
    n := node.NewNode()

    n.On("echo", func (msg node.Message) error {
        recv_body, err := utils.DecodeMessageBody[echoMessageBody](msg.Body)
        if err != nil {
            return err
        }

        resp_body := echoMessageBody{
            BaseMessageBody: node.BaseMessageBody{ Type: "echo_ok" },
            Echo: recv_body.Echo,
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }

        return n.Reply(msg, map_resp_body)
    })

    if err := n.Run(); err != nil {
        fmt.Fprintln(os.Stderr, "An error occurred:", err)
        os.Exit(1)
    }
}
