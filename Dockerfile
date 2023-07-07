FROM alpine as setup
RUN addgroup --gid 1234 -S appgroup && \
    adduser --uid 1234 -S appuser -G appgroup

FROM scratch as production
COPY --from=setup /etc/passwd /etc/passwd
COPY user-svc /user-svc
USER appuser
EXPOSE 1234
ENTRYPOINT ["/user-svc"]
