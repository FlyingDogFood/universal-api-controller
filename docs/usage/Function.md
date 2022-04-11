# Function
A Function defines a list of actions to execute.

## Spec

### params
Params is an array of parameters for the Function.  
A parameter is made of `name` and `description`.  
`name`: Name of the parameter.  
`description`: Describes the purpos of the parameter.  

### actions
Is a list of Functions or EnpointTemplates to execute when the Function is executed.  
For configuration consult [Function or EndpointTemplate Reference](Function_or_EndpointTemplate_Reference.md).  

## Samples
```yaml
apiVersion: universal-api-controller.io/v1alpha1
kind: Function
metadata:
  name: function-sample
spec:
  params:
    - name: host
      description: Host to Configure(Can be Cluster intern or extern)
    - name: port
      description: Port of the API
  actions:
    - name: create
      ref:
        name: endpointtemplate-sample
        kind: EndpointTemplate
      params:
        - name: host
          value: {{ .Parameters.host }}
        - name: port
          value: "{{ .Parameters.port }}"
```