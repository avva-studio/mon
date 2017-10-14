FROM scratch

COPY ./bin/gohmoney-rest /

CMD ["/gohmoney-rest", "/.gohmoney/.gohmoneydbconnectionstring"]