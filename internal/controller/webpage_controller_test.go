package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/mszalbach/controller-runtime-template/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// https://onsi.github.io/ginkgo/#writing-specs
// TODO not sure how this works with more tests and more controllers to not interfere with each other.
var _ = Describe("Application Controller", func() {
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
			Expect(k8sClient.Create(ctx, namespace)).To(Succeed())

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
			Expect(k8sClient.Create(ctx, app)).To(Succeed())

			// Then
			pod := &corev1.Pod{}

			Eventually(func() error {
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      "test-app",
					Namespace: "test-namespace",
				}, pod)
			}, timeout, interval).Should(Succeed())

			// Verify pod
			Expect(pod.Spec.Containers).To(HaveLen(1))
			Expect(pod.Spec.Containers[0].Image).To(Equal("nginx:latest"))
			Expect(pod.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(80)))

			// Verify owner reference
			// There is no garbage collector in env test so the pod can not be deleted when the app is deleted. So this reference check should be enough
			Expect(pod.OwnerReferences).To(HaveLen(1))
			Expect(pod.OwnerReferences[0].Name).To(Equal("test-app"))
			Expect(pod.OwnerReferences[0].UID).To(Equal(app.UID))
		})
	})
})
