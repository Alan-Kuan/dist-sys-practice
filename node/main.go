package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

func NewNode() (*Node, error) {
    var err error

    n := &Node{
        nodeId: "",
        nodeIds: nil,
        nextMsgId: 0,
        nextMsgIdLock: new(sync.Mutex),
        handlers: map[string]Handler{},
        wg: new(sync.WaitGroup),
    }

    if err = n.on("init", n.initHandler); err != nil {
        return nil, err
    }
    if err = n.on("echo", n.echoHandler); err != nil {
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

        n.wg.Add(1)
        go func() {
            defer n.wg.Done()

            if err := handler(recv_msg); err != nil {
                fmt.Fprintf(os.Stderr,
                    "An error occurred when handling '%s' message: %s\n",
                    recv_msg.Body.Type, err)
            }
        }()
        n.wg.Wait()
    }

    return nil
}

func (n *Node) on(msg_type string, handler Handler) error {
    if _, exists := n.handlers[msg_type]; exists {
        return fmt.Errorf("Handler for this message type already exists.")
    }
    n.handlers[msg_type] = handler
    return nil
}

func (n *Node) reply(recv_msg Message, resp_body MessageBody) error {
    resp_body.InReplyTo = &recv_msg.Body.MsgId
    return n.send(recv_msg.Src, resp_body)
}

func (n *Node) send(dest string, body MessageBody) error {
    n.nextMsgIdLock.Lock()
    body.MsgId = n.nextMsgId
    n.nextMsgId++
    n.nextMsgIdLock.Unlock()

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
