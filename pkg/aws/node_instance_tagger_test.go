package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/ouzi-dev/node-tagger/pkg/mocks"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	nodeName              = "ip-some-private-ip.region.compute.internal"
	noInstancesFoundError = "No instances found for the node with private dns: " +
		"ip-some-private-ip.region.compute.internal"
	multipleInstancesFoundError = "More than one instances found with private dns: " +
		"ip-some-private-ip.region.compute.internal. Cannot proceed with tagging"
	instanceID = "i-some-id"
	Once       = 1
)

var inputNode = &corev1.Node{
	ObjectMeta: metav1.ObjectMeta{
		Name: nodeName,
	},
}

var inputTags = map[string]string{
	"tag1": "value1",
	"tag2": "value2",
}

var errGeneric = errors.New("error")

type createTagsInputMatcher struct {
	x *ec2.CreateTagsInput
}

func CreateTagsInputMatcher(input *ec2.CreateTagsInput) gomock.Matcher {
	return &createTagsInputMatcher{input}
}

// Matches ensures that the input to the method matches the actual call
// Compares the resource and tags if the input
func (e createTagsInputMatcher) Matches(x interface{}) bool {
	resourceMatches := *x.(*ec2.CreateTagsInput).Resources[0] == *e.x.Resources[0]
	tagsMatch := true
	tags := x.(*ec2.CreateTagsInput).Tags
	tagsToMatch := e.x.Tags

	for _, tag := range tags {
		tagMatches := false

		for _, tagToMatch := range tagsToMatch {
			if *tagToMatch.Key == *tag.Key && *tagToMatch.Value == *tag.Value {
				tagMatches = true
			}
		}

		if !tagMatches {
			tagsMatch = false
			break
		}
	}

	return resourceMatches && tagsMatch
}

func (e createTagsInputMatcher) String() string {
	return fmt.Sprintf("input does not match %v", e.x)
}

func TestEnsureInstanceNodeHasTags_ReturnsError_If_DescribeInstances_Returns_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2Client := mocks.NewMockEC2API(ctrl)

	subject := NewNodeInstanceTagger(mockEc2Client)

	expectedDescribeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(inputNode.Name)},
			},
		},
	}

	mockEc2Client.
		EXPECT().
		DescribeInstances(&expectedDescribeInstancesInput).
		Return(nil, errGeneric).
		Times(Once)

	err := subject.EnsureInstanceNodeHasTags(inputNode, inputTags)

	assert.EqualError(t, err, errGeneric.Error())
}

func TestEnsureInstanceNodeHasTags_ReturnsError_If_DescribeInstances_Returns_No_Matches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2Client := mocks.NewMockEC2API(ctrl)

	subject := NewNodeInstanceTagger(mockEc2Client)

	expectedDescribeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(inputNode.Name)},
			},
		},
	}

	describeInstancesOutput := ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{},
			},
		},
	}

	mockEc2Client.
		EXPECT().
		DescribeInstances(&expectedDescribeInstancesInput).
		Return(&describeInstancesOutput, nil).
		Times(Once)

	err := subject.EnsureInstanceNodeHasTags(inputNode, inputTags)

	assert.EqualError(t, err, noInstancesFoundError)
}

func TestEnsureInstanceNodeHasTags_ReturnsError_If_DescribeInstances_Returns_Multiple_Matches(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2Client := mocks.NewMockEC2API(ctrl)

	subject := NewNodeInstanceTagger(mockEc2Client)

	expectedDescribeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(inputNode.Name)},
			},
		},
	}

	describeInstancesOutput := ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:     aws.String("i-one"),
						PrivateDnsName: aws.String(nodeName),
					},
					{
						InstanceId:     aws.String("i-two"),
						PrivateDnsName: aws.String(nodeName),
					},
				},
			},
		},
	}

	mockEc2Client.
		EXPECT().
		DescribeInstances(&expectedDescribeInstancesInput).
		Return(&describeInstancesOutput, nil).
		Times(Once)

	err := subject.EnsureInstanceNodeHasTags(inputNode, inputTags)

	assert.EqualError(t, err, multipleInstancesFoundError)
}

func TestEnsureInstanceNodeHasTags_ReturnsNoError_If_InstanceAlreadyTagged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2Client := mocks.NewMockEC2API(ctrl)

	subject := NewNodeInstanceTagger(mockEc2Client)

	expectedDescribeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(inputNode.Name)},
			},
		},
	}

	describeInstancesOutput := ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:     aws.String(instanceID),
						PrivateDnsName: aws.String(nodeName),
						Tags: []*ec2.Tag{
							{
								Key:   aws.String("tag1"),
								Value: aws.String("value1"),
							},
							{
								Key:   aws.String("tag2"),
								Value: aws.String("value2"),
							},
						},
					},
				},
			},
		},
	}

	mockEc2Client.
		EXPECT().
		DescribeInstances(&expectedDescribeInstancesInput).
		Return(&describeInstancesOutput, nil).
		Times(Once)

	err := subject.EnsureInstanceNodeHasTags(inputNode, inputTags)

	assert.NoError(t, err, multipleInstancesFoundError)
}

func TestEnsureInstanceNodeHasTags_ReturnsError_If_InstanceFailsTagging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2Client := mocks.NewMockEC2API(ctrl)

	subject := NewNodeInstanceTagger(mockEc2Client)

	expectedDescribeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(inputNode.Name)},
			},
		},
	}

	describeInstancesOutput := ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:     aws.String(instanceID),
						PrivateDnsName: aws.String(nodeName),
						Tags:           []*ec2.Tag{},
					},
				},
			},
		},
	}

	expectedCreateTagsInput := ec2.CreateTagsInput{
		Resources: []*string{aws.String(instanceID)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("tag1"),
				Value: aws.String("value1"),
			},
			{
				Key:   aws.String("tag2"),
				Value: aws.String("value2"),
			},
		},
	}

	mockEc2Client.
		EXPECT().
		DescribeInstances(&expectedDescribeInstancesInput).
		Return(&describeInstancesOutput, nil).
		Times(Once)

	mockEc2Client.
		EXPECT().
		CreateTags(CreateTagsInputMatcher(&expectedCreateTagsInput)).
		Return(nil, errGeneric).
		Times(Once)

	err := subject.EnsureInstanceNodeHasTags(inputNode, inputTags)

	assert.EqualError(t, err, errGeneric.Error())
}

func TestEnsureInstanceNodeHasTags_ReturnsNoError_If_InstanceSucceedsTagging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEc2Client := mocks.NewMockEC2API(ctrl)

	subject := NewNodeInstanceTagger(mockEc2Client)

	expectedDescribeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(inputNode.Name)},
			},
		},
	}

	describeInstancesOutput := ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:     aws.String(instanceID),
						PrivateDnsName: aws.String(nodeName),
						Tags:           []*ec2.Tag{},
					},
				},
			},
		},
	}

	expectedCreateTagsInput := ec2.CreateTagsInput{
		Resources: []*string{aws.String(instanceID)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("tag1"),
				Value: aws.String("value1"),
			},
			{
				Key:   aws.String("tag2"),
				Value: aws.String("value2"),
			},
		},
	}

	mockEc2Client.
		EXPECT().
		DescribeInstances(&expectedDescribeInstancesInput).
		Return(&describeInstancesOutput, nil).
		Times(Once)

	mockEc2Client.
		EXPECT().
		CreateTags(CreateTagsInputMatcher(&expectedCreateTagsInput)).
		Return(nil, nil).
		Times(Once)

	err := subject.EnsureInstanceNodeHasTags(inputNode, inputTags)

	assert.NoError(t, err)
}
