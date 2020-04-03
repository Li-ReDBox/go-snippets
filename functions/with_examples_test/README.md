# Deploy cloud function with `examples_test.go`

It may not be common to have `examples_test.go` in source files of a cloud function. The common practice I have seen
from [google-cloud-go](https://github.com/googleapis/google-cloud-go) is to name the package in `examples_test.go` to
`package_test`, not `package` as seen in normal test files. But I suspect `gcloud functions deploy` can not exclude
this file. Or maybe the builder on the cloud side cannot handle this, the deployment would fail. The error messages
could not be less helpful very often. The common error message could be as simple as:

`ERROR: (gcloud.functions.deploy) OperationError: code=3, message=Build failed: `

As you can see, it does not even include the actual message. But in one attempt I accidentally get a more meaningful
message (I could not reproduce this error message):

`ERROR: (gcloud.functions.deploy) OperationError: code=3, message=Build failed: Cannot parse function source dir:  Multiple packages in user code directory: demo_test != demo`

I found dashboard may have more information. That means the error is picked up in the remote building process. Which
means sometime `gcloud` did not communicate well with remote, could not get error from it. That's my guess why the
error message sometime not complete.

This error can appear either the code is in a module and just a normal package.

The fix of this is to exclude this file by `.gcloudignore` or use `--ignore-file` argument. 

A funny note, I *accidentally* set `--ignore-file=examples_test.go`, there was no error message. Not very friendly again.

## Error

```shell
gcloud functions deploy ex-test --entry-point=Display --runtime=go113 --trigger-http --verbosity=debug

DEBUG: Running [gcloud.functions.deploy] with arguments: [--entry-point: "Display", --runtime: "go113", --trigger-http: "True", --verbosity: "debug", NAME: "ex-test"]
INFO: Not using ignore file.
INFO: Not using ignore file.
Deploying function (may take a while - up to 2 minutes)...failed.                                                                                                                    
DEBUG: (gcloud.functions.deploy) OperationError: code=3, message=Build failed: 
Traceback (most recent call last):
  File "google-cloud-sdk/lib/googlecloudsdk/calliope/cli.py", line 983, in Execute
    resources = calliope_command.Run(cli=self, args=args)
  File "google-cloud-sdk/lib/googlecloudsdk/calliope/backend.py", line 807, in Run
    resources = command_instance.Run(args)
  File "google-cloud-sdk/lib/surface/functions/deploy.py", line 343, in Run
    return _Run(args, track=self.ReleaseTrack())
  File "google-cloud-sdk/lib/surface/functions/deploy.py", line 297, in _Run
    on_every_poll=[TryToLogStackdriverURL])
  File "google-cloud-sdk/lib/googlecloudsdk/api_lib/functions/util.py", line 318, in CatchHTTPErrorRaiseHTTPExceptionFn
    return func(*args, **kwargs)
  File "google-cloud-sdk/lib/googlecloudsdk/api_lib/functions/util.py", line 369, in WaitForFunctionUpdateOperation
    on_every_poll=on_every_poll)
  File "google-cloud-sdk/lib/googlecloudsdk/api_lib/functions/operations.py", line 151, in Wait
    on_every_poll)
  File "google-cloud-sdk/lib/googlecloudsdk/api_lib/functions/operations.py", line 121, in _WaitForOperation
    sleep_ms=SLEEP_MS)
  File "google-cloud-sdk/lib/googlecloudsdk/core/util/retry.py", line 219, in RetryOnResult
    result = func(*args, **kwargs)
  File "google-cloud-sdk/lib/googlecloudsdk/api_lib/functions/operations.py", line 73, in _GetOperationStatus
    raise exceptions.FunctionsError(OperationErrorToString(op.error))
FunctionsError: OperationError: code=3, message=Build failed: 
ERROR: (gcloud.functions.deploy) OperationError: code=3, message=Build failed: 
```