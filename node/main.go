package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

func NewNode() *Node {
    n := &Node{
        nextMsgIdLock: new(sync.Mutex),
        logLock: new(sync.Mutex),
        handlers: make(map[string]Handler),
        wg: new(sync.WaitGroup),
    }

    n.handlers["init"] = n.initHandler
    n.handlers["echo"] = n.echoHandler
    n.handlers["topology"] = n.topologyHandler

    return n
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
                n.log("An error occurred when handling '%s' message: %s\n",
                    recv_msg.Body.Type, err)
            }
        }()
        n.wg.Wait()
    }

    return nil
}

func (n *Node) log(msg string, a ...any) {
    n.logLock.Lock()
    fmt.Fprintf(os.Stderr, msg, a...)
    n.logLock.Unlock()
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
