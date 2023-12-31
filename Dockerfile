FROM golang:1.21

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

COPY . .

RUN go build -o main .

CMD ["./main"]


#cli 
#docker build -t api .
#docker run -p 3000:3000 api