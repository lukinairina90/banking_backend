FROM golang:1.19.3-alpine as builder

RUN mkdir /build
COPY . /build
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/banking-app cmd/main.go

# generate clean, final image for end users
FROM alpine:3.11.3
COPY --from=builder /bin/banking-app /bin/banking-app
COPY docker/rbac/model.conf /rbac/model.conf
COPY docker/rbac/policy.csv /rbac/policy.csv

# executable
ENTRYPOINT [ "/bin/banking-app" ]