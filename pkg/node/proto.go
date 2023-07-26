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

type initMessageBody struct {
    BaseMessageBody
    NodeId      string      `json:"node_id"`
    NodeIds     []string    `json:"node_ids"`
}
