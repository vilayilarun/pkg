// package main

// import (
// 	// other imports...
// 	"github.com/go-logr/logr"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/client-go/kubernetes/scheme"
// 	_ "k8s.io/client-go/plugin/pkg/client/auth"
// 	ctrl "sigs.k8s.io/controller-runtime"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// 	"sigs.k8s.io/controller-runtime/pkg/healthz"
// 	"sigs.k8s.io/controller-runtime/pkg/log/zap"

// 	meshcomv1alpha1 "github.com/vilayilarun/pkg/api/v1alpha1"
// 	"github.com/vilayilarun/pkg/controllers"
// )

// type MeshReconciler struct {
// 	client.Client
// 	Scheme *runtime.Scheme
// 	Log    logr.Logger
// }

// var (
// 	scheme   = runtime.NewScheme()
// 	setupLog = ctrl.Log.WithName("setup")
// )

// func init() {
// 	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

// 	utilruntime.Must(meshcomv1alpha1.AddToScheme(scheme))
// 	//+kubebuilder:scaffold:scheme
// }

// func main() {
// 	var metricsAddr string
// 	var enableLeaderElection bool
// 	var probeAddr string
// 	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
// 	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
// 	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
// 		"Enable leader election for controller manager. "+
// 			"Enabling this will ensure there is only one active controller manager.")
// 	opts := zap.Options{
// 		Development: true,
// 	}
// 	opts.BindFlags(flag.CommandLine)
// 	flag.Parse()

// 	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
// 	// mgrOptions := manager.Options{
// 	// 	MetricsBindAddress: "0.0.0.0:8080",
// 	// 	// Add other options here if needed
// 	// }

// 	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
// 		Scheme: scheme,
// 		// MetricsBindAddress:     metricsAddr,
// 		// Port:                   9443,
// 		HealthProbeBindAddress: probeAddr,
// 		LeaderElection:         enableLeaderElection,
// 		LeaderElectionID:       "cdbd3f8f.mesh.com",
// 	})
// 	if err != nil {
// 		setupLog.Error(err, "unable to start manager")
// 		os.Exit(1)
// 	}

// 	if err = (&controllers.MeshReconciler{
// 		Client: mgr.GetClient(),
// 		Log:    logr.Log.WithName("controllers").WithName("MyOperator"),
// 		Scheme: mgr.GetScheme(),
// 	}).SetupWithManager(mgr); err != nil {
// 		setupLog.Error(err, "unable to create controller", "controller", "Mesh")
// 		os.Exit(1)
// 	}
// 	//+kubebuilder:scaffold:builder

// 	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
// 		setupLog.Error(err, "unable to set up health check")
// 		os.Exit(1)
// 	}
// 	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
// 		setupLog.Error(err, "unable to set up ready check")
// 		os.Exit(1)
// 	}

//		setupLog.Info("starting manager")
//		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
//			setupLog.Error(err, "problem running manager")
//			os.Exit(1)
//		}
//	}
package main

import (
	"flag"
	logr "github.com/go-logr/logr"
	"os"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	meshcomv1alpha1 "github.com/vilayilarun/pkg/api/v1alpha1"
	"github.com/vilayilarun/pkg/controllers"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	//+kubebuilder:scaffold:imports
)

type MeshReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

var (
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme.Scheme))

	utilruntime.Must(meshcomv1alpha1.AddToScheme(scheme.Scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	// mgrOptions := manager.Options{
	// 	MetricsBindAddress: "0.0.0.0:8080",
	// 	// Add other options here if needed
	// }

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme.Scheme,
		// MetricsBindAddress:     metricsAddr,
		// Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "cdbd3f8f.mesh.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.MeshReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("MyOperator"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Mesh")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
