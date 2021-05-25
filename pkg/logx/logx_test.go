// Package logx provides ...
package logx

import "testing"

func init() {
	std = NewLogging(WithShortFile())
}

func TestInfo(t *testing.T) {
	std.Info("test")
}
