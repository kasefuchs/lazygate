package utils

import (
	"fmt"
	"strings"
)

// TagsToLabels converts tags to labels.
func TagsToLabels(tags []string) map[string]string {
	labels := make(map[string]string, len(tags))
	for _, tag := range tags {
		if parts := strings.SplitN(tag, "=", 2); len(parts) == 2 {
			left, right := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			labels[left] = right
		}
	}

	return labels
}

func TransformLabel(key, oldPrefix, newPrefix string) string {
	if strings.HasPrefix(key, oldPrefix+".") {
		trimmed := strings.TrimPrefix(key, oldPrefix+".")
		return fmt.Sprintf("%s.%s", newPrefix, trimmed)
	}

	return key
}

func TransformLabels(labels map[string]string, oldPrefix, newPrefix string) map[string]string {
	transformedLabels := make(map[string]string)
	for key, value := range labels {
		newKey := TransformLabel(key, oldPrefix, newPrefix)
		transformedLabels[newKey] = value
	}

	return transformedLabels
}

// TransformEnvVars заменяет префикс в переменных окружения
func TransformEnvVars(envVars []string, oldPrefix, newPrefix string) []string {
	var result []string

	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		if strings.HasPrefix(key, oldPrefix) {
			trimmed := strings.TrimPrefix(key, oldPrefix)
			newKey := newPrefix + trimmed
			result = append(result, fmt.Sprintf("%s=%s", newKey, value))
		} else {
			result = append(result, env)
		}
	}

	return result
}
