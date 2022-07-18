FROM golang:buster as build 
WORKDIR /src 
ADD . . 
RUN make run-test && make build-linux

FROM alpine:3.16.0 as pub
WORKDIR /app 

COPY --from=build /src/dist/url-redirection ./
RUN chmod +x ./url-redirection

EXPOSE 80
CMD ["./url-redirection"]