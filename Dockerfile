# syntax=docker/dockerfile:1
FROM golang:1.16
RUN mkdir /app
WORKDIR /app
ADD . /app

RUN make all

# execute program
CMD [ "./chain-indexing"]
