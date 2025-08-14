# ╔═════════════════════════════════════════════════════╗
# ║ IMAGE                                               ║
# ╚═════════════════════════════════════════════════════╝
FROM scratch
WORKDIR /app
COPY eui64-calculator .
USER 65534:65534
EXPOSE 8080
ENTRYPOINT ["./eui64-calculator"]
