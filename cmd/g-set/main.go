package main

import (
	"encoding/json"
	"fmt"
	"os"
    "sync"
	"time"

	"alan-kuan/dist-sys-practice/pkg/node"
	"alan-kuan/dist-sys-practice/pkg/utils"
)

func main() {
    n := newGSetNode()

    n.On("add", makeAddMessageHandler(n))
    n.On("read", makeReadMessageHandler(n))
    n.On("replicate", makeReplicateMessageHandler(n))

    n.every(5, func () {
        n.crdtLock.RLock()
        messages := n.crdt.read()
        n.crdtLock.RUnlock()

        my_id := n.GetNodeId()
        for _, curr_id := range n.GetNodeIds() {
            if curr_id == my_id {
                continue
            }

            body := replicateMessageBody{
                BaseMessageBody: node.BaseMessageBody{ Type: "replicate" },
                Value: messages,
            }
            raw_body, _ := json.Marshal(body)

            n.Send(curr_id, raw_body)
        }
    })

    n.runPeriodicTasks()

    if err := n.Run(); err != nil {
        fmt.Fprintln(os.Stderr, "An error occurred:", err)
        os.Exit(1)
    }
}

func newGSetNode() *gSetNode {
    n := &gSetNode{
        Node: node.NewNode(),
        crdt: *newGSet(),
        crdtLock: new(sync.RWMutex),
    }
    return n
}

func (n *gSetNode) every(interval int, callback func()) {
    task := periodicTask{ interval, callback }
    n.periodicTasks = append(n.periodicTasks, task)
}

func (n *gSetNode) runPeriodicTasks() {
    for i := range n.periodicTasks {
        go func (task *periodicTask) {
            for {
                task.callback()
                time.Sleep(time.Duration(task.interval) * time.Second)
            }
        }(&n.periodicTasks[i])
    }
}

func makeAddMessageHandler(n *gSetNode) node.Handler {
    return func (msg node.Message) error {
        recv_body, err := utils.DecodeMessageBody[addMessageBody](msg.Body)
        if err != nil {
            return err
        }

        n.crdtLock.Lock()
        n.crdt.add(recv_body.Element)
        n.crdtLock.Unlock()

        resp_body := node.BaseMessageBody{
            Type: "add_ok",
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }
        if err := n.Reply(msg, map_resp_body); err != nil {
            return err
        }

        return nil
    }
}

func makeReadMessageHandler(n *gSetNode) node.Handler {
    return func (msg node.Message) error {
        n.crdtLock.RLock()
        messages := n.crdt.read()
        n.crdtLock.RUnlock()

        resp_body := readMessageBody{
            BaseMessageBody: node.BaseMessageBody{ Type: "read_ok" },
            Value: messages,
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }
        if err := n.Reply(msg, map_resp_body); err != nil {
            return err
        }

        return nil
    }
}

func makeReplicateMessageHandler(n *gSetNode) node.Handler {
    return func (msg node.Message) error {
        recv_body, err := utils.DecodeMessageBody[replicateMessageBody](msg.Body)
        if err != nil {
            return err
        }

        n.crdtLock.Lock()
        n.crdt.merge(recv_body.Value)
        n.crdtLock.Unlock()

        return nil
    }
}
