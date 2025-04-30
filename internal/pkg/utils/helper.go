package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func SetOwnerRefence(owner metav1.Object, object metav1.Object, scheme *runtime.Scheme) error {
	return ctrl.SetControllerReference(owner, object, scheme)
}
