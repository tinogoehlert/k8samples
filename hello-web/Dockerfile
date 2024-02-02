FROM golang as builder

WORKDIR /

COPY . .

RUN go get .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# deployment image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

LABEL author="Tino Göhlert"

WORKDIR /root/
COPY --from=builder /app .

CMD [ "./app" ]

EXPOSE 8080
