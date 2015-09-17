package main

import "encoding/json"

func decode(payload []byte, value interface{}) error {
	var err error

	if err = json.Unmarshal(payload, value); err != nil {
		return err
	}

	return nil
}
