package main

import (
	"github.com/AndreiMartynenko/chat-server/cli/cmd/root"
	"github.com/AndreiMartynenko/common/pkg/closer"
)

func main() {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	root.Execute()
}
