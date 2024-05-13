# syntax=docker/dockerfile:1
FROM golang:1.22
WORKDIR /src
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /paymentService 

FROM scratch
COPY --from=0 /paymentService /bin/paymentService
CMD ["/bin/paymentService"]