FROM gcr.io/distroless/static:nonroot
COPY ginko        /usr/local/bin/ginko
COPY ginko-admin  /usr/local/bin/ginko-admin
COPY memserver    /usr/local/bin/memserver
COPY memmcp       /usr/local/bin/memmcp
COPY memctl       /usr/local/bin/memctl

VOLUME ["/data"]
ENV GINKO_DB=/data/ginko.db
EXPOSE 8787

ENTRYPOINT ["/usr/local/bin/ginko"]
CMD ["serve", "--db", "/data/ginko.db"]
