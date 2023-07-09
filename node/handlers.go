package node

func (n *Node) initHandler(msg Message) error {
    n.nodeId = msg.Body.NodeId
    n.nodeIds = msg.Body.NodeIds

    return n.reply(msg, MessageBody{
        Type: "init_ok",
    })
}

func (n *Node) echoHandler(msg Message) error {
    return n.reply(msg, MessageBody{
        Type: "echo_ok",
        Echo: msg.Body.Echo,
    })
}
