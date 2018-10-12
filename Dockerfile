# build stage
FROM golang:alpine AS build-env
ADD . /src
RUN apk add make git
RUN cd /src && make install && make build

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/bin/mail_service /app/
CMD ./mail_service