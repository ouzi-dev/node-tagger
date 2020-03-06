package aws

import (
	corev1 "k8s.io/api/core/v1"
)

type NodeTagger interface {
	EnsureInstanceNodeHasTags(node *corev1.Node, tags map[string]string) error
}
