FROM scratch
WORKDIR /app
COPY rancher-info-agent /usr/bin/
ENTRYPOINT ["rancher-info-agent"]