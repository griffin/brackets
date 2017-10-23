package env

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestEnv_GenerateValidator(t *testing.T) {
	e := New()
	t.Errorf(e.GenerateSelector(12))

}
