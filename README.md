# SingleStore Helm Chart

This repository contains a Helm chart for deploying SingleStore (MemSQL) via a Kubernetes Operator.

Refer to [Deploying SingleStore](https://docs.singlestore.com/docs/deploy-kubernetes/) and the [SingleStore Operator Reference](https://docs.singlestore.com/docs/operator-reference/) for more information.

## Chart Parameters

All of the parameters that can be specified in this chart are described in the chart's [README](singlestore/README.md). 

## Testing

This repository provides snapshot tests for this chart under [template_test](tests/template/template_test.go).
