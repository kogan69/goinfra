package config_test

import (
	"testing"

	"github.com/kogan69/goinfra/config"
	"github.com/stretchr/testify/assert"
)

func Test_Configure(t *testing.T) {
	assert.Equal(t, "", config.GetEnvValue("DEBUG_LEVEL", ""))
	err := config.InitEnvFromFile(".env.test")
	assert.NoError(t, err)
	assert.Equal(t, "info", config.GetEnvValue("DEBUG_LEVEL", ""))
}

func Test_Bad_Config_File_Path(t *testing.T) {
	err := config.InitEnvFromFile(".env.")
	assert.Error(t, err)
}
