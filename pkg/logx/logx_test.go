// Package logx provides ...
package logx

import (
	"context"
	"fmt"
	"testing"
)

func init() {
	std = NewLogging(WithService("test"))
}

func TestLogging(t *testing.T) {
	SetLevel(LevelDebug)
	Debug("test")
	Debugf("test: %s", "hello")

	Info("test")
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

	WithData(map[string]any{
		"hello": "world",
	}).Infof("hahaha: %s", "test")
	WithContext(context.Background()).Info("test")
}

func TestNewLogging(t *testing.T) {
	logger := NewLogging(
		WithLevel(LevelError),
		WithService("test-service"),
		WithFlag(FlagFile),
	)
	logger.Debug("test")
	logger.Debugf("test: %s", "hello")

	logger.Info("test")
	logger.Infof("test: %s", "hello")

	logger.Warning("test")
	logger.Warningf("test: %s", "hello")

	logger.Error("test")
	logger.Errorf("test: %s", "hello")
}

func BenchmarkLogging(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info(fmt.Sprint(i))
	}
}
