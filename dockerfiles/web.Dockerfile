FROM golang:1.24 AS build-stage

WORKDIR /workdir
COPY app/go.mod app/go.sum ./
RUN go mod download
COPY app/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app

WORKDIR /workdir2
COPY health/go.mod ./
RUN go mod download
COPY health/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /health


####
FROM gcr.io/distroless/base-debian12 AS build-release-stage

WORKDIR /templates
COPY app/templates/*.html app/templates/*tmpl .

WORKDIR /
COPY --from=build-stage /app /app
COPY --from=build-stage /health /health

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
HEALTHCHECK --interval=10s --timeout=1s --start-interval=1s CMD ["/health"]

