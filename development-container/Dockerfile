FROM golang:1.23.4-alpine

# add git and necessary build tools
RUN apk add --no-cache git curl

# install air for hot-reloading
RUN curl -sSf https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

# set working directory
WORKDIR /app

# copy go mod and sum files
COPY ../go.mod ../go.sum ./

# download dependencies
RUN go mod download

# copy the source code
COPY .. .

# expose the port your application runs on
EXPOSE 8088

# use air for hot-reloading in development
CMD ["air"]
