package resources

import (
	webappv1 "github.com/hoon77/crd-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func BuildService(webapp *webappv1.WebApp) *corev1.Service {
	labels := map[string]string{
		"app": webapp.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Name,
			Namespace: webapp.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Port:     80,
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: 80,
					},
				},
			},
		},
	}
}
