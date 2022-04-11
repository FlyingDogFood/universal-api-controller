# Config
A Config represents a API configuration to make.  
Creating a Config object executes a ConfigTemplate with the defined parameters.  

## Spec

### ref
`name`: Defines the name of the ConfigTemplate in the same namespace.  
### params
Params is an array of parameters to define to execute the ConfigTemplate.  
A parameter is made of `name` and `value`.  
`name`: Name of the parameter, has to be the same as in ConfigTemplate.  
`value`: Value of the Parameter, is used for execution of ConfigTemplate.  

## Samples
```yaml
apiVersion: universal-api-controller.io/v1alpha1
kind: Config
metadata:
  name: config-sample
spec:
  ref:
    name: configtemplate-sample
  params:
    - name: host
      value: localhost
    - name: port
      value: "8080"
```