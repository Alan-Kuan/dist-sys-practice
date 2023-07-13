package node

import (
	"sync"
	"time"
)

func (n *Node) initHandler(msg Message) error {
    recv_body, err := decodeMessageBody[InitMessageBody](msg.Body)
    if err != nil {
        return err
    }

    n.nodeId = recv_body.NodeId
    n.nodeIds = recv_body.NodeIds

    resp_body := BaseMessageBody {
        Type: "init_ok",
    }
    map_resp_body, err := encodeMessageBodyToMap(&resp_body)
    if err != nil {
        return err
    }

    return n.reply(msg, map_resp_body)
}

func (n *Node) echoHandler(msg Message) error {
    recv_body, err := decodeMessageBody[EchoMessageBody](msg.Body)
    if err != nil {
        return err
    }

    resp_body := EchoMessageBody{
        BaseMessageBody: BaseMessageBody{ Type: "echo_ok" },
        Echo: recv_body.Echo,
    }
    map_resp_body, err := encodeMessageBodyToMap(&resp_body)
    if err != nil {
        return err
    }

    return n.reply(msg, map_resp_body)
}

func (n *Node) topologyHandler(msg Message) error {
    recv_body, err := decodeMessageBody[TopologyMessageBody](msg.Body)
    if err != nil {
        return err
    }

    n.neighbors = recv_body.Topology[n.nodeId]
    n.log("My neighbors: %v\n", n.neighbors)

    resp_body := BaseMessageBody{
        Type: "topology_ok",
    }
    map_resp_body, err := encodeMessageBodyToMap(&resp_body)
    if err != nil {
        return err
    }

    return n.reply(msg, map_resp_body)
}

func (n *Node) broadcastHandler(msg Message) error {
    // 1. reply ok
    resp_body := BaseMessageBody {
        Type: "broadcast_ok",
    }
    map_resp_body, err := encodeMessageBodyToMap(&resp_body)
    if err != nil {
        return err
    }
    if err := n.reply(msg, map_resp_body); err != nil {
        return err
    }

    // 2. check the message
    recv_body, err := decodeMessageBody[BroadcastMessageBody](msg.Body)
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

    n.log("Received message: %v\n", recv_body.Message)

    // 3. gossip to neighbors
    gossip_body := BroadcastMessageBody{
        BaseMessageBody: BaseMessageBody{ Type: "broadcast" },
        Message: recv_body.Message,
    }
    map_gossip_body, err := encodeMessageBodyToMap(&gossip_body)
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
            n.rpc(dest, map_gossip_body, func(msg Message) error {
                recv_body, err := decodeMessageBody[BaseMessageBody](msg.Body)
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

func (n *Node) readHandler(msg Message) error {
    messages := make([]any, 0)

    n.messagesLock.Lock()
    for message := range n.messages {
        messages = append(messages, message)
    }
    n.messagesLock.Unlock()

    resp_body := ReadMessageBody{
        BaseMessageBody: BaseMessageBody{ Type: "read_ok" },
        Messages: messages,
    }
    map_resp_body, err := encodeMessageBodyToMap(&resp_body)
    if err != nil {
        return err
    }

    return n.reply(msg, map_resp_body)
}
