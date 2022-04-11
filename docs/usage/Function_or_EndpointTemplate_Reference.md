# Function or EndpointTemplate Reference
Defines a Function or EndpointTemplate to execute with the defined parameters.  

## name
Defines a name for the step.  

## ref
Defines a reference to a Function or Endpointtemplate. It consistsof `name` and `kind`.  
`name`: Is the name of the EndpointTemplate or Function in the same namespace.  
`kind`: Is `Function` or `EndpointTemplate`  

## params
Defines a list of parameters for the Execution of the Function or EndpointTemplate. With the parameters `name` and `value`.  
`name`: Name of the parameter, has to be the same as in the referenced Function or EndpointTemplate.  
`value`: Value of the Parameter, is used for execution of Function or EndpointTemplate. Can be a go Template and can include parameters or http responses, for details see [here](Parameters.md).  

## Samples
```yaml
name: create-update
ref:
  name: function-sample
  kind: Function
params:
  - name: host
    value: {{ .Parameters.host }}
  - name: port
    value: "{{ .Parameters.port }}"
```