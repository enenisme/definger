package utils

import (
	"testing"

	"github.com/enenisme/definger/logger"
	"github.com/stretchr/testify/assert"
)

func Test_json2json(t *testing.T) {
	oldJsonPath := "../test/rule.json"
	newJsonPath := "../test/new.json"
	logger := logger.NewLogger(logger.LogLevelDebug)
	jsonData := Json2Json(oldJsonPath, newJsonPath, logger)
	assert.NotNil(t, jsonData)
}
