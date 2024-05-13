FROM golang:1.22
WORKDIR /src
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /orderService

FROM scratch
COPY --from=0 /orderService /bin/orderService
EXPOSE 8090
CMD ["/bin/orderService"]