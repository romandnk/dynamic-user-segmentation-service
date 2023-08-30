FROM golang:1.20 as build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o "./bin/dynamic-user-segmentation" ./cmd/dynamic-user-segmentation

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/dynamic-user-segmentation ./bin/

COPY ./configs/dynamic-user-segmentation.yaml ./configs/
COPY ./static/reports ./static/reports

CMD ["./bin/dynamic-user-segmentation", "-config", "./configs/dynamic-user-segmentation.yaml"]