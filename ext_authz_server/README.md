# ext_authz_server

## Notes
1. `authv3.CheckResponse.Status` should use grpc codes but not http status code, and use `google.golang.org/genproto/googleapis/rpc/status` to fill in.
2. `authv3.CheckResponse.Status` definitively means **ALLOW** or **DENY** but not the `Status` in `authv3.CheckResponse.HttpResponse`
3. When the reqeust is denied by external authorization server, if you want to get not only http status code but response body, you can set it in the body of `authv3.CheckResponse_DeniedResponse`:
    ```go
    return &authv3.CheckResponse{
        // ...
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
            // ...
			DeniedResponse: &authv3.DeniedHttpResponse{
				Body:   fmt.Sprintf(`{"code": %s, "message": "request denied"}`, code),
			},
		},
	}
    ```