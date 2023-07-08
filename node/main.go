package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
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

type Handler func(Message) error

type Node struct {
    nodeId      string
    nodeIds     []string
    nextMsgId   int
    lock        *sync.Mutex
    handlers    map[string]Handler
}

func NewNode() (*Node, error) {
    n := &Node{
        nodeId: "",
        nodeIds: nil,
        nextMsgId: 0,
        lock: new(sync.Mutex),
        handlers: map[string]Handler{},
    }

    err := n.On("init", func (msg Message) error {
        n.nodeId = msg.Body.NodeId
        n.nodeIds = msg.Body.NodeIds

        return n.Reply(msg, MessageBody{ Type: "init_ok" })
    })
    if err != nil {
        return nil, err
    }

    return n, nil
}

func (n *Node) Run() error {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        var recv_msg Message

        err := json.Unmarshal(scanner.Bytes(), &recv_msg)
        if err != nil {
            return err
        }

        handler, ok := n.handlers[recv_msg.Body.Type]
        if !ok {
            return fmt.Errorf("No handler for message type '%s'",
                recv_msg.Body.Type)
        }

        if err := handler(recv_msg); err != nil {
            return err
        }
    }

    return nil
}

func (n *Node) On(msg_type string, handler Handler) error {
    if _, exists := n.handlers[msg_type]; exists {
        return fmt.Errorf("Handler for this message type already exists.")
    }
    n.handlers[msg_type] = handler
    return nil
}

func (n *Node) Reply(recv_msg Message, resp_body MessageBody) error {
    resp_body.InReplyTo = &recv_msg.Body.MsgId
    return n.Send(recv_msg.Src, resp_body)
}

func (n *Node) Send(dest string, body MessageBody) error {
    body.MsgId = n.nextMsgId
    n.nextMsgId++

    resp_msg := Message{
        Src: n.nodeId,
        Dest: dest,
        Body: body,
    }

    raw_resp_msg, err := json.Marshal(resp_msg)
    if err != nil {
        return err
    }

    fmt.Println(string(raw_resp_msg))

    return nil
}
