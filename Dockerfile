FROM scratch

COPY ./bin/accounting-rest-serve /

CMD ["/accounting-rest-serve"]