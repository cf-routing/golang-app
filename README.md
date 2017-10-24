# an app for testing (m)TLS with the cf router

## Without any TLS
```
cf push golang-app
```

## Using one-way TLS (server cert only)
Before running `cf push`
- In `manifest.yml` set
  ```yaml
  TLS: true
  ```
- Include `server.crt` and `server.key` files in this directory

## Using mTLS (server cert and expecting client cert)
Before running `cf push`
- In `manifest.yml`
    ```yaml
    TLS: true
    MTLS: true
    ```
- Include `server.crt` and `server.key` files in this directory
- Include a `ca.crt` file in this directory (CA that signed the expected client certs)

## Specifying Cipher suites while TLS or MTLS is enabled
- When `CIPHER` key is not provided in `manifest.yml` the [golang default](https://golang.org/pkg/crypto/tls/#pkg-constants) ciphers are loaded.
- To specify a cipher suite for the app, use `CIPHER` in `manifest.yml`; Gorouter must be configured to support the same cipher suite. Note: only one cipher suite is supported by `CIPHER`. For e.g. in `manifest.yml`
 
   ```yaml
    CIPHER: "ECDHE-RSA-AES128-GCM-SHA256"
    ```
