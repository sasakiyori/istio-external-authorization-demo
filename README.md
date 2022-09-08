# Istio External Authorization Server Demo
A demo for istio grpc external authorization server.

## Precondition
Please ensure you have installed k8s and istio related.

## Usage
1. run [`ext_authz_server/pack.sh`](ext_authz_server/pack.sh) to build a external authorization image to docker hub. Note that you should change the docker user name. You can use [my image](https://hub.docker.com/r/sasakiyori/ext-authz-server/tags) built by this repository's code as well.
    ```bash
    #!/bin/bash

    # ext_authz_server/pack.sh
    # change to your docker hub user name
    USERNAME="sasakiyori"
    ```

2. apply [`config/ext-authz-server.yaml`](config/ext-authz-server.yaml) to build the `Service` and `Deployment` of external authorization server:  `kubectl apply -f ext-authz-server.yaml`

3. config map processing
   - `kubectl edit configmap istio -n istio-system`
   -  add `extensionProviders` config by [`config/istio-config-map.yaml`](config/istio-config-map.yaml)
   - `kubectl rollout restart deployment/istiod -n istio-system`

4. add istio ingress config by [`config/istio-ingressgateway.yaml`](config/istio-ingressgateway.yaml)

5. check if your external authorization server runs normally

## References
- [istio sample](https://github.com/istio/istio/tree/master/samples/extauthz)
- [authorization-policy](https://istio.io/latest/docs/reference/config/security/authorization-policy/)
- [envoy auth proto](https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/auth/v3/external_auth.proto)
