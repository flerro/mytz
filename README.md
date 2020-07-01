# mytz

Timezones as a service 

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Installing dependencies & building the target 

The `sam` commands are wrapped inside of the `Makefile`. To execute this simply run

    make build

## Invoke the service

Invoking function locally through local API Gateway:

```bash
❯ make serve

Mounting MyTzFunction at http://127.0.0.1:3000/tz [GET]
You can now browse to the above endpoints to invoke your functions. You do not need to restart/reload SAM CLI while working on your functions, changes will be reflected instantly/automatically. You only need to restart SAM CLI if you update your AWS SAM template
2020-07-01 23:03:36  * Running on http://127.0.0.1:3000/ 

...

❯ curl "localhost:3000/tz?time=2019-07-01T12:32&zone=Europe/Rome&other=Europe/Kiev,Pacific/Auckland" | jq .
{
  "referenceTime": "01 Jul 19 12:32 CEST",
  "timezones": {
    "Europe/Kiev": "01 Jul 19 13:32 EEST",
    "Pacific/Auckland": "01 Jul 19 22:32 NZST"
  }
}

```

## Deployment 

Modify parameters (e.g. S3 bucket name, aws profile) in the `Makefile`, then run: 

    make deploy