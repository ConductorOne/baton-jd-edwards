FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-jd-edwards"]
COPY baton-jd-edwards /