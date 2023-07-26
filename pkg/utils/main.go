package utils

import (
    "encoding/json"
)

func DecodeMessageBody[B any](raw_body json.RawMessage) (*B, error) {
    var body B

    if err := json.Unmarshal(raw_body, &body); err != nil {
        return nil, err
    }

    return &body, nil
}

func EncodeMessageBodyToMap[B any](body *B) (*map[string]any, error) {
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
