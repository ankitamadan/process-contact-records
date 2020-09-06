process-contact-records
------------

A small lambda designed to read off a the connect trace records from the kinesis stream and insert them into the database.

Environment variables:
* DB_HOST
* DB_NAME
* DB_USER
* DB_PASS

## Building

_Remember to build your handler executable for Linux!_
```
GOOS=linux GOARCH=amd64 go build -o main
zip main.zip main
```

## Deploying

```
aws lambda update-function-code --function-name process-contact-trace-records --zip-file fileb://main.zip
```
