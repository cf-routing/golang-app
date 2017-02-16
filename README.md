# golang-app
About

Simple app to test headers passed into the application.

Instructions to deploy as this application.

Run the following commands

git clone https://github.com/cf-routing/golang-app.git
cd doctorroute
cf api api.domain.com
cf auth [your_username] [your_password]
cf create-org testorg
cf target -o testorg
cf create-space testspace
cf target -o testorg -s testspace
cf push golang

After staging the app successfully, run the command `cf app golang` to get the URL for the app.

```
curl -vv golang.cfdomain.com
> Host: golang.cfdomain.com
> User-Agent: curl/7.50.1
> Accept: */*
> X-Forwarded-Host:test
>
< HTTP/1.1 200 OK
< Content-Length: 562
< Content-Type: text/plain; charset=utf-8
< Date: Thu, 16 Feb 2017 19:01:49 GMT
< X-Vcap-Request-Id: 5b6fbcef-5a30-4cc6-7045-bb627473f997
go world!
```

This application exposes endpoint `/headers` and prints all the headers received.
