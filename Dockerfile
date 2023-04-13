FROM golang:1.20 as build

# This might need to be adjusted based on device that used
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64
RUN apt update && apt install -y build-essential

WORKDIR /app
COPY . .
COPY ./config ./config

RUN go mod download
RUN CGO_ENABLED=0 go build -o app cmd/main.go
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest

FROM mcr.microsoft.com/playwright:v1.30.0-focal

# TINI used to prevent zombie process and memory leaks on running browser
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

WORKDIR /app

RUN mkdir -p /var/log/behance
RUN chmod 775 /var/log/behance

COPY --from=build /go/bin/playwright /app/playwright

COPY --from=build /app/app .
COPY --from=build /app/config ./config

RUN /app/playwright install

EXPOSE 8080

ENTRYPOINT ["/tini", "--"]
CMD ["/app/app"]