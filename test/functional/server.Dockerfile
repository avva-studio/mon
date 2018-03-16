FROM scratch

COPY ./bin/accounting-rest /

CMD ["/accounting-rest", "serve"]