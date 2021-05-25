// Package logx provides ...
package logx

import "testing"

func init() {
	std = NewLogging()
}

func TestInfo(t *testing.T) {
	Infof("test")
}
