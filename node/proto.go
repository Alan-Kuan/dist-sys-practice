package node

import (
    "sync"
)

type Handler func(Message) error

type Node struct {
    nodeId          string
    nodeIds         []string
    nextMsgId       int
    nextMsgIdLock   *sync.Mutex
    handlers        map[string]Handler
    wg              *sync.WaitGroup
}

type Message struct {
    Src     string      `json:"src"`
    Dest    string      `json:"dest"`
    Body    MessageBody `json:"body"`
}

type MessageBody struct {
    MsgId       int                 `json:"msg_id"`
    Type        string              `json:"type"`
    InReplyTo   *int                `json:"in_reply_to,omitempty"`

    // init
    NodeId      string              `json:"node_id,omitempty"`
    NodeIds     []string            `json:"node_ids,omitempty"`

    // echo
    Echo        string              `json:"echo,omitempty"`

    // topology
    Topology    map[string][]string `json:"topology,omitempty"`

    // broadcast
    Message     interface {}        `json:"message,omitempty"`

    // read
    Messages    []interface {}      `json:"messages,omitempty"`
}
