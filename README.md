# dist-sys-practice
My practice of writing distributed systems in Go.

I validated my systems with a benmark tool, [maelstrom](https://github.com/jepsen-io/maelstrom).
Besides, there was a guide of writing different kinds of distributed workloads provided along with the tool.
Therefore, I followed it and implemented mine step by step.

## Workloads
| Kind | Description |
| --- | --- |
| Echo | Such node replies the same message a client just sends. |
| Broadcast | Such node maintains a set of elements. Whenever it receives an "add" request from a client, it broadcasts the request to neighboring nodes. Clients can send a "read" request to any node to obtain the synchronized set of elements. |
| G-set | Such node maintains a set of elements. Contrary to previous workload, it replicates the whole set to every other node periodically. Clients can also obtain the set in the same way. |
| G-counter | Such node maintains 2 maps of values. One keeps every cumulative positive delta values of every node; while the other keeps every negative ones. Clients can send an "add" request to any node with a positive/negative delta value, and a "read" request to any node to obtain the sum of all values. |

## Quick Notes
### G-Counter

每個 node 維護兩個 vector 去紀錄每個 node 目前的 counter 值
一個 vector 負責「加」的請求，另個負責「減」的請求。

> 為什麼不直接記一個值去同步就好？
 
> 做不到。和其他節點分享時，無法知道哪個才是最新的，尤其今天的狀況允許「加」和「減」。
