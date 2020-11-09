FROM golang:1.13.7-buster
WORKDIR /usr/src/app
COPY bin/service /bin/service
CMD [ "/bin/service" ]