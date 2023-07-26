package main

import (
    "alan-kuan/dist-sys-practice/pkg/node"
)

type echoMessageBody struct {
    node.BaseMessageBody
    Echo    string  `json:"echo"`
}
