# Start by building the application.
FROM golang:1.19.2 as build
LABEL stage=dockerbuilder
WORKDIR /app/example
COPY . .

# Build the binary
RUN make build

# Now copy it into our base image.
FROM alpine:3.9

# Copy bin file
COPY --from=build /app/example/dist/example /app/example
RUN mkdir /logs

EXPOSE 8000
ENTRYPOINT ["/app/example"]