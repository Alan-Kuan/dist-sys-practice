package main

func newGCounter() *gCounter {
    counter := &gCounter{
        increments: make(map[string]int),
        decrements: make(map[string]int),
    }
    return counter
}

func (g_counter *gCounter) read() int {
    sum := 0
    for _, value := range g_counter.increments {
        sum += value
    }
    for _, value := range g_counter.decrements {
        sum -= value
    }
    return sum
}

func (g_counter *gCounter) add(node_id string, delta int) {
    if delta >= 0 {
        add(&g_counter.increments, node_id, delta)
    } else {
        add(&g_counter.decrements, node_id, -delta)
    }
}

func add(dest_map *map[string]int, node_id string, pos_delta int) {
    if _, exists := (*dest_map)[node_id]; exists {
        (*dest_map)[node_id] += pos_delta
    } else {
        (*dest_map)[node_id] = pos_delta
    }
}

func (g_counter *gCounter) merge(other_increments *map[string]int,
                                 other_decrements *map[string]int) {
    merge(&g_counter.increments, other_increments)
    merge(&g_counter.decrements, other_decrements)
}

func merge(dest_map *map[string]int, other_map *map[string]int) {
    for node_id, other_value := range *other_map {
        if value, exists := (*dest_map)[node_id]; exists {
            if other_value > value {
                (*dest_map)[node_id] = other_value
            }
        } else {
            (*dest_map)[node_id] = other_value
        }
    }
}
