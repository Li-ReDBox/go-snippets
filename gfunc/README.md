### local test environment

`cmd` directory contains local run set up.

If you have no dependencies, or not using module, you may try to use relative paths in `cmd/main.go`. Otherwise, you have to use `replace` to import cloud function to be tested.

### Deploy

Maintain a minimal deployment by adding what should not be uploaded:
if you have dependencies, do not ignore go.mod

gcloud would be able to find the entry-point function in the current package. It seems it does not care about package's name.

gcloud functions deploy test --entry-point Demo --trigger-http --runtime go111
