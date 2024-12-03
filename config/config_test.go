package config

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	configs := []string{
		"example",
		"dev",
		"staging",
		"prod",
	}

	for _, c := range configs {
		t.Run(c, func(t *testing.T) {
			b, err := ioutil.ReadFile(fmt.Sprintf("../../config.yml.%s", c))
			assert.NoError(t, err)
			var out interface{}
			err = yaml.Unmarshal(b, &out)
			if err != nil {
				t.Fatalf("got error on config.yml.%s \n%v", c, err)
			}
		})
	}
}
