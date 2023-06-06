FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
COPY ./func/ /func
WORKDIR /func

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /func/async-blob-copy .

RUN rm ./async-blob-copy.go
RUN rm ./go.sum
RUN rm ./go.mod
RUN rm ./local.settings.json

FROM mcr.microsoft.com/azure-functions/dotnet:4.9.1-appservice

ENV AzureWebJobsScriptRoot=/home/site/wwwroot \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true

COPY --from=builder /func/ /home/site/wwwroot

EXPOSE 8080
ENTRYPOINT [ "/home/site/wwwroot/async-blob-copy" ]
