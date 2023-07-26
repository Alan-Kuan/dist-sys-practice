package main

import (
    "sync"

	"alan-kuan/dist-sys-practice/pkg/node"
)

type gSetNode struct {
    *node.Node
    periodicTasks   []periodicTask
    crdt            gSet
    crdtLock        *sync.RWMutex
}

type periodicTask struct {
    interval    int
    callback    func()
}

type gSet struct {
    set map[any]struct{}  // work as a set
}

type addMessageBody struct {
    node.BaseMessageBody
    Element any `json:"element"`
}

type readMessageBody struct {
    node.BaseMessageBody
    Value   []any   `json:"value"`
}

type replicateMessageBody readMessageBody
