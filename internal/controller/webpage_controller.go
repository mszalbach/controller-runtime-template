// Package controller contains all controllers for this operator
package controller

import (
	"context"

	api "github.com/mszalbach/controller-runtime-template/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
)

// Permissions needed for this controller
// TODO erzeugt nur die Rolle aber nicht das Binding und die ServiceAccounts daher vermutlich nicht so hilfreich?
// +kubebuilder:rbac:groups=mszalbach.github.com,resources=*,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=mszalbach.github.com,resources=*/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mszalbach.github.com,resources=*/finalizers,verbs=update;patch
// +kubebuilder:rbac:groups=apps,resources=pods,verbs=get;list;watch;create;update;patch;delete

// WebPageReconciler reconciles a WebPage object
type WebPageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile method for WebPageReconciler
func (r *WebPageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx).WithValues("webpage", req.NamespacedName)
	log.Info("reconciling webpage")

	var webpage api.WebPage
	if err := r.Get(ctx, req.NamespacedName, &webpage); err != nil {
		log.Error(err, "unable to get webpage")
		return ctrl.Result{}, err
	}

	var pod corev1.Pod
	podFound := true
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "unable to get pod")
			return ctrl.Result{}, err
		}
		podFound = false
	}

	if podFound {
		if err := r.Delete(ctx, &pod); err != nil {
			log.Error(err, "unable to delete pod")
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}
	templ := webpage.DeepCopy()
	pod.Name = req.Name
	pod.Namespace = req.Namespace

	// TODO use cdr content for webpage
	pod.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:  req.Name,
				Image: templ.Spec.Image,
				Ports: []corev1.ContainerPort{
					{
						Name:          "http",
						Protocol:      corev1.ProtocolTCP,
						ContainerPort: 80,
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&webpage, &pod, r.Scheme); err != nil {
		log.Error(err, "unable to set pod's owner reference")
		return ctrl.Result{}, err
	}

	if err := r.Create(ctx, &pod); err != nil {
		log.Error(err, "unable to create pod")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager helps to install this controller to a manager
func (r *WebPageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.WebPage{}).
		Named("webpage").
		Complete(r)
}
