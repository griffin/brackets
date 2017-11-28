FROM scratch
#ADD ca-certificates.crt /etc/ssl/certs/
ADD main /
ADD public /
ADD config.yml /
CMD ["/main"]

EXPOSE 8080