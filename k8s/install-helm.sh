#!/bin/bash

set -e

CHART_NAME="majesticcoding"
NAMESPACE="streaming"

# Package the Helm chart
helm lint ./helm/$CHART_NAME
helm dependency update ./helm/$CHART_NAME
helm install $CHART_NAME ./helm/$CHART_NAME --create-namespace --namespace $NAMESPACE

echo "Helm chart $CHART_NAME installed to namespace $NAMESPACE."