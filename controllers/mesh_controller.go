

package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	meshcomv1alpha1 "github.com/vilayilarun/pkg/api/v1alpha1"
	v1alpha1 "github.com/vilayilarun/pkg/api/v1alpha1"
)

const (
	frontendName = "frontend"
	backendName  = "backend"
	appName      = "app"
)

// MeshReconciler reconciles a Mesh object
type MeshReconciler struct {
	client client.Client
	Scheme *runtime.Scheme
}

func (r *MeshReconciler) createConfigMap(ctx context.Context, instance *v1alpha1.Mesh, name, label string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": label,
			},
		},
		Data: data,
	}

	if err := ctrl.SetControllerReference(instance, configMap, r.client.Scheme()); err != nil {
		return err
	}

	foundConfigMap := &corev1.ConfigMap{}
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
		err = r.client.Create(ctx, configMap)
		if err != nil {
			r.Log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
			return err
		}
	} else if err != nil {
		r.Log.Error(err, "Failed to get ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
		return err
	}

	return nil
}
func (r *MeshReconciler) createSecret(ctx context.Context, instance *v1alpha1.Mesh, name, label string, data map[string]string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": label,
			},
		},
		StringData: data,
	}

	if err := ctrl.SetControllerReference(instance, secret, r.client.Scheme()); err != nil {
		return err
	}

	foundSecret := &corev1.Secret{}
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: instance.Namespace}, foundSecret)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		err = r.client.Create(ctx, secret)
		if err != nil {
			r.Log.Error(err, "Failed to create new Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
			return err
		}
	} else if err != nil {
		r.Log.Error(err, "Failed to get Secret", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
		return err
	}

	return nil
}

