package controllers

import (
	"context"

	"github.com/kubernetes-sigs/kernel-module-management/controllers/notifier"
	"github.com/kubernetes-sigs/kernel-module-management/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const DaemonSetReconcilerName = "DaemonSet"

type DaemonConfig struct {
	config.Daemon

	Image string
}

type DaemonReconciler struct {
	client    client.Client
	namespace string
	n         *notifier.Notifier[DaemonConfig]
	scheme    *runtime.Scheme
}

func NewDaemonReconciler(client client.Client, namespace string, n *notifier.Notifier[DaemonConfig], scheme *runtime.Scheme) *DaemonReconciler {
	return &DaemonReconciler{
		client:    client,
		namespace: namespace,
		n:         n,
	}
}

func (r *DaemonReconciler) Reconcile(ctx context.Context, _ ctrl.Request) (ctrl.Result, error) {
	logger := ctrl.LoggerFrom(ctx)

	logger.Info("Getting config here")
	logger.Info("Reconciling DaemonSet here")

	const dsName = "kmm-daemon"

	ds := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dsName,
			Namespace: r.namespace,
		},
	}

	labels := map[string]string{
		"app.kubernetes.io/component":  "daemon",
		"app.kubernetes.io/part-of":    "kernel-module-management",
		"app.kubernetes.io/managed-by": "kmm-control-plane",
	}

	opRes, err := controllerutil.CreateOrPatch(ctx, r.client, &ds, func() error {
		cfg := r.n.Get()

		ds.Spec = appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: v1.PodSpec{
					NodeSelector: cfg.NodeSelector,
					Tolerations:  cfg.Tolerations,
					Containers: []v1.Container{
						{
							Name:  "daemon",
							Image: cfg.Image,
							Env: []v1.EnvVar{
								{
									Name: "GRPC_SERVER_ADDR",
									Value:
								},
								{
									Name: "NODENAME",
									ValueFrom: &v1.EnvVarSource{
										FieldRef: &v1.ObjectFieldSelector{FieldPath: "spec.nodeName"},
									},
								},
							},
							SecurityContext: &v1.SecurityContext{
								Privileged: pointer.Bool(true),
							},
						},
					},

					//Volumes:                       nil,
					//InitContainers:                nil,
					//EphemeralContainers:           nil,
					//RestartPolicy:                 "",
					//TerminationGracePeriodSeconds: nil,
					//ActiveDeadlineSeconds:         nil,
					//DNSPolicy:                     "",
					//ServiceAccountName:            "",
					//DeprecatedServiceAccount:      "",
					//AutomountServiceAccountToken:  nil,
					//NodeName:                      "",
					//HostNetwork:                   false,
					//HostPID:                       false,
					//HostIPC:                       false,
					//ShareProcessNamespace:         nil,
					//SecurityContext:               nil,
					//ImagePullSecrets:              nil,
					//Hostname:                      "",
					//Subdomain:                     "",
					//Affinity:                      nil,
					//SchedulerName:                 "",
					//HostAliases:                   nil,
					//PriorityClassName:             "",
					//Priority:                      nil,
					//DNSConfig:                     nil,
					//ReadinessGates:                nil,
					//RuntimeClassName:              nil,
					//EnableServiceLinks:            nil,
					//PreemptionPolicy:              nil,
					//Overhead:                      nil,
					//TopologySpreadConstraints:     nil,
					//SetHostnameAsFQDN:             nil,
					//OS:                            nil,
					//HostUsers:                     nil,
					//SchedulingGates:               nil,
					//ResourceClaims:                nil,
				},
			},
		}

		owner := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       appsv1.DeploymentSpec{},
			Status:     appsv1.DeploymentStatus{},
		}

		return controllerutil.SetControllerReference(owner, &ds, r.scheme)
	})

	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info(
		"DaemonSet reconciled",
		"namespace", r.namespace,
		"name", dsName,
		"result", opRes,
	)

	return ctrl.Result{}, nil
}

func (r *DaemonReconciler) SetupWithManager(mgr manager.Manager) error {
	reqs := make([]reconcile.Request, 1)
	h := handler.EnqueueRequestsFromMapFunc(func(_ context.Context, _ client.Object) []reconcile.Request {
		return reqs
	})

	ch := r.n.EventChannel()

	return ctrl.NewControllerManagedBy(mgr).
		WatchesRawSource(&source.Channel{Source: ch}, h).
		Watches(&appsv1.DaemonSet{}, h).
		Named(DaemonSetReconcilerName).
		Complete(r)
}
