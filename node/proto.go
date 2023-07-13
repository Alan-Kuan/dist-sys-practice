package node

import (
	"encoding/json"
	"sync"
)

type Handler func(Message) error

type Node struct {
    nodeId          string
    nodeIds         []string

    nextMsgId       int
    nextMsgIdLock   *sync.Mutex

    logLock         *sync.Mutex

    handlers        map[string]Handler
    callbacks       map[int]Handler
    wg              *sync.WaitGroup

    neighbors       []string
    
    messages        map[any]struct{}  // work as a set
    messagesLock    *sync.Mutex
}

type Message struct {
    Src     string          `json:"src"`
    Dest    string          `json:"dest"`
    Body    json.RawMessage `json:"body"`
}

type BaseMessageBody struct {
    MsgId       *int    `json:"msg_id,omitempty"`
    Type        string  `json:"type"`
    InReplyTo   *int    `json:"in_reply_to,omitempty"`
}

type InitMessageBody struct {
    BaseMessageBody
    NodeId      string      `json:"node_id"`
    NodeIds     []string    `json:"node_ids"`
}

type EchoMessageBody struct {
    BaseMessageBody
    Echo    string  `json:"echo"`
}

type TopologyMessageBody struct {
    BaseMessageBody
    Topology    map[string][]string `json:"topology"`
}

type BroadcastMessageBody struct {
    BaseMessageBody
    Message     any     `json:"message"`
}

type ReadMessageBody struct {
    BaseMessageBody
    Messages    []any   `json:"messages"`
}

type MessageBody interface {
    BaseMessageBody |
    InitMessageBody |
    EchoMessageBody |
    TopologyMessageBody |
    BroadcastMessageBody |
    ReadMessageBody
}

func decodeMessageBody[B MessageBody](raw_body json.RawMessage) (*B, error) {
    var body B

    if err := json.Unmarshal(raw_body, &body); err != nil {
        return nil, err
    }

    return &body, nil
}

func encodeMessageBodyToMap[B MessageBody](body *B) (*map[string]any, error) {
    var map_body map[string]any

    raw_body, err := json.Marshal(*body)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(raw_body, &map_body); err != nil {
        return nil, err
    }

    return &map_body, nil
}
