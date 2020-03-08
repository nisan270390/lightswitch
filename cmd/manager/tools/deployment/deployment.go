package deployment

import (
	lightswitchv1alpha1 "LightSwitch/pkg/apis/lightswitch/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	extensionv1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateLightSwitchDeployment creates light switch deployment
func CreateLightSwitchDeployment(cr *lightswitchv1alpha1.LightSwitch) *extensionv1.Deployment {
	name := cr.Spec.ServiceName + "-light-switch"
	var terminationGracePeriod int64 = 110

	labels := map[string]string{
		"app": name,
	}

	handler := corev1.Handler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: cr.Spec.HealthcheckPath,
			Port: intstr.IntOrString{
				IntVal: cr.Spec.Port,
			},
		},
	}

	return &extensionv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},

		Spec: extensionv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: cr.Spec.PodAnnotations,
				},

				Spec: corev1.PodSpec{
					ServiceAccountName:            cr.Spec.ServiceName,
					Affinity:                      affinity(name),
					PriorityClassName:             "riskified-critical-spots",
					Tolerations:                   toleration(),
					TerminationGracePeriodSeconds: &terminationGracePeriod,
					Containers: []corev1.Container{
						{
							Name:            "light-switch",
							Image:           cr.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: cr.Spec.Port,
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler:             handler,
								InitialDelaySeconds: 120,
								PeriodSeconds:       20,
								TimeoutSeconds:      5,
								FailureThreshold:    3,
							},
						},
					},
				},
			},
		},
	}
}

func affinity(name string) *corev1.Affinity {
	podAntiAffinityTermSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "app",
				Operator: "In",
				Values: []string{
					name,
				},
			},
		},
	}

	return &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 90,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &podAntiAffinityTermSelector,
						TopologyKey:   "failure-domain.beta.kubernetes.io/zone",
					},
				},
				{
					Weight: 80,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &podAntiAffinityTermSelector,
						TopologyKey:   "beta.kubernetes.io/instance-type",
					},
				},
				{
					Weight: 70,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &podAntiAffinityTermSelector,
						TopologyKey:   "kubernetes.io/hostname",
					},
				},
			},
		},

		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				{
					Weight: 100,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "node-role.kubernetes.io/spot-worker",
								Operator: "In",
								Values: []string{
									"true",
								},
							},
						},
					},
				},
			},
		},
	}
}

func toleration() []corev1.Toleration {
	return []corev1.Toleration{
		{
			Key:      "spot-instance",
			Operator: "Equal",
			Value:    "True",
			Effect:   "NoSchedule",
		},
	}
}
