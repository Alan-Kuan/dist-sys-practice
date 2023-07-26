package main

import (
    "sync"

    "alan-kuan/dist-sys-practice/pkg/node"
)

type broadcastNode struct {
    *node.Node

    neighbors       []string
    
    messages        map[any]struct{}  // work as a set
    messagesLock    *sync.Mutex
}

type topologyMessageBody struct {
    node.BaseMessageBody
    Topology    map[string][]string `json:"topology"`
}

type broadcastMessageBody struct {
    node.BaseMessageBody
    Message     any     `json:"message"`
}

type readMessageBody struct {
    node.BaseMessageBody
    Messages    []any   `json:"messages"`
}
