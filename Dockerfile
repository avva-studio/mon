FROM scratch

COPY ./bin/monserve /

CMD ["/monserve"]