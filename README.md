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