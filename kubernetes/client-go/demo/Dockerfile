FROM golang:1.22 as builder

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . .

RUN CGO_ENABLED=0 go build -o ingress-manager main.go

FROM alpine:3.17.9

WORKDIR /app

COPY --from=builder /app/ingress-manager .

CMD ["./ingress-manager"]