FROM golang:1.19 as build

ENV NAME "go_scheduler"
WORKDIR /opt/${NAME}
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/${NAME} ./cmd/api/main.go

FROM scratch
ENV NAME "go_scheduler"
COPY --from=build /opt/${NAME}/bin/${NAME} /${NAME}
COPY --from=build /opt/${NAME}/configs/dev.env /dev.env
COPY --from=build /opt/${NAME}/schemas /schemas

CMD ["./go_scheduler"]