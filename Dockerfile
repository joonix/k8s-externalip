FROM scratch

ADD k8s-externalip-init /app/init
ENTRYPOINT [ "/app/init" ]
