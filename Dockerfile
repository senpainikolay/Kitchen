FROM golang:latest

RUN mkdir /app
COPY . /app 
WORKDIR /app 

RUN export GO111MODULE=on  
RUN go mod tidy 
EXPOSE 8081 
CMD go run src/main.go


 

