package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	t.Run("config test", func(t *testing.T) {
		b, err := os.ReadFile("../config.yml.example")
		assert.NoError(t, err)
		var out interface{}
		err = yaml.Unmarshal(b, &out)
		if err != nil {
			t.Fatalf("got error on config.yml.example \n%v", err)
		}
	})
}
