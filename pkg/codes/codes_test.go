// Package codes provides ...
package codes

import (
	"testing"
)

var (
	testCode   Code = 2000
	formatCode Code = 2001
)

func init() {
	trans := MemoryTranslator{
		Code2Desc: map[string]map[Code]string{
			LangZhCN: {
				testCode:   "错误测试",
				formatCode: "测试: %s %s",
			},
		},
	}
	SetTranslator(trans)
}

func TestTrEmbeded(t *testing.T) {
	str := Unknown.Tr(LangZhCN)
	t.Log(str)
	str = Unknown.Tr(LangEnUS)
	t.Log(str)
}

func TestTrEmbededArgs(t *testing.T) {
	str := Unknown.Tr(LangZhCN, "你做什么")
	t.Log(str)
	str = Unknown.Tr(LangEnUS, "what are you doing")
	t.Log(str)
}

func TestTrEmbededNoLang(t *testing.T) {
	str := Unknown.Tr("zh-hk")
	t.Log(str)
}

func TestTrCustom(t *testing.T) {
	str := testCode.Tr(LangZhCN)
	t.Log(str)
}

func TestTrCustomFormat(t *testing.T) {
	str := formatCode.Tr(LangZhCN, "hello", "world")
	t.Log(str)
}

func TestTrCustomErrFormat(t *testing.T) {
	str := formatCode.Tr(LangZhCN, 1)
	t.Log(str)
}

func TestTrCustomNoLang(t *testing.T) {
	str := testCode.Tr(LangEnUS)
	t.Log(str)
}
