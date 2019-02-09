FROM golang:latest as build

RUN go get github.com/gorilla/mux && mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o runner .

FROM scratch

COPY --from=build /app/runner /runner
WORKDIR /
EXPOSE 8080
CMD ["/runner"]