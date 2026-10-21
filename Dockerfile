FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/tracker ./cmd/tracker
RUN mkdir /data

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/tracker /tracker
COPY --from=build --chown=65532:65532 /data /data
ENV STATE_FILE=/data/state.json
ENTRYPOINT ["/tracker"]
