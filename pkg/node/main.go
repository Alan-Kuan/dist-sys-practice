package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
    "sync"

    "alan-kuan/dist-sys-practice/pkg/utils"
)

func NewNode() *Node {
    n := &Node{
        nextMsgIdLock: new(sync.Mutex),
        logLock: new(sync.Mutex),
        handlers: make(map[string]Handler),
        callbacks: make(map[int]Handler),
        wg: new(sync.WaitGroup),
    }

    n.On("init", func (msg Message) error {
        recv_body, err := utils.DecodeMessageBody[initMessageBody](msg.Body)
        if err != nil {
            return err
        }

        n.nodeId = recv_body.NodeId
        n.nodeIds = recv_body.NodeIds

        resp_body := BaseMessageBody {
            Type: "init_ok",
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }

        return n.Reply(msg, map_resp_body)
    })

    return n
}

func (n *Node) GetNodeId() string {
    return n.nodeId
}

func (n *Node) GetNodeIds() []string {
    return n.nodeIds
}

func (n *Node) Log(msg string, a ...any) {
    n.logLock.Lock()
    fmt.Fprintf(os.Stderr, msg, a...)
    n.logLock.Unlock()
}

func (n *Node) On(msg_type string, handler Handler) error {
    if _, exists := n.handlers[msg_type]; exists {
        return fmt.Errorf("Handler of this message type already exists.")
    }
    n.handlers[msg_type] = handler
    return nil
}

func (n *Node) Run() error {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        var recv_msg Message

        if err := json.Unmarshal(scanner.Bytes(), &recv_msg); err != nil {
            return err
        }

        recv_body, err := utils.DecodeMessageBody[BaseMessageBody](recv_msg.Body)
        if err != nil {
            return err
        }

        var handler Handler

        if recv_body.InReplyTo != nil {
            // handle a reply message
            handler, _ = n.callbacks[*recv_body.InReplyTo]
        }
        if handler == nil {
            // handle a new message
            handler, _ = n.handlers[recv_body.Type]
        }
        if handler == nil {
            return fmt.Errorf("No handler for message type '%s'",
                recv_body.Type)
        }

        n.wg.Add(1)
        go func() {
            defer n.wg.Done()

            if err := handler(recv_msg); err != nil {
                n.Log("An error occurred when handling '%s' message: %s\n",
                    recv_body.Type, err)
            }
        }()
    }
    if err := scanner.Err(); err != nil {
        return err
    }

    n.wg.Wait()

    return nil
}

func (n *Node) Rpc(dest string, map_body *map[string]any, handler Handler) error {
    n.nextMsgIdLock.Lock()
    msg_id := n.nextMsgId
    n.nextMsgId++
    n.nextMsgIdLock.Unlock()

    n.callbacks[msg_id] = handler
    (*map_body)["msg_id"] = msg_id

    raw_body, err := json.Marshal(*map_body)
    if err != nil {
        return err
    }

    return n.Send(dest, raw_body)
}

func (n *Node) Reply(recv_msg Message, map_resp_body *map[string]any) error {
    recv_body, err := utils.DecodeMessageBody[BaseMessageBody](recv_msg.Body)
    if err != nil {
        return err
    }

    (*map_resp_body)["in_reply_to"] = *recv_body.MsgId

    n.nextMsgIdLock.Lock()
    (*map_resp_body)["msg_id"] = n.nextMsgId
    n.nextMsgId++
    n.nextMsgIdLock.Unlock()

    raw_resp_body, err := json.Marshal(*map_resp_body)
    if err != nil {
        return err
    }

    return n.Send(recv_msg.Src, raw_resp_body)
}

func (n *Node) Send(dest string, raw_resp_body json.RawMessage) error {
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
