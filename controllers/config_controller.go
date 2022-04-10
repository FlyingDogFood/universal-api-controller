/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	universalapicontrolleriov1alpha1 "github.com/flyingdogfood/universal-api-controller/api/v1alpha1"
)

// ConfigReconciler reconciles a Config object
type ConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=universal-api-controller.io,resources=configs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=universal-api-controller.io,resources=configs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=universal-api-controller.io,resources=configs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Config object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	requeueAfter, _ := time.ParseDuration("60s")

	log.Info("Reconciling")
	// Get config object
	var config universalapicontrolleriov1alpha1.Config
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		log.Error(err, "unable to fetch Config")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Create Status with current Timestamp
	status := newStatus()

	finalizerName := "universal-api-controller.io/finalizer"

	// Check if Object is under pending deletion and add finalizer if not present
	pendingDeletion, err := r.handleFinalizer(ctx, config, finalizerName)
	if err != nil {
		log.Error(err, "failed to add finalizer")
		status["error"] = "Could not add finalizer"
		if suberr := r.updateStatus(ctx, config, status); suberr != nil {
			log.Error(suberr, "Failed to update status")
			return ctrl.Result{RequeueAfter: requeueAfter}, suberr
		}
		return ctrl.Result{RequeueAfter: requeueAfter}, err
	}

	// Get configTemplate object
	var configTemplate universalapicontrolleriov1alpha1.ConfigTemplate
	if err := r.Get(ctx, types.NamespacedName{Name: config.Spec.Ref.Name, Namespace: config.Namespace}, &configTemplate); err != nil {
		log.Error(err, "unable to fetch ConfigTemplate")
		status["error"] = "Could not find ConfigTemplate: " + config.Spec.Ref.Name
		status["details"] = err
		if suberr := r.updateStatus(ctx, config, status); suberr != nil {
			log.Error(suberr, "Failed to update status")
			return ctrl.Result{RequeueAfter: requeueAfter}, suberr
		}
		return ctrl.Result{RequeueAfter: requeueAfter}, err
	}
	log.Info("Fetched ConfigTemplate", "ConfigTemplate", configTemplate)

	var functionsOrEndpointTemplates []universalapicontrolleriov1alpha1.FunctionOrEndpointTemplateRef
	if pendingDeletion {
		functionsOrEndpointTemplates = configTemplate.Spec.Delete
	} else {
		functionsOrEndpointTemplates = configTemplate.Spec.Reconcile
	}

	parameters := newParameters()
	parameters, err = parameters.generateParameters(config.Spec.Params)
	if err != nil {
		log.Error(err, "Failed to generate Parameters")
		status["error"] = "Failed to generate Parameters"
		status["details"] = err
		if suberr := r.updateStatus(ctx, config, status); suberr != nil {
			log.Error(suberr, "Failed to update status")
			return ctrl.Result{RequeueAfter: requeueAfter}, suberr
		}
		return ctrl.Result{RequeueAfter: requeueAfter}, err
	}

	_, status, err = r.executeFunctionsOrEndpointTemplates(ctx, functionsOrEndpointTemplates, parameters, config.Namespace)
	if err != nil {
		if suberr := r.updateStatus(ctx, config, status); suberr != nil {
			log.Error(suberr, "Failed to update status")
			return ctrl.Result{RequeueAfter: requeueAfter}, suberr
		}
		return ctrl.Result{RequeueAfter: requeueAfter}, err
	}

	if pendingDeletion {
		controllerutil.RemoveFinalizer(&config, finalizerName)
		if err := r.Update(ctx, &config); err != nil {
			log.Error(err, "unable to delete finalizer")
			status["error"] = "Could not remove Finalizer"
			if suberr := r.updateStatus(ctx, config, status); suberr != nil {
				log.Error(suberr, "Failed to update status")
				return ctrl.Result{RequeueAfter: requeueAfter}, suberr
			}
			return ctrl.Result{RequeueAfter: requeueAfter}, err
		}
		// return and don't requeue as Object is deleted
		return ctrl.Result{}, nil
	}

	if err := r.updateStatus(ctx, config, status); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{RequeueAfter: requeueAfter}, err
	}

	log.Info("Finished Reconciliation", "RequeAfter", requeueAfter.String())
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

