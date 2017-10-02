# golang-app
About

Simple app to test headers passed into the application and one and two way tls connections.

Instructions to deploy as this application.


Run the following commands

```
git clone https://github.com/cf-routing/golang-app.git
cd doctorroute
cf api api.domain.com
cf auth [your_username] [your_password]
cf create-org testorg
cf target -o testorg
cf create-space testspace
cf target -o testorg -s testspace
cf push golang
```

## Run without tls verification

- Update `deployment-vars.yml` by removing `mtls_ca` and `ultimate_ca`
- Run `./generate_manifest deployment-vars.yml`
- Run `cf push potato`

## Provide certs for one-way TLS connections

- Update `deployment-vars.yml` with an `ultimate_ca` comprised of a Certificate Authority and a server cert generated using that CA.
- In `deployment-vars.yml` set `tls_enabled` to `true`
- In `deployment-vars.yml`, remove the `mtls_ca` field from `deployment-vars.yml`
- Run `./generate_manifest deployment-vars.yml -o operations/tls.yml`
- Run `cf push potato`
(If you do not care about the validity of the certificate, you can just run `./generate_manifest potato.yml && cf push potato`)

## Provide certs for MTLS connections

- Update `deployment-vars.yml` with an `ultimate_ca` comprised of a Certificate Authority and a server cert generated using that CA.
- In `deployment-vars.yml` set `tls_enabled` to `true`
- In `deployment-vars.yml`, set `mtls` to true
- In `deployment-vars.yml`, update `mtls_ca` with a CA cert that has been used to generate the certs for `router.backends.cert_chain` and `router.backends.private_key`.
- Run `./generate_manifest deployment-vars.yml -o operations/tls.yml -o operations/mtls.yml`
- Run `cf push potato`

Note that the server certs generated for the backend are generated with the common name `golangSSL`, which the `private_instance_id` of the
route will need to be set to.

This application exposes endpoint `/headers` and prints all the headers received.
