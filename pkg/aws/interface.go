package aws

import (
	corev1 "k8s.io/api/core/v1"
)

//nolint
//go:generate mockgen -package=mocks -destination ../mocks/mock_instance_tagger.go github.com/ouzi-dev/node-tagger/pkg/aws NodeTagger

type NodeTagger interface {
	EnsureInstanceNodeHasTags(node *corev1.Node, tags map[string]string) error
}
