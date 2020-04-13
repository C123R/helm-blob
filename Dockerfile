FROM scratch
COPY helm-blob /
ENTRYPOINT ["/helm-blob"]