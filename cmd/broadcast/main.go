package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"alan-kuan/dist-sys-practice/pkg/node"
	"alan-kuan/dist-sys-practice/pkg/utils"
)

func main() {
    n := newBroadcastNode()

    n.On("topology", makeTopologyHandler(n))
    n.On("broadcast", makeBroadcastHandler(n))
    n.On("read", makeReadHandler(n))

    if err := n.Run(); err != nil {
        fmt.Fprintln(os.Stderr, "An error occurred:", err)
        os.Exit(1)
    }
}

func newBroadcastNode() *broadcastNode {
    n := &broadcastNode{
        Node: node.NewNode(),
        messages: make(map[any]struct{}),
        messagesLock: new(sync.Mutex),
    }
    return n
}

func makeTopologyHandler(n *broadcastNode) node.Handler {
    return func (msg node.Message) error {
        recv_body, err := utils.DecodeMessageBody[topologyMessageBody](msg.Body)
        if err != nil {
            return err
        }
        n.neighbors = recv_body.Topology[n.GetNodeId()]
        n.Log("My neighbors: %v\n", n.neighbors)

        resp_body := node.BaseMessageBody{
            Type: "topology_ok",
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }

        return n.Reply(msg, map_resp_body)
    }
}

func makeBroadcastHandler(n *broadcastNode) node.Handler {
    return func (msg node.Message) error {
        // 1. reply ok
        resp_body := node.BaseMessageBody {
            Type: "broadcast_ok",
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }
        if err := n.Reply(msg, map_resp_body); err != nil {
            return err
        }

        // 2. check the message
        recv_body, err := utils.DecodeMessageBody[broadcastMessageBody](msg.Body)
        if err != nil {
            return err
        }

        n.messagesLock.Lock()
        if _, exists := n.messages[recv_body.Message]; exists {
            n.messagesLock.Unlock()
            return nil
        }
        n.messages[recv_body.Message] = struct{}{}
        n.messagesLock.Unlock()

        n.Log("Received message: %v\n", recv_body.Message)

        // 3. gossip to neighbors
        gossip_body := broadcastMessageBody{
            BaseMessageBody: node.BaseMessageBody{ Type: "broadcast" },
            Message: recv_body.Message,
        }
        map_gossip_body, err := utils.EncodeMessageBodyToMap(&gossip_body)
        if err != nil {
            return err
        }

        unacked_neighbors := make(map[string]struct{})
        unacked_lock := new(sync.Mutex)

        for _, neighbor := range n.neighbors {
            if neighbor == msg.Src {
                continue
            }
            unacked_neighbors[neighbor] = struct{}{}
        }

        // retry loop
        for len(unacked_neighbors) > 0 {
            for dest := range unacked_neighbors {
                n.Rpc(dest, map_gossip_body, func(msg node.Message) error {
                    recv_body, err := utils.DecodeMessageBody[node.BaseMessageBody](msg.Body)
                    if err != nil {
                        return err
                    }

                    if recv_body.Type == "broadcast_ok" {
                        unacked_lock.Lock()
                        delete(unacked_neighbors, dest)
                        unacked_lock.Unlock()
                    }

                    return nil
                })
            }

            time.Sleep(time.Second)
        }

        return nil
    }
}

func makeReadHandler(n *broadcastNode) node.Handler {
    return func (msg node.Message) error {
        n.messagesLock.Lock()
        messages := utils.MapToSlice(n.messages)
        n.messagesLock.Unlock()

        resp_body := readMessageBody{
            BaseMessageBody: node.BaseMessageBody{ Type: "read_ok" },
            Messages: messages,
        }
        map_resp_body, err := utils.EncodeMessageBodyToMap(&resp_body)
        if err != nil {
            return err
        }

        return n.Reply(msg, map_resp_body)
    }
}
