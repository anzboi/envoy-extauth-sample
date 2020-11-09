# Envoy extauth sample

This repo is used to experiment with envoy extauth configuration

## Run

You will need go and docker-compose installed

```sh
# build binaries and run envoy, auth.go and service.go
make up

# bring down (if cmd-C doesn't do it)
docker-compose down
```

## Testing

The service exposes a basic greeter service defined in [proto/hello.proto](./proto/hello.proto) along with an http gateway and reflection.

```sh
# reflection
grpcurl localhost:8000 list

# grpcurl (with some headers)
grpcurl -H "authhorization: bearer <token>" -H "x-b3-traceid: abcde" -d '{"name": "bob"}' localhost:8000 hello.Greeter/GetGreeting

# curl through http gateway
curl -H "authhorization: bearer <token>" -H "x-b3-traceid: abcde" http://localhost:8000/v1/greeting?name=bob
```

You will see a lot of logs for these request, this is what you should expect to see

1. Envoy logs indicating the request is passing through envoy
1. Auth server logs indicating details of the auth request that it recieved
1. Service logs indicating that the end service recieved a request (with an extra header added by the auth service).

## Information about extauth

All envoy config for this running example is found in `front-envoy.yaml`. See the [envoy docs](https://www.envoyproxy.io/docs/envoy/latest/api-v2/config/filter/http/ext_authz/v2/ext_authz.proto) for more information

Extauth is envoys method of delegating request authentication and authorization to an externa auth service (hence the name extauth). To do this, envoy is configured with an `envoy.filters.http.ext_authz` filter in its filter chain.

```yaml
- name: envoy.filters.http.ext_authz
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz
    http_service:
      server_uri:
        uri: authz:8080
        cluster: ext-authz
        timeout: 0.25s
      authorization_request:
        headers_to_add:
          - key: Method
            value: "%REQ(:METHOD)%"
          - key: Path
            value: "%REQ(:PATH)%"
      authorization_response:
        allowed_upstream_headers:
          patterns:
            - exact: authorization
```

When this filter is activated, envoy sends the request as is (almost) to the extauth server. This server is responsible for determining whether the request should be allowed through or blocked.

The original request is sent to the auth server NOT modified EXCEPT for the method and path (these need to be modified to hit the auth server). The `authorization_request` block adds the original path and method to the headers sent to the auth server.

The auth server can add headers to its response, and the auth filter will add those headers in before sending the request down to the target service. These allowed headers are listed in the `allows_upstream_headers` block. This allows the auth server to add or change headers (for example, exchange an opaque token for an unencrypted JWT).

Finally, we configure this cluster in envoy to tell it where the auth service is being hosted

```yaml
- name: ext-authz
  connect_timeout: 0.25s
  type: strict_dns
  lb_policy: round_robin
  hosts:
  - socket_address:
      address: authz
      port_value: 8080
```
