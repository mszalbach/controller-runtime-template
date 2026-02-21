package controller

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/mszalbach/controller-runtime-template/api/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

func TestControllers(t *testing.T) {
	if testing.Short() {
		t.Skip("-short was passed, skipping Controllers")
	}
	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	t := GinkgoT()

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx, cancel = context.WithCancel(context.TODO())

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	err = api.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	require.NoError(t, err)
	assert.NotNil(t, k8sClient)

	// Start controller manager
	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	require.NoError(t, err)

	// TODO move this to the test suite of the controller
	// TODO install it or call it directly via reconcile, so it is possible to test the states between
	err = (&WebPageReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
	}).SetupWithManager(k8sManager)
	require.NoError(t, err)

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		assert.NoError(t, err, "failed to run manager")
	}()
})

var _ = AfterSuite(func() {
	t := GinkgoT()
	cancel()
	By("tearing down the test environment")
	err := testEnv.Stop()
	require.NoError(t, err)
})
