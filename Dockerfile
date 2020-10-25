FROM golang:alpine
RUN mkdir ./moviespectrum
ADD cmd ./moviespectrum/cmd/
ADD internal ./moviespectrum/internal/
ADD go.mod ./moviespectrum/go.mod
ADD go.sum ./moviespectrum/go.sum
RUN cd ./moviespectrum \
    && apk add --no-cache ffmpeg-dev gcc musl-dev pkgconfig \
    && go build -o ./rest-service cmd/moviespectrum/rest_service.go

FROM alpine
RUN mkdir /opt/moviespectrum && apk add --no-cache ffmpeg-libs
COPY --from=0 /go/moviespectrum/rest-service /opt/moviespectrum/rest-service
ADD assets /opt/moviespectrum/assets
WORKDIR /opt/moviespectrum
EXPOSE 8080/tcp
ENTRYPOINT ["/opt/moviespectrum/rest-service"]
