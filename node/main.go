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
        messages: make(map[any]struct{}),
        messagesLock: new(sync.Mutex),
    }

    n.handlers["init"] = n.initHandler
    n.handlers["echo"] = n.echoHandler
    n.handlers["topology"] = n.topologyHandler
    n.handlers["broadcast"] = n.broadcastHandler
    n.handlers["read"] = n.readHandler

    return n
}

func (n *Node) Run() error {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        var recv_msg Message

        if err := json.Unmarshal(scanner.Bytes(), &recv_msg); err != nil {
            return err
        }

        recv_body, err := decodeMessageBody[BaseMessageBody](recv_msg.Body)
        if err != nil {
            return err
        }

        handler, ok := n.handlers[recv_body.Type]
        if !ok {
            return fmt.Errorf("No handler for message type '%s'",
                recv_body.Type)
        }

        n.wg.Add(1)
        go func() {
            defer n.wg.Done()

            if err := handler(recv_msg); err != nil {
                n.log("An error occurred when handling '%s' message: %s\n",
                    recv_body.Type, err)
            }
        }()
    }
    n.wg.Wait()

    return nil
}

func (n *Node) log(msg string, a ...any) {
    n.logLock.Lock()
    fmt.Fprintf(os.Stderr, msg, a...)
    n.logLock.Unlock()
}

func (n *Node) reply(recv_msg Message, map_resp_body *map[string]any) error {
    recv_body, err := decodeMessageBody[BaseMessageBody](recv_msg.Body)
    if err != nil {
        return err
    }

    (*map_resp_body)["in_reply_to"] = *recv_body.MsgId

    n.nextMsgIdLock.Lock()
    (*map_resp_body)["msg_id"] = n.nextMsgId
    n.nextMsgId++
    n.nextMsgIdLock.Unlock()

    raw_resp_body, err := json.Marshal(map_resp_body)
    if err != nil {
        return err
    }

    return n.send(recv_msg.Src, raw_resp_body)
}

func (n *Node) send(dest string, raw_resp_body json.RawMessage) error {
    resp_msg := Message{
        Src: n.nodeId,
        Dest: dest,
        Body: raw_resp_body,
    }

    raw_resp_msg, err := json.Marshal(resp_msg)
    if err != nil {
        return err
    }

    fmt.Println(string(raw_resp_msg))

    return nil
}
