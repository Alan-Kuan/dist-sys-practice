package main

import (
	"encoding/json"
	"sync"

	"alan-kuan/dist-sys-practice/pkg/node"
)

type gCounterNode struct {
    *node.Node
    periodicTasks   []periodicTask
    crdt            *gCounter
    crdtLock        *sync.RWMutex
}

type periodicTask struct {
    interval    int
    callback    func()
}

type gCounter struct {
    increments  map[string]int
    decrements  map[string]int
}

type addMessageBody struct {
    node.BaseMessageBody
    Delta   int `json:"delta"`
}

type readMessageBody struct {
    node.BaseMessageBody
    Value   int `json:"value"`
}

type replicateMessageBody struct {
    node.BaseMessageBody
    Increments  json.RawMessage `json:"increments"`
    Decrements  json.RawMessage `json:"decrements"`
}
