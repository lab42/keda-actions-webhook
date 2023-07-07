FROM alpine as setup
RUN addgroup --gid 1234 -S appgroup && \
    adduser --uid 1234 -S appuser -G appgroup && \
    apk --no-cache add ca-certificates

FROM scratch as production
COPY --from=setup /etc/passwd /etc/passwd
COPY --from=setup /etc/ssl /etc/ssl
COPY keda-actions-webhook /keda-actions-webhook
USER appuser
EXPOSE 1234
ENTRYPOINT ["/keda-actions-webhook"]
