# ConfigTemplate
A ConfigTemplate represents a API configuration template which can combine multiple API calls to get the API to desired state.  

## Spec

### params
Params is an array of parameters for the ConfigTemplate.  
A parameter is made of `name` and `description`.  
`name`: Name of the parameter.  
`description`: Describes the purpos of the parameter.  

### reconcile
Is a list of Functions or EnpointTemplates to execute when a Config referencing the ConfigTemplate is created. It's also repeatedly executed to ensure the configuration keeps in desired state.  
For configuration consult [Function or EndpointTemplate Reference](Function_or_EndpointTemplate_Reference.md).  

### delete
Is a list of Functions or EnpointTemplates to execute when a Config referencing the ConfigTemplate is deleted.  
For configuration consult [Function or EndpointTemplate Reference](Function_or_EndpointTemplate_Reference.md).  

## Samples
```yaml
apiVersion: universal-api-controller.io/v1alpha1
kind: ConfigTemplate
metadata:
  name: configtemplate-sample
spec:
  params:
    - name: host
      description: Host to Configure(Can be Cluster intern or extern)
    - name: port
      description: Port of the API
  reconcile:
    - name: create-update
      ref:
        name: function-sample
        kind: Function
      params:
        - name: host
          value: {{ .Parameters.host }}
        - name: port
          value: "{{ .Parameters.port }}"
  delete:
    - name: delete
      ref:
        name: endpointtemplate-sample
        kind: EndpointTemplate
      params:
        - name: host
          value: {{ .Parameters.host }}
        - name: port
          value: "{{ .Parameters.port }}"
```