FROM golang:1.22.10 as builder

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY . .

RUN CGO_ENABLED=0 go build -o ingress-manager main.go

FROM busybox:1.37.0

WORKDIR /app

COPY --from=builder /app/ingress-manager .

CMD ["./ingress-manager"]