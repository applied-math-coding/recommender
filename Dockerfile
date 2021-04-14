FROM node:14 AS build_client
WORKDIR /client/
COPY ./client .
RUN yarn install --production
RUN yarn build

FROM golang:1.16 AS build_go
WORKDIR /src/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch as bin
EXPOSE 8080
WORKDIR /root/
COPY --from=build_go /src/app .
COPY --from=build_client /client/build ./client/build
CMD ["./app"]

