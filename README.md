# universal-api-controller
A kubernetes controller to declarative controll the configuration of applications that don't have their own crds for configuration.
This controllers enables you to make a series of API calls to controll a state of a cluster external or internal APIs, without the need to write your own kubernetes controller.

## Docs
- [Documentation](docs)
    - [Installation](docs/install)
    - [Usage](docs/usage)

## What's missing
- Tests
- Documentation
- Probably still full of bugs

## Use Cases
### Configure external oder internal API configuration
You can use the universal-api-controller to controll the state of Cluster external or internal APIs to a defined state.
Examples:
- Controll the user assignment to groups in an IAM system declaratively over Config objects.

### Use as an Operator
You can use the universal-api-controller as a replacement to write a operator. You can define multiple ConfigTemplates vor different Actions/Configurations.
Examples:
- You have deployed a software/application over a helm chart e.g.(Hashicorp Vault). Now you can define EndpointTemplates, Functions and ConfigTemplates to for example create Namespaces in vault by creating Config Objects.