FROM golang:latest

RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get -u google.golang.org/api/drive/v3
RUN go get -u golang.org/x/oauth2/...
RUN go get -u github.com/jteeuwen/go-bindata/...

ADD ./src/api /go/src/api

#Binarize client secret key...
RUN go-bindata -pkg data -o ./src/api/app/data/data.go ./src/api/app/data/
RUN cp ./src/api/app/data/client_secret.json .
RUN rm ./src/api/app/data/client_secret.json

RUN go install api
RUN go test -v api/app/items
RUN go test -v api/app/files
RUN go test api/app/auth
CMD ["api"]
