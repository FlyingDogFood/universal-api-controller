package controllers

import (
	"context"
	"time"

	universalapicontrolleriov1alpha1 "github.com/flyingdogfood/universal-api-controller/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Config controller", func() {

	var config universalapicontrolleriov1alpha1.Config
	//var configTemplate universalapicontrolleriov1alpha1.ConfigTemplate
	//var endpointTemplate universalapicontrolleriov1alpha1.EndpointTemplate

	config = universalapicontrolleriov1alpha1.Config{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "universal-api-controller.io/v1alpha1",
			Kind:       "Config",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-config",
			Namespace: "default",
		},
		Spec: universalapicontrolleriov1alpha1.ConfigSpec{
			Params: []universalapicontrolleriov1alpha1.Param{
				universalapicontrolleriov1alpha1.Param{
					Name:  "host",
					Value: "google.com",
				},
				universalapicontrolleriov1alpha1.Param{
					Name:  "port",
					Value: "443",
				},
			},
			Ref: universalapicontrolleriov1alpha1.ConfigTemplateRef{
				Name: "test-configTemplate",
			},
		},
	}

	/*configTemplate = universalapicontrolleriov1alpha1.ConfigTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "universal-api-controller.io/v1alpha1",
			Kind:       "ConfigTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configTemplate",
			Namespace: "default",
		},
		Spec: universalapicontrolleriov1alpha1.ConfigTemplateSpec{
			Params: []universalapicontrolleriov1alpha1.ParamDefinition{
				universalapicontrolleriov1alpha1.ParamDefinition{
					Name: "host",
				},
				universalapicontrolleriov1alpha1.ParamDefinition{
					Name: "port",
				},
			},
			Reconcile: []universalapicontrolleriov1alpha1.FunctionOrEndpointTemplateRef{
				universalapicontrolleriov1alpha1.FunctionOrEndpointTemplateRef{
					Name: "test-configTemplate",
					Ref: universalapicontrolleriov1alpha1.ObjectRef{
						Name: "test-endpointTemplate",
						Type: "EndpointTemplate",
					},
					Params: []universalapicontrolleriov1alpha1.Param{
						universalapicontrolleriov1alpha1.Param{
							Name:  "host",
							Value: "{{ .Parameters.host }}",
						},
						universalapicontrolleriov1alpha1.Param{
							Name:  "port",
							Value: "{{ .Parameters.port }}",
						},
					},
				},
			},
			Delete: []universalapicontrolleriov1alpha1.FunctionOrEndpointTemplateRef{
				universalapicontrolleriov1alpha1.FunctionOrEndpointTemplateRef{
					Name: "test-configTemplate",
					Ref: universalapicontrolleriov1alpha1.ObjectRef{
						Name: "test-endpointTemplate",
						Type: "EndpointTemplate",
					},
					Params: []universalapicontrolleriov1alpha1.Param{
						universalapicontrolleriov1alpha1.Param{
							Name:  "host",
							Value: "{{ .Parameters.host }}",
						},
						universalapicontrolleriov1alpha1.Param{
							Name:  "port",
							Value: "{{ .Parameters.port }}",
						},
					},
				},
			},
		},
	}

	endpointTemplate = universalapicontrolleriov1alpha1.EndpointTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "universal-api-controller.io/v1alpha1",
			Kind:       "EndpointTemplate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-endpointTemplate",
			Namespace: "default",
		},
		Spec: universalapicontrolleriov1alpha1.EndpointTemplateSpec{
			Params: []universalapicontrolleriov1alpha1.ParamDefinition{
				universalapicontrolleriov1alpha1.ParamDefinition{
					Name: "host",
				},
				universalapicontrolleriov1alpha1.ParamDefinition{
					Name: "port",
				},
			},
			Method: "GET",
			URL:    "https://{{ .Parameters.host }}:{{ .Parameters.port }}/",
		},
	}*/

	Context("When creating Config", func() {
		It("Should create error status as ConfigTemplate is missing", func() {
			By("Creating Congig")
			ctx := context.Background()
			Expect(k8sClient.Create(ctx, &config)).Should(Succeed())
			By("Getting the created Config")
			Eventually(func() bool {
				createdConfig := universalapicontrolleriov1alpha1.Config{}
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "test-config"}, &createdConfig)
				if err != nil {
					return false
				}
				return true
			}, time.Second*10, time.Microsecond*100).Should(BeTrue())
			By("Checking Config Status")
			Eventually(func() string {
				createdConfig := universalapicontrolleriov1alpha1.Config{}
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: "test-config"}, &createdConfig)
				if err != nil {
					return ""
				}
				var status apiextensions.JSON
				input := apiextensionsv1.JSON(createdConfig.Status)
				var scope conversion.Scope
				apiextensionsv1.Convert_v1_JSON_To_apiextensions_JSON(&input, &status, scope)
				switch mystatus := status.(type) {
				case map[string]interface{}:
					switch err := mystatus["error"].(type) {
					case string:
						return err
					}
				}
				return ""
			}, time.Second*10, time.Microsecond*100).Should(BeEquivalentTo("Could not find ConfigTemplate: test-configTemplate"))
		})
	})
})
