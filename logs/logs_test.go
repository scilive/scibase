package logs_test

import (
	"testing"

	"github.com/scilive/scibase/logs"
)

func TestLog(t *testing.T) {
	logs.Log.Info().Msg("test")
}
