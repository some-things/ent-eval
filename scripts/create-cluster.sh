#!/usr/bin/env bash

MYSQL_CHART_VERSION="9.1.4"
PROJECT_NAMESPACE="ent-eval"

main() {
  echo "Creating kind cluster"
  kind create cluster

  echo "Creating namespace"
  kubectl create ns $PROJECT_NAMESPACE

  echo "Deploying mysql chart"
  helm upgrade --install my-mysql bitnami/mysql --version $MYSQL_CHART_VERSION \
    --namespace $PROJECT_NAMESPACE --create-namespace --wait

  echo "Connecting via telepresence"
  telepresence quit && telepresence connect --no-report

  cat <<EOF
Deployment complete - you can now connect to the cluster via telepresence!

MySQL Username: root
MySQL Password: \$(kubectl get secret --namespace ent-eval my-mysql -o jsonpath="{.data.mysql-root-password}" | base64 -d)

export MYSQL_ROOT_PASSWORD=\$(kubectl get secret --namespace ent-eval my-mysql -o jsonpath="{.data.mysql-root-password}" | base64 -d)
mysql -h my-mysql.ent-eval.svc.cluster.local -uroot -p"\$MYSQL_ROOT_PASSWORD"
EOF
}

main
