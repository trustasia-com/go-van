package confx

import "testing"

func TestExpandEnvValue(t *testing.T) {
	lookup := func(key string) (string, bool) {
		values := map[string]string{
			"SET":   "value",
			"EMPTY": "",
		}
		value, ok := values[key]
		return value, ok
	}

	tests := map[string]string{
		"${SET}":                   "value",
		"${MISSING}":               "",
		"${SET:-default}":          "value",
		"${EMPTY:-default}":        "default",
		"${MISSING:-default}":      "default",
		"${FIRST:-a}-${SECOND:-b}": "a-b",
		"${INVALID-NAME}":          "${INVALID-NAME}",
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			if got := expandEnvValue(input, lookup); got != want {
				t.Fatalf("expandEnvValue(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

func TestExpandEnv(t *testing.T) {
	t.Setenv("CONFX_HOST", "database.internal")
	t.Setenv("CONFX_PORT", "")

	input := []byte(`dsn: "postgres://${CONFX_HOST}:${CONFX_PORT:-5432}/app"`)
	want := `dsn: "postgres://database.internal:5432/app"`
	if got := string(ExpandEnv(input)); got != want {
		t.Fatalf("ExpandEnv() = %q, want %q", got, want)
	}
}