func (r *MeshReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("Mesh", request.NamespacedName)

	// Fetch the Mesh instance
	instance := &v1alpha1.Mesh{}
	err := r.client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Mesh not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get Mesh")
		return reconcile.Result{}, err
	}

	// Define frontend deployment
	frontendDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend",
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "frontend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "frontend",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "frontend-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "frontend-config",
									},
								},
							},
						},
						corev1.Volume{
							Name: "frontend-secrets",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "frontend-secrets",
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "frontend",
							Image: instance.Spec.FrontendImage,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "frontend-config",
									MountPath: "/etc/frontend",
								},
								corev1.VolumeMount{
									Name:      "frontend-secrets",
									MountPath: "/etc/frontend/secrets",
								},
							},
						},
					},
				},
			},
		},
	}

	// Define backend deployment
	backendDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "backend",
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "backend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "backend",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "backend-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "backend-config",
									},
								},
							},
						},
						corev1.Volume{
							Name: "backend-secrets",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "backend-secrets",
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "backend",
							Image: instance.Spec.BackendImage,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "backend-config",
									MountPath: "/etc/backend",
								},
								corev1.VolumeMount{
									Name:      "backend-secrets",
									MountPath: "/etc/backend/secrets",
								},
							},
						},
					},
				},
			},
		},
	}

	// Define app deployment
	appDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "app",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "app",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "app-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "app-config",
									},
								},
							},
						},
						corev1.Volume{
							Name: "app-secrets",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "app-secrets",
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "app",
							Image: instance.Spec.AppImage,
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "app-config",
									MountPath: "/etc/app",
								},
								corev1.VolumeMount{
									Name:      "app-secrets",
									MountPath: "/etc/app/secrets",
								},
							},
						},
					},
				},
			},
		},
	}

	// Define configmaps
	frontendConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend-config",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": "frontend",
			},
		},
		Data: map[string]string{
			"config.yaml": "frontend configuration",
		},
	}

	backendConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "backend-config",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": "backend",
			},
		},
		Data: map[string]string{
			"config.yaml": "backend configuration",
		},
	}

	appConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-config",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": "app",
			},
		},
		Data: map[string]string{
			"config.yaml": "app configuration",
		},
	}

	// Define secrets
	frontendSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend-secrets",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": "frontend",
			},
		},
		StringData: map[string]string{
			"secret.yaml": "frontend secret",
		},
	}

	backendSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "backend-secrets",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": "backend",
			},
		},
		StringData: map[string]string{
			"secret.yaml": "backend secret",
		},
	}

	appSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-secrets",
			Namespace: instance.Namespace,
			Labels: map[string]string{
				"app": "app",
			},
		},
		StringData: map[string]string{
			"secret.yaml": "app secret",
		},
	}

	// Set Mesh instance as the owner and controller
	if err := ctrl.SetControllerReference(instance, frontendDeployment, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, backendDeployment, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, appDeployment, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, frontendConfigMap, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, backendConfigMap, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, appConfigMap, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, frontendSecret, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, backendSecret, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, appSecret, r.client.Scheme()); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the frontend deployment already exists, if not create a new one
	foundFrontendDeployment := &appsv1.Deployment{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "frontend", Namespace: instance.Namespace}, foundFrontendDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Frontend Deployment", "Deployment.Namespace", frontendDeployment.Namespace, "Deployment.Name", frontendDeployment.Name)
		err = r.client.Create(ctx, frontendDeployment)
		if err != nil {
			log.Error(err, "Failed to create new Frontend Deployment", "Deployment.Namespace", frontendDeployment.Namespace, "Deployment.Name", frontendDeployment.Name)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Frontend Deployment")
		return reconcile.Result{}, err
	}

	// Check if the backend deployment already exists, if not create a new one
	foundBackendDeployment := &appsv1.Deployment{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "backend", Namespace: instance.Namespace}, foundBackendDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Backend Deployment", "Deployment.Namespace", backendDeployment.Namespace, "Deployment.Name", backendDeployment.Name)
		err = r.client.Create(ctx, backendDeployment)
		if err != nil {
			log.Error(err, "Failed to create new Backend Deployment", "Deployment.Namespace", backendDeployment.Namespace, "Deployment.Name", backendDeployment.Name)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Backend Deployment")
		return reconcile.Result{}, err
	}

	// Check if the app deployment already exists, if not create a new one
	foundAppDeployment := &appsv1.Deployment{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "app", Namespace: instance.Namespace}, foundAppDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new App Deployment", "Deployment.Namespace", appDeployment.Namespace, "Deployment.Name", appDeployment.Name)
		err = r.client.Create(ctx, appDeployment)
		if err != nil {
			log.Error(err, "Failed to create new App Deployment", "Deployment.Namespace", appDeployment.Namespace, "Deployment.Name", appDeployment.Name)
			return reconcile.Result{}, err
		}

		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get App Deployment")
		return reconcile.Result{}, err
	}

	// Check if the frontend configmap already exists, if not create a new one
	foundFrontendConfigMap := &corev1.ConfigMap{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "frontend-config", Namespace: instance.Namespace}, foundFrontendConfigMap)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Frontend ConfigMap", "ConfigMap.Namespace", frontendConfigMap.Namespace, "ConfigMap.Name", frontendConfigMap.Name)
		err = r.client.Create(ctx, frontendConfigMap)
		if err != nil {
			log.Error(err, "Failed to create new Frontend ConfigMap", "ConfigMap.Namespace", frontendConfigMap.Namespace, "ConfigMap.Name", frontendConfigMap.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Frontend ConfigMap")
		return reconcile.Result{}, err
	}

	// Check if the backend configmap already exists, if not create a new one
	foundBackendConfigMap := &corev1.ConfigMap{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "backend-config", Namespace: instance.Namespace}, foundBackendConfigMap)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Backend ConfigMap", "ConfigMap.Namespace", backendConfigMap.Namespace, "ConfigMap.Name", backendConfigMap.Name)
		err = r.client.Create(ctx, backendConfigMap)
		if err != nil {
			log.Error(err, "Failed to create new Backend ConfigMap", "ConfigMap.Namespace", backendConfigMap.Namespace, "ConfigMap.Name", backendConfigMap.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Backend ConfigMap")
		return reconcile.Result{}, err
	}

	// Check if the app configmap already exists, if not create a new one
	foundAppConfigMap := &corev1.ConfigMap{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "app-config", Namespace: instance.Namespace}, foundAppConfigMap)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new App ConfigMap", "ConfigMap.Namespace", appConfigMap.Namespace, "ConfigMap.Name", appConfigMap.Name)
		err = r.client.Create(ctx, appConfigMap)
		if err != nil {
			log.Error(err, "Failed to create new App ConfigMap", "ConfigMap.Namespace", appConfigMap.Namespace, "ConfigMap.Name", appConfigMap.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get App ConfigMap")
		return reconcile.Result{}, err
	}

	// Check if the frontend secret already exists, if not create a new one
	foundFrontendSecret := &corev1.Secret{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "frontend-secrets", Namespace: instance.Namespace}, foundFrontendSecret)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Frontend Secret", "Secret.Namespace", frontendSecret.Namespace, "Secret.Name", frontendSecret.Name)
		err = r.client.Create(ctx, frontendSecret)
		if err != nil {
			log.Error(err, "Failed to create new Frontend Secret", "Secret.Namespace", frontendSecret.Namespace, "Secret.Name", frontendSecret.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Frontend Secret")
		return reconcile.Result{}, err
	}

	// Check if the backend secret already exists, if not create a new one
	foundBackendSecret := &corev1.Secret{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "backend-secrets", Namespace: instance.Namespace}, foundBackendSecret)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Backend Secret", "Secret.Namespace", backendSecret.Namespace, "Secret.Name", backendSecret.Name)
		err = r.client.Create(ctx, backendSecret)
		if err != nil {
			log.Error(err, "Failed to create new Backend Secret", "Secret.Namespace", backendSecret.Namespace, "Secret.Name", backendSecret.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get Backend Secret")
		return reconcile.Result{}, err
	}

	// Check if the app secret already exists, if not create a new one
	foundAppSecret := &corev1.Secret{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "app-secrets", Namespace: instance.Namespace}, foundAppSecret)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new App Secret", "Secret.Namespace", appSecret.Namespace, "Secret.Name", appSecret.Name)
		err = r.client.Create(ctx, appSecret)
		if err != nil {
			log.Error(err, "Failed to create new App Secret", "Secret.Namespace", appSecret.Namespace, "Secret.Name", appSecret.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get App Secret")
		return reconcile.Result{}, err
	}

	// Deployment and ConfigMaps/Secrets created successfully
	// Reconciliation is complete
	return reconcile.Result{}, nil
}
