# svcf
Adds go-flags support to voi-svc

Uppdates https://github.com/voi-oss/svc

with an pre-init-stage that parses github.com/jessevdk/go-flags

Also add a nullworker which can be embedded in service which only want the init-stage (blocks until terminated.)