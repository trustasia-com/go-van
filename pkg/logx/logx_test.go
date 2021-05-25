// Package logx provides ...
package logx

import (
	"context"
	"testing"
)

func init() {
	std = NewLogging(WithService("test"))
}

func TestLogging(t *testing.T) {
	Infof("test")
	Infof("test: %s", "hello")

	Warning("test")
	Warningf("test: %s", "hello")

	Error("test")
	Errorf("test: %s", "hello")

	// Fatal("test")
	// Fatalf("test: %s", "hello")
}

func TestEntry(t *testing.T) {
	e := NewEntry(NewLogging())
	e.Info("test")

	WithData(map[string]interface{}{
		"hello": "world",
	}).Infof("hahaha: %s", "test")
	WithContext(context.Background()).Info("test")
}