func (r *ConfigReconciler) executeFunctionsOrEndpointTemplates(ctx context.Context, functionsOrEndpointTemplates []universalapicontrolleriov1alpha1.FunctionOrEndpointTemplateRef, parameters Parameters, namespace string) (Parameters, Status, error) {
	log := log.FromContext(ctx)

	status := make(Status)

	for _, functionOrEndpointTemplate := range functionsOrEndpointTemplates {
		funcParameters, err := parameters.generateParameters(functionOrEndpointTemplate.Params)
		if err != nil {
			log.Error(err, "Could not generate Parameters for "+functionOrEndpointTemplate.Name)
			status["error"] = "Could not generate Parameters for " + functionOrEndpointTemplate.Name
			return parameters, status, err
		}
		if functionOrEndpointTemplate.Ref.Type == "Function" {
			var function universalapicontrolleriov1alpha1.Function
			if err := r.Get(ctx, types.NamespacedName{Name: functionOrEndpointTemplate.Ref.Name, Namespace: namespace}, &function); err != nil {
				log.Error(err, "unable to fetch Function "+functionOrEndpointTemplate.Ref.Name)
				status["error"] = "Could not find Function: " + functionOrEndpointTemplate.Ref.Name
				return parameters, status, err
			}
			resParameters, resStatus, err := r.executeFunctionsOrEndpointTemplates(ctx, function.Spec.Actions, funcParameters, namespace)
			parameters.merge(resParameters, function.Name)
			status.merge(resStatus, function.Name)
			if err != nil {
				log.Error(err, "Failed to execute Function "+function.Name)
				status["error"] = "Failed to execute Function " + function.Name
				return parameters, status, err
			}
		} else if functionOrEndpointTemplate.Ref.Type == "EndpointTemplate" {
			var endpointTemplate universalapicontrolleriov1alpha1.EndpointTemplate
			if err := r.Get(ctx, types.NamespacedName{Name: functionOrEndpointTemplate.Ref.Name, Namespace: namespace}, &endpointTemplate); err != nil {
				log.Error(err, "unable to fetch EndpointTemplate")
				status["error"] = "Could not find EndpointTemplate: " + functionOrEndpointTemplate.Ref.Name
				return parameters, status, err
			}
			httpResponse, status, err := r.executeEndpointTemplate(ctx, endpointTemplate, funcParameters)
			parameters.Responses[functionOrEndpointTemplate.Name] = httpResponse
			status.merge(status, functionOrEndpointTemplate.Name)
			if err != nil {
				log.Error(err, "Failed to execute EndpointTemplate "+endpointTemplate.Name)
				status["error"] = "Failed to execute EndpointTemplate " + endpointTemplate.Name
				return parameters, status, err
			}
		} else {
			err := errors.New("Not supported type: " + functionOrEndpointTemplate.Ref.Type)
			log.Error(err, "Type is not supported")
			status["error"] = "Not Supportet Type: " + functionOrEndpointTemplate.Ref.Type
			return parameters, status, err
		}
	}
	return parameters, status, nil
}

func (r *ConfigReconciler) executeEndpointTemplate(ctx context.Context, endpointTemplate universalapicontrolleriov1alpha1.EndpointTemplate, parameters Parameters) (HttpResponse, Status, error) {
	log := log.FromContext(ctx)
	status := make(Status)
	var err error

	log.Info("executing EndpointTemplate: " + endpointTemplate.Name)

	httpMethod, err := templateString(endpointTemplate.Spec.Method, parameters)
	if err != nil {
		log.Error(err, "Cannot template httpMethod")
		status["error"] = "Cannot template httpMethod"
		status["details"] = err
		return HttpResponse{}, status, err
	}

	httpURL, err := templateString(endpointTemplate.Spec.URL, parameters)
	if err != nil {
		log.Error(err, "Cannot template httpURL")
		status["error"] = "Cannot template httpURL"
		status["details"] = err
		return HttpResponse{}, status, err
	}

	httpBody, err := templateString(endpointTemplate.Spec.Body, parameters)
	if err != nil {
		log.Error(err, "Cannot template httpBody")
		status["error"] = "Cannot template httpBody"
		status["details"] = err
		return HttpResponse{}, status, err
	}

	httpRequest, _ := http.NewRequest(httpMethod, httpURL, strings.NewReader(httpBody))

	for key, value := range endpointTemplate.Spec.Headers {
		httpHeaderValue, err := templateString(value, parameters)
		if err != nil {
			log.Error(err, "Cannot template Header: "+key)
			status["error"] = "Cannot template Header" + key
			return HttpResponse{}, status, err
		}

		httpRequest.Header.Add(key, httpHeaderValue)
	}

	httpClient := http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		log.Error(err, "Failed doing HTTP Request")
		status["error"] = "Failed doing HTTP Request"
		status["details"] = err
		return HttpResponse{}, status, err
	}
	response, err := fromHttpResponse(*httpResponse)
	if err != nil {
		log.Error(err, "Failed phrasing HTTP Body")
		status["error"] = "Failed phrasing HTTP Body"
		status["details"] = err
	}
	return response, status, err
}

func (r *ConfigReconciler) handleFinalizer(ctx context.Context, config universalapicontrolleriov1alpha1.Config, finalizerName string) (bool, error) {
	// Check if Config is under Deletion
	if config.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&config, finalizerName) {
			controllerutil.AddFinalizer(&config, finalizerName)
			return false, r.Update(ctx, &config)
		}
		return false, nil
	}
	return true, nil
}

func (r *ConfigReconciler) updateStatus(ctx context.Context, config universalapicontrolleriov1alpha1.Config, status Status) error {
	var err error
	config.Status.Raw, err = status.bytes()
	if err != nil {
		return err
	}
	return r.Status().Update(ctx, &config)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&universalapicontrolleriov1alpha1.Config{}).
		Complete(r)
}
