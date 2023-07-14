package utils

import (
	"bytes"
	"text/template"

	"github.com/kubebuilder-demo/api/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func parseTemplate(templateName string, app *v1beta1.App) []byte {
	tmpl, err := template.ParseFiles(templateName + ".yml")
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, app)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func NewDeployment(app *v1beta1.App) *appv1.Deployment {
	d := &appv1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", app), d)
	if err != nil {
		panic(err)
	}
	return d
}

func NewIngress(app *v1beta1.App) *netv1.Ingress {
	i := &netv1.Ingress{}
	err := yaml.Unmarshal(parseTemplate("ingress", app), i)
	if err != nil {
		panic(err)
	}
	return i
}

func NewService(app *v1beta1.App) *corev1.Service {
	s := &corev1.Service{}
	err := yaml.Unmarshal(parseTemplate("service", app), s)
	if err != nil {
		panic(err)
	}
	return s
}

// newDeployment creates a new Deployment for a App resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the App resource that 'owns' it.
func newDeployment(app *v1beta1.App) *appv1.Deployment {
	labels := map[string]string{
		"app": app.GetObjectMeta().GetName(),
	}
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.GetObjectMeta().GetName(),
			Namespace: app.GetObjectMeta().GetNamespace(),
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &app.Spec.Replicas,
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
							Name:  app.GetObjectMeta().GetName(),
							Image: app.Spec.Image,
						},
					},
				},
			},
		},
	}
}

func newService(app *v1beta1.App) *corev1.Service {
	labels := map[string]string{
		"app": app.GetObjectMeta().GetName(),
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.GetObjectMeta().GetName(),
			Namespace: app.GetObjectMeta().GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.IntOrString{IntVal: 80},
				},
			},
		},
	}
}

func newIngress(app *v1beta1.App) *v1.Ingress {
	pathType := v1.PathTypePrefix
	return &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.GetObjectMeta().GetName(),
			Namespace: app.GetObjectMeta().GetNamespace(),
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: app.GetObjectMeta().GetName(),
											Port: v1.ServiceBackendPort{
												Number: 8080,
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
}
