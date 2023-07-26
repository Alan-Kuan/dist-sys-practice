package main

import (
    "alan-kuan/dist-sys-practice/pkg/utils"
)

func newGSet() *gSet {
    set := &gSet{
        set: make(map[any]struct{}),
    }
    return set
}

func (g_set *gSet) read() []any {
    return utils.MapToSlice(&g_set.set)
}

func (g_set *gSet) add(element any) {
    g_set.set[element] = struct{}{}
}

func (g_set *gSet) merge(other_slice []any) {
    for _, element := range other_slice {
        g_set.set[element] = struct{}{}
    }
}
