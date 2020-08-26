# mytz

Timezones as a service, a simple serverless application using GO

The `main.go` file contains the business logic + AWS realated code (better separation of concerns may be needed for a non-demo application). Application components are defined in `template.yaml`, their lifecycle is completely handled by the CloudFormation service on AWS (no manual activity on the AWS Console).

## Requirements

* [Golang](https://golang.org)
* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) already configured with Administrator permission
* [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [Docker](https://www.docker.com/community-edition) for local testing (optional)

## Build / Run locally / Deploy

The provided `Makefile` wraps the `sam` cli utility to execute several tasks:

- build the go function 
  ```
  make build
  ``` 

- invoke function locally with no parameters for testing 
  ```
  make run
  ``` 

- start local server for testing via a HTTP client
  ```
  make serve
  ``` 

- deploy on AWS (please review deployment parameters first) 
  ```
  make deploy
  ``` 

- remove all resources for the application (please **be aware** that S3 bucket used for deployment must be deleted manually) 
  ```
  make destroy
  ```

If `make` command is not available on your system, use `sam` cli directly (tip: `Makefile` format is quite readable ;)

## Testing

The Lambda function logic can alsbe tested via `go test`

```bash
❯ cd tz
❯ go test
PASS
ok      mytz    0.003s
```

HTTP invocation can be performed locally using the `SAM CLI`:

```bash

# Start local server

❯ make serve

...
Mounting MyTzFunction at http://127.0.0.1:3000/ [GET]
You can now browse to the above endpoints to invoke your functions. You do not need to restart/reload SAM CLI while working on your functions, changes will be reflected instantly/automatically. You only need to restart SAM CLI if you update your AWS SAM templat
2020-08-26 09:51:26  * Running on http://127.0.0.1:3000/ (Press CTRL+C to quit)

# Then invoke the function in a different terminal

❯ http "localhost:3000/tz?time=2019-07-01T12:32&zone=Europe/Rome&other=Europe/Kiev,Pacific/Auckland"
HTTP/1.0 200 OK
Content-Length: 314
Content-Type: application/json
Date: Wed, 26 Aug 2020 08:05:52 GMT
Server: Werkzeug/1.0.1 Python/3.7.8

{
  "ref": {
    "localTime": "01 Jul 19 12:32 UTC",
    "zone": "UTC"
  },
  "others": [
    {
      "localTime": "01 Jul 19 05:32 PDT",
      "zone": "America/Los_Angeles"
    },
    {
      "localTime": "01 Jul 19 08:32 EDT",
      "zone": "America/New_York"
    },
    {
      "localTime": "01 Jul 19 20:32 CST",
      "zone": "Asia/Shanghai"
    },
    {
      "localTime": "01 Jul 19 15:32 +03",
      "zone": "Europe/Istanbul"
    }
  ]
}

```

For more complex scenarios testing on AWS is mandatory ;)