package node

import "encoding/json"

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

    gossip_body := BroadcastMessageBody{
        BaseMessageBody: BaseMessageBody{ Type: "broadcast" },
        Message: recv_body.Message,
    }
    raw_gossip_body, err := json.Marshal(&gossip_body)
    if err != nil {
        return err
    }

    for _, neighbor := range n.neighbors {
        if neighbor == msg.Src {
            continue
        }
        n.send(neighbor, raw_gossip_body)
    }

    // don't reply to messages from other server nodes
    if recv_body.MsgId == nil {
        return nil
    }

    resp_body := BaseMessageBody {
        Type: "broadcast_ok",
    }
    map_resp_body, err := encodeMessageBodyToMap(&resp_body)
    if err != nil {
        return err
    }

    return n.reply(msg, map_resp_body)
}

func (n *Node) readHandler(msg Message) error {
    messages := make([]any, len(n.messages))

    n.messagesLock.Lock()
    for message := range n.messages {
        if message == nil {
            continue
        }
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
