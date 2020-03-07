package node

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/ouzi-dev/node-tagger/pkg/aws"
	"github.com/ouzi-dev/node-tagger/pkg/flags"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_node")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Node Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	reconciler, err := newReconciler(mgr)
	if err != nil {
		return err
	}

	return add(mgr, reconciler)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	awsSession, err := aws.GetAwsSessionFromEnv()
	if err != nil {
		return nil, err
	}

	return &ReconcileNode{
		client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		nodeTagger: aws.NewNodeInstanceTagger(ec2.New(awsSession)),
	}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("node-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Node
	err = c.Watch(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNode implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNode{}

// ReconcileNode reconciles a Node object
type ReconcileNode struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	nodeTagger aws.NodeTagger
}

// Reconcile reads that state of the cluster for a Node object and adds tags to the underlying instances if necessary
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNode) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Node.Name", request.Name)
	reqLogger.Info("Reconciling Node")

	// Fetch the Node instance
	instance := &corev1.Node{}

	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Node is not an aws one so skip it. Return and don't requeue
	if !isAwsNode(instance) {
		reqLogger.Info("Node is not an AWS node. Skipping")
		return reconcile.Result{}, nil
	}

	err = r.nodeTagger.EnsureInstanceNodeHasTags(instance, flags.InstanceTags)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Node already tagged properly
	reqLogger.Info("Node instance has been tagged")

	return reconcile.Result{}, nil
}

func isAwsNode(node *corev1.Node) bool {
	return strings.HasPrefix(node.Spec.ProviderID, "aws")
}
