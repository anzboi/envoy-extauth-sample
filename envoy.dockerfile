FROM hub.artifactory.gcp.anz/envoyproxy/envoy:v1.14.4

CMD /usr/local/bin/envoy -c /etc/front-envoy.yaml --service-cluster front-proxy