package main

import (
	"bufio"
    "encoding/json"
	"fmt"
	"os"
)

type MessageBody struct {
    MsgId       int         `json:"msg_id"`
    Type        string      `json:"type"`
    NodeId      string      `json:"node_id,omitempty"`
    NodeIds     []string    `json:"node_ids,omitempty"`
    InReplyTo   *int        `json:"in_reply_to,omitempty"`
    Echo        string      `json:"echo,omitempty"`
}

type Message struct {
    Src     string      `json:"src"`
    Dest    string      `json:"dest"`
    Body    MessageBody `json:"body"`
}

func reply(recv_msg Message, resp_body MessageBody) {
    resp_msg := Message{
        Src: recv_msg.Dest,
        Dest: recv_msg.Src,
        Body: resp_body,
    }
    raw_resp_msg, _ := json.Marshal(resp_msg)
    fmt.Println(string(raw_resp_msg))
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    next_msg_id := 0

    for scanner.Scan() {
        var recv_msg Message

        json.Unmarshal(scanner.Bytes(), &recv_msg)

        resp_body := MessageBody{
            MsgId: next_msg_id,
            InReplyTo: &recv_msg.Body.MsgId,
        }
        next_msg_id++

        switch recv_msg.Body.Type {
        case "init":
            resp_body.Type = "init_ok"
            reply(recv_msg, resp_body)
        case "echo":
            resp_body.Type = "echo_ok"
            resp_body.Echo = recv_msg.Body.Echo
            reply(recv_msg, resp_body)
        }
    }
}
