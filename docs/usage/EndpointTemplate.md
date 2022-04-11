# ConfigTemplate
A ConfigTemplate represents a API endpoint template which can can be used in multiple Functions or ConfigTemplates and can makes a http call with the provided parameters.

## Spec

### params
Params is an array of parameters for the EndpointTemplate.  
A parameter is made of `name` and `description`.  
`name`: Name of the parameter.  
`description`: Describes the purpos of the parameter.  

### method
Defines the http method. Can be a go Template and can include parameters for details see [here](Parameters.md).

### url 
Defines the utl for the http call. Can be a go Template and can include parameters for details see [here](Parameters.md).  

### headers
Defines the http headers for the http calls. Headers is a map of strings to strings. The map values can be a go Template and can include parameters for details see [here]c(Parameters.md).  

### body
Defines the http body. Can be a go Template and can include parameters for details see [here](Parameters.md).  

## Samples
```yaml
apiVersion: universal-api-controller.io/v1alpha1
kind: EndpointTemplate
metadata:
  name: endpointtemplate-sample
spec:
  params:
    - name: host
      description: Host to Configure(Can be Cluster intern or extern)
    - name: port
      description: Port of the API
  method: GET
  url: https://{{ .Parameters.host }}:{{ .Parameters.port }}/
```