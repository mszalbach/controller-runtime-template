package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/mszalbach/controller-runtime-template/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// https://onsi.github.io/ginkgo/#writing-specs
// TODO controller manuell testen und reconcile loop aufrufen, für mehr controlle
var _ = Describe("Application Controller", func() {
	t := GinkgoT()
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating an Application", func() {
		It("Should create a Deployment", func() {
			ctx := context.Background()

			// Given
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-namespace",
				},
			}
			require.NoError(t, k8sClient.Create(ctx, namespace))

			// When
			app := &api.WebPage{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-app",
					Namespace: "test-namespace",
				},
				Spec: api.WebPageSpec{
					Image: "nginx:latest",
				},
			}
			require.NoError(t, k8sClient.Create(ctx, app))

			// Then
			pod := &corev1.Pod{}
			assert.EventuallyWithT(t, func(c *assert.CollectT) {
				err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      "test-app",
					Namespace: "test-namespace",
				}, pod)
				assert.NoError(c, err)
			}, timeout, interval)

			// Verify pod
			assert.Len(t, pod.Spec.Containers, 1)
			assert.Equal(t, "nginx:latest", pod.Spec.Containers[0].Image)
			assert.Equal(t, int32(80), pod.Spec.Containers[0].Ports[0].ContainerPort)

			// Verify owner reference
			// There is no garbage collector in env test so the pod can not be deleted when the app is deleted. So this reference check should be enough
			assert.Len(t, pod.OwnerReferences, 1)
			assert.Equal(t, "test-app", pod.OwnerReferences[0].Name)
			assert.Equal(t, app.UID, pod.OwnerReferences[0].UID)
		})
	})
})
