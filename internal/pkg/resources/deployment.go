package resources

import (
	webappv1 "github.com/hoon77/crd-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func BuildDeployment(webapp *webappv1.WebApp) *appsv1.Deployment {
	labels := map[string]string{
		"app": webapp.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Name,
			Namespace: webapp.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: webapp.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webapp",
							Image: webapp.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}
