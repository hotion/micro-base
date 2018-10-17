package service

import (
	"context"
	"fmt"
)

// PostHello hello word 测试demo
func (s basicService) PostHello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("hello %s", name), nil
}
