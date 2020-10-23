FROM golang
RUN mkdir ./moviespectrum
ADD cmd ./moviespectrum/cmd/
ADD internal ./moviespectrum/internal/
ADD go.mod ./moviespectrum/go.mod
ADD go.sum ./moviespectrum/go.sum
RUN cd ./moviespectrum && go build -o ./rest-service cmd/moviespectrum/rest_service.go

FROM alpine
RUN mkdir /opt/moviespectrum
COPY --from=0 /go/moviespectrum/rest-service /opt/moviespectrum/rest-service
ENTRYPOINT ["/opt/moviespectrum/rest-service"]
