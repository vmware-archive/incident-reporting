#!/bin/bash

if [ -z "$(which jq)" ]; then
    echo "ERROR: jq must be available on the PATH, but is not."
    exit 1
fi

marker="generated_"
serviceaccount=${1?"must pass a service account"}
namespace=${2:-default}

# look up the service account secret
serviceaccount_secret=$( kubectl get serviceaccount -n $namespace $serviceaccount -o json | jq -r '.secrets[0].name' )

# decode the secrets
serviceaccount_token=$( kubectl get secret -n $namespace $serviceaccount_secret -o json | jq -r '.data.token | @base64d' )
serviceaccount_ca=$( kubectl get secret -n $namespace $serviceaccount_secret -o json | jq -r '.data."ca.crt" | @base64d' )
serviceaccount_namespace=$( kubectl get secret -n $namespace $serviceaccount_secret -o json | jq -r '.data.namespace | @base64d' )

express() {
    indented=$( echo  "${!1}" | sed 's/^/  /')
    echo -e "${marker}${1}: |\n${indented}"
}

express serviceaccount_token
express serviceaccount_ca
express serviceaccount_namespace
