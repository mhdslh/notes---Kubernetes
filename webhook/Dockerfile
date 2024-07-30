FROM ubuntu

COPY server/server .
COPY cert/ca-cert.pem .
COPY cert/ca-key.pem .

ENTRYPOINT ["/server"]