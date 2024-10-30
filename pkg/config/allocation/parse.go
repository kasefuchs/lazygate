package allocation

import (
	"github.com/go-playground/validator/v10"
	"github.com/kasefuchs/lazygate/pkg/utils"
	"github.com/traefik/paerser/parser"
)

// Root label name.
const rootName = "lazygate"

// Dynamic config validator.
var validate = validator.New()

// ParseLabels parses labels to dynamic config.
func ParseLabels(labels map[string]string) (*Config, error) {
	cfg := DefaultConfig()

	if err := parser.Decode(labels, cfg, rootName, rootName); err != nil {
		return nil, err
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ParseTags parses tags to dynamic config.
func ParseTags(tags []string) (*Config, error) {
	labels := utils.TagsToLabels(tags)

	return ParseLabels(labels)
}
