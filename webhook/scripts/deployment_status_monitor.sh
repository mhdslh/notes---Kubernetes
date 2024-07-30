#! /bin/bash

function wait() {
    until [[ -n $(kubectl get deployment $1 -n $2 -o jsonpath='{range .status.conditions[*]}{"["}{.type}{":"}{.status}{"]"}{end}' | grep "Available:True") ]]
    do
        echo "namespace:$2 deployment:$1 is not ready yet"
        sleep 2
    done
}

function main() {
    local short_opts="h"
    local long_opts="help,ns:"
    local opts
    local namespace="default"
    local deployment=""

    if [ "${#}" -lt 1 ]; then
        echo "Usage: deployment_status_monitor.sh <deployment name> [--ns <namespace name>]"
        exit 1
    fi

    if ! opts=$(getopt -s bash -o "${short_opts}" -l "${long_opts}" -- "${@}"); then
        echo "Usage: deployment_status_monitor.sh <deployment name> [--ns <namespace name>]"
        exit 1
    fi
    eval set -- "${opts}"

    while true; do
        case $1 in
            -h | --help)
                echo "Usage: deployment_status_monitor.sh <deployment name> [--ns <namespace name>]"
                exit 0
                ;;
            --ns)
                namespace=$2
                shift 2
                ;;
            --)
                shift
                if [ -z $1 ]
                then
                    echo "<deployment name> is required"
                    echo "Usage: deployment_status_monitor.sh <deployment name> [--ns <namespace name>]"
                    exit 1
                fi
                deployment=$1
                break
                ;;
        esac
    done

    time wait ${deployment} ${namespace}

    exit 0
}

main "${@}"