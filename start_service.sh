#!/bin/sh
/bin/service &
envoy -c /etc/service-envoy.yaml --service-cluster service${SERVICE_NAME}