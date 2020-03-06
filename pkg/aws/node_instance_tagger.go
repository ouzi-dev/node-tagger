package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

type nodeInstanceTagger struct {
	ec2Client ec2iface.EC2API
}

func NewNodeInstanceTagger(ec2Client ec2iface.EC2API) NodeTagger {
	return &nodeInstanceTagger{
		ec2Client: ec2Client,
	}
}

func (n *nodeInstanceTagger) EnsureInstanceNodeHasTags(node *v1.Node, tags map[string]string) error {
	describeInstancesInput := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(node.Name)},
			},
		},
	}

	describeInstancesOutput, err := n.ec2Client.DescribeInstances(describeInstancesInput)
	if err != nil {
		return err
	}

	if len(describeInstancesOutput.Reservations) == 0 || len(describeInstancesOutput.Reservations[0].Instances) == 0 {
		return errors.Errorf("No instances found fot the node")
	}

	if len(describeInstancesOutput.Reservations) > 1 || len(describeInstancesOutput.Reservations[0].Instances) > 1 {
		return errors.Errorf("More than one instances found. Cannot proceed with tagging")
	}

	existingTags := describeInstancesOutput.Reservations[0].Instances[0].Tags

	if instanceAlreadyTagged(tags, existingTags) {
		return nil
	}

	newTags := convertDesiredTagsToAwsTags(tags)
	instanceID := describeInstancesOutput.Reservations[0].Instances[0].InstanceId

	createTagsInput := &ec2.CreateTagsInput{
		Resources: []*string{instanceID},
		Tags:      newTags,
	}

	_, err = n.ec2Client.CreateTags(createTagsInput)
	if err != nil {
		return err
	}

	return nil
}

func instanceAlreadyTagged(requestedTags map[string]string, existingTags []*ec2.Tag) bool {
	for requestedTagKey, requestedTagValue := range requestedTags {
		existingKeyFound := false

		for _, existingTag := range existingTags {
			if *existingTag.Key == requestedTagKey && *existingTag.Value == requestedTagValue {
				existingKeyFound = true
			}
		}

		if !existingKeyFound {
			return false
		}
	}

	return true
}

func convertDesiredTagsToAwsTags(requestedTags map[string]string) []*ec2.Tag {
	resultTags := []*ec2.Tag{}

	for tagKey, tagValue := range requestedTags {
		targetTag := &ec2.Tag{
			Key:   aws.String(tagKey),
			Value: aws.String(tagValue),
		}
		resultTags = append(resultTags, targetTag)
	}

	return resultTags
}
