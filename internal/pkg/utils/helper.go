package utils

import (
	webappv1 "github.com/hoon77/crd-operator/api/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func SetOwnerRefence(owner metav1.Object, object metav1.Object, scheme *runtime.Scheme) error {
	return ctrl.SetControllerReference(owner, object, scheme)
}

func GetCommonLabels(webapp *webappv1.WebApp) map[string]string {
	return map[string]string{
		"app": webapp.Name,
	}
}

func PtrPathType(p networkingv1.PathType) *networkingv1.PathType {
	return &p
}
