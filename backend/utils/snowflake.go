package utils

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func InitSnowflake() error {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("snowflake node init failed: %w", err)
	}
	return nil
}

func GenerateID() string {
	return node.Generate().String()
}
