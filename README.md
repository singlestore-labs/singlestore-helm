# SingleStore Helm Chart

This repository contains Helm chart for installing SingleStore via Kubernetes Operator. For more information on Kubernetes Operator please refer to [SingleStore Kubernetes deployment](https://docs.singlestore.com/db/v9.0/deploy/kubernetes/) and [SingleStore Operator Reference](https://docs.singlestore.com/db/v9.0/reference/singlestore-operator-reference/).

## Chart parameters

All of the parameters that can be specified in the chart are described in chart's [README](singlestore/README.md). 

## Testing

Repository also provides snapshot tests for the chart under [template_test](tests/template/template_test.go)