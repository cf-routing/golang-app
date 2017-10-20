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
- When no cipher suite is provided, the [golang default](https://golang.org/pkg/crypto/tls/#pkg-constants) ciphers are loaded.
- When the set of ciphers are specified and shared with the gorouter, the most preferred cipher (first one in the list) is used.
- In `manifest.yml`
    ```yaml
    CIPHER: "ECDHE-RSA-AES128-GCM-SHA256:TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
    ```
