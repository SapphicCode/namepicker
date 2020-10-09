# build
FROM golang:alpine AS build

RUN apk add --no-cache git
COPY . /project
RUN cd /project && go build -o /bin/namepicker ./cmd/namepicker

# deploy
FROM alpine:latest

RUN apk add --no-cache curl
RUN mkdir /usr/share/names && \
    curl -s https://www.ssa.gov/oact/babynames/names.zip | unzip -d /usr/share/names -
COPY --from=build /bin/namepicker /bin/namepicker

WORKDIR /work
COPY templates /work/templates

EXPOSE 8080/tcp
ENV GIN_MODE=release
CMD ["/bin/namepicker"]
