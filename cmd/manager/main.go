/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/kubernetes-sigs/kernel-module-management/controllers/notifier"
	"github.com/kubernetes-sigs/kernel-module-management/internal/config"
	"github.com/kubernetes-sigs/kernel-module-management/internal/proto"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	clusterv1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	v1beta12 "github.com/kubernetes-sigs/kernel-module-management/api/v1beta1"
	"github.com/kubernetes-sigs/kernel-module-management/controllers"
	"github.com/kubernetes-sigs/kernel-module-management/internal/build"
	"github.com/kubernetes-sigs/kernel-module-management/internal/build/job"
	"github.com/kubernetes-sigs/kernel-module-management/internal/cmd"
	"github.com/kubernetes-sigs/kernel-module-management/internal/constants"
	"github.com/kubernetes-sigs/kernel-module-management/internal/daemonset"
	"github.com/kubernetes-sigs/kernel-module-management/internal/filter"
	"github.com/kubernetes-sigs/kernel-module-management/internal/metrics"
	"github.com/kubernetes-sigs/kernel-module-management/internal/module"
	"github.com/kubernetes-sigs/kernel-module-management/internal/preflight"
	"github.com/kubernetes-sigs/kernel-module-management/internal/registry"
	"github.com/kubernetes-sigs/kernel-module-management/internal/sign"
	signjob "github.com/kubernetes-sigs/kernel-module-management/internal/sign/job"
	"github.com/kubernetes-sigs/kernel-module-management/internal/statusupdater"
	"github.com/kubernetes-sigs/kernel-module-management/internal/utils"
	//+kubebuilder:scaffold:imports
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1beta12.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	logger := klogr.New().WithName("kmm")

	ctrl.SetLogger(logger)

	setupLogger := logger.WithName("setup")

	var configFile string

	flag.StringVar(&configFile, "config", "", "The path to the configuration file.")

	klog.InitFlags(flag.CommandLine)

	flag.Parse()

	commit, err := cmd.GitCommit()
	if err != nil {
		setupLogger.Error(err, "Could not get the git commit; using <undefined>")
		commit = "<undefined>"
	}

	operatorNamespace := cmd.GetEnvOrFatalError(constants.OperatorNamespaceEnvVar, setupLogger)

	managed, err := GetBoolEnv("KMM_MANAGED")
	if err != nil {
		setupLogger.Error(err, "could not determine if we are running as managed; disabling")
		managed = false
	}

	serverAddr := cmd.GetEnvOrFatalError("GRPC_SERVER_ADDR", logger)
	daemonImage := cmd.GetEnvOrFatalError("RELATED_IMAGES_DAEMON", logger)

	ctx := ctrl.SetupSignalHandler()

	ctx, cancel := context.WithCancel(ctx)

	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterKMMServer(
		grpcServer,
		proto.NewServerImpl(logger.WithName("grpc-server")),
	)

	go func() {
		setupLogger.Info("Starting gRPC server", "address", serverAddr)

		if err = grpcServer.Serve(lis); err != nil {
			cmd.FatalError(setupLogger, err, "GRPC server error")
		}
	}()

	setupLogger.Info("Creating manager", "git commit", commit)

	cfg, err := config.ParseFile(configFile)
	if err != nil {
		cmd.FatalError(setupLogger, err, "could not parse the configuration file", "path", configFile)
	}

	options := cfg.ManagerOptions()
	options.Scheme = scheme
	options.LeaderElectionNamespace = "default"

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), *options)
	if err != nil {
		cmd.FatalError(setupLogger, err, "unable to create manager")
	}

	client := mgr.GetClient()

	filterAPI := filter.New(client)

	metricsAPI := metrics.New()
	metricsAPI.Register()

	registryAPI := registry.NewRegistry()
	jobHelperAPI := utils.NewJobHelper(client)
	buildHelperAPI := build.NewHelper()

	buildAPI := job.NewBuildManager(
		client,
		job.NewMaker(client, buildHelperAPI, jobHelperAPI, scheme),
		jobHelperAPI,
		registryAPI,
	)

	signAPI := signjob.NewSignJobManager(
		client,
		signjob.NewSigner(client, scheme, jobHelperAPI),
		jobHelperAPI,
		registryAPI,
	)

	daemonAPI := daemonset.NewCreator(client, constants.KernelLabel, scheme)
	kernelAPI := module.NewKernelMapper(buildHelperAPI, sign.NewSignerHelper())

	mc := controllers.NewModuleReconciler(
		client,
		buildAPI,
		signAPI,
		daemonAPI,
		kernelAPI,
		metricsAPI,
		filterAPI,
		statusupdater.NewModuleStatusUpdater(client),
		operatorNamespace,
	)

	if err = mc.SetupWithManager(mgr, constants.KernelLabel); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.ModuleReconcilerName)
	}

	nodeKernelReconciler := controllers.NewNodeKernelReconciler(client, constants.KernelLabel, filterAPI)

	if err = nodeKernelReconciler.SetupWithManager(mgr); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.NodeKernelReconcilerName)
	}

	if err = controllers.NewPodNodeModuleReconciler(client, daemonAPI).SetupWithManager(mgr); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.PodNodeModuleReconcilerName)
	}

	if err = controllers.NewNodeLabelModuleVersionReconciler(client).SetupWithManager(mgr); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.NodeLabelModuleVersionReconcilerName)
	}

	preflightStatusUpdaterAPI := statusupdater.NewPreflightStatusUpdater(client)
	preflightAPI := preflight.NewPreflightAPI(client, buildAPI, signAPI, registryAPI, preflightStatusUpdaterAPI, kernelAPI)

	if err = controllers.NewPreflightValidationReconciler(client, filterAPI, metricsAPI, preflightStatusUpdaterAPI, preflightAPI).SetupWithManager(mgr); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.PreflightValidationReconcilerName)
	}

	if managed {
		setupLogger.Info("Starting as managed")

		if err = clusterv1alpha1.Install(scheme); err != nil {
			cmd.FatalError(setupLogger, err, "could not add the Cluster API to the scheme")
		}

		if err = controllers.NewNodeKernelClusterClaimReconciler(client).SetupWithManager(mgr); err != nil {
			cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.NodeKernelClusterClaimReconcilerName)
		}
	}

	if err = (&v1beta12.Module{}).SetupWebhookWithManager(mgr); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create webhook", "webhook", "Module")
	}

	dcfg := controllers.DaemonConfig{
		Daemon: cfg.Daemon,
		Image:  daemonImage,
	}

	n := notifier.New(dcfg)

	dr := controllers.NewDaemonReconciler(client, operatorNamespace, n)

	if err = dr.SetupWithManager(mgr); err != nil {
		cmd.FatalError(setupLogger, err, "unable to create controller", "name", controllers.DaemonReconciler{})
	}

	//+kubebuilder:scaffold:builder

	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		cmd.FatalError(setupLogger, err, "unable to set up health check")
	}
	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		cmd.FatalError(setupLogger, err, "unable to set up ready check")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err = watcher.Add(configFile); err != nil {
		cmd.FatalError(setupLogger, err, "Could not add the config file to the watched files")
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				setupLogger.Info("inotify event", "event", ev)
				n.Set(dcfg)
			case err = <-watcher.Errors:
				setupLogger.Error(err, "inotify error")
				cancel()
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	setupLogger.Info("starting manager")
	if err = mgr.Start(ctx); err != nil {
		cmd.FatalError(setupLogger, err, "problem running manager")
	}
}

func GetBoolEnv(s string) (bool, error) {
	envValue := os.Getenv(s)

	if envValue == "" {
		return false, nil
	}

	managed, err := strconv.ParseBool(envValue)
	if err != nil {
		return false, fmt.Errorf("%q: invalid value for %s", envValue, s)
	}

	return managed, nil
}
