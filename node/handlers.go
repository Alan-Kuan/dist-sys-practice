package node

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
    n.messages = append(n.messages, recv_body.Message)
    n.messagesLock.Unlock()
    n.log("Received message: %v\n", recv_body.Message)

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
    copy(messages, n.messages)
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
