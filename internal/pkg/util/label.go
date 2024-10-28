package util

import (
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
