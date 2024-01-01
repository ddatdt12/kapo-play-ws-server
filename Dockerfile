FROM golang:1.21-alpine as builder

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod go.sum ./

# download Go modules and dependencies
RUN go mod download

# copy the code into the container
COPY . .

# compile application
RUN go build -o /main .


# Run the tests in the container
# FROM builder as test
# RUN go test -v ./...


FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=builder /main /

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 3001

ENV PORT=3001 \
    GIN_MODE=release \
    APP_ENV=production

USER nonroot:nonroot

ENTRYPOINT ["/main"]
