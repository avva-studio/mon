FROM scratch

COPY ./bin/main /

CMD ["/main", "/.gohmoney/.gohmoneydbconnectionstring"]