package resources

import (
	webappv1 "github.com/hoon77/crd-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConfigMapSuffix = "-configmap"
)

func BuildConfigMap(webapp *webappv1.WebApp) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Name + ConfigMapSuffix,
			Namespace: webapp.Namespace,
		},
		Data: webapp.Spec.ConfigData,
	}
}
