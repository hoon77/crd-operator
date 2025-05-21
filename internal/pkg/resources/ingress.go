package resources

import (
	webappv1 "github.com/hoon77/crd-operator/api/v1"
	"github.com/hoon77/crd-operator/internal/pkg/utils"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	IngressTLSSecretNameMaSuffix = "-tls"
	NginxClassName               = "nginx"
	NginxRewriteTargetAnoKey     = "nginx.ingress.kubernetes.io/rewrite-target"
)

func BuildIngress(webapp *webappv1.WebApp) *networkingv1.Ingress {
	if webapp.Spec.Ingress == nil || !webapp.Spec.Ingress.Enabled {
		return nil
	}

	path := "/"
	if webapp.Spec.Ingress.Path != "" {
		path = webapp.Spec.Ingress.Path
	}

	port := webapp.Spec.Ingress.Port
	if port == 0 {
		port = 80
	}

	className := webapp.Spec.Ingress.ClassName
	if className != "" {
		className = NginxClassName
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        webapp.Name,
			Namespace:   webapp.Namespace,
			Labels:      utils.GetCommonLabels(webapp),
			Annotations: GetRewriteTargetAnnotations(className, webapp.Spec.Ingress.RewriteTarget),
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &className,
			Rules: []networkingv1.IngressRule{
				{
					Host: webapp.Spec.Ingress.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     path,
									PathType: utils.PtrPathType(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: webapp.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: webapp.Spec.Ingress.Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if webapp.Spec.Ingress.TLS {
		ingress.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts:      []string{webapp.Spec.Ingress.Host},
				SecretName: webapp.Name + IngressTLSSecretNameMaSuffix,
			},
		}
	}

	return ingress
}

func GetRewriteTargetAnnotations(className string, rewriteTarget string) map[string]string {
	var annotations map[string]string
	if rewriteTarget == "" {
		rewriteTarget = "/"
	}
	switch className {
	case NginxClassName:
		annotations = map[string]string{
			NginxRewriteTargetAnoKey: rewriteTarget,
		}
	}

	return annotations
}
