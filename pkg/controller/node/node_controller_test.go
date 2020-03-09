package node

import (
	"errors"
	"testing"

	"github.com/ouzi-dev/node-tagger/pkg/flags"

	"github.com/golang/mock/gomock"
	"github.com/ouzi-dev/node-tagger/pkg/mocks"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	NodeKind       = "Node"
	NodeAPIVersion = "v1"
	name           = "Node"
)

type testReconcileItem struct {
	testName          string
	resource          *corev1.Node
	expectedError     error
	shouldTagInstance bool
}

var inputTags = map[string]string{
	"tag1": "value1",
	"tag2": "value2",
}

var tests = []testReconcileItem{
	{
		testName: "non-aws node",
		resource: &corev1.Node{
			TypeMeta: metav1.TypeMeta{
				Kind:       NodeKind,
				APIVersion: NodeAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: corev1.NodeSpec{
				ProviderID: "gce://project/region/gke-cluster",
			},
		},
		shouldTagInstance: false,
	},
	{
		testName: "aws node error tagging",
		resource: &corev1.Node{
			TypeMeta: metav1.TypeMeta{
				Kind:       NodeKind,
				APIVersion: NodeAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: corev1.NodeSpec{
				ProviderID: "aws:///az/i-instance-id",
			},
		},
		expectedError:     errors.New("error"),
		shouldTagInstance: true,
	},
	{
		testName: "aws node success tagging",
		resource: &corev1.Node{
			TypeMeta: metav1.TypeMeta{
				Kind:       NodeKind,
				APIVersion: NodeAPIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: corev1.NodeSpec{
				ProviderID: "aws:///az/i-instance-id",
			},
		},
		shouldTagInstance: true,
	},
}

func TestReconcileNode_Reconcile(t *testing.T) {
	for _, testData := range tests {
		// pin testData var in this scope
		testData := testData
		t.Run(testData.testName, func(t *testing.T) {
			flags.InstanceTags = inputTags
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockNodeTagger := mocks.NewMockNodeTagger(ctrl)

			// Register operator types with the runtime scheme.
			s := scheme.Scheme

			// Objects to track in the fake client.
			objs := []runtime.Object{
				testData.resource,
			}

			// Create a fake client to mock API calls.
			cl := fake.NewFakeClientWithScheme(s, objs...)

			// Create a ReconcileNode object with the scheme and fake client.
			r := &ReconcileNode{
				client:     cl,
				scheme:     s,
				nodeTagger: mockNodeTagger,
			}

			// Mock request to simulate Reconcile() being called on an event for a
			// watched resource .
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      name,
					Namespace: "",
				},
			}

			numberOfTimesToTagInstance := 0
			if testData.shouldTagInstance {
				numberOfTimesToTagInstance = 1
			}

			mockNodeTagger.
				EXPECT().EnsureInstanceNodeHasTags(testData.resource, inputTags).
				Return(testData.expectedError).
				Times(numberOfTimesToTagInstance)

			_, err := r.Reconcile(req)

			assert.Equal(t, testData.expectedError, err)
		})
	}
}
