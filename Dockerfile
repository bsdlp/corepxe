FROM dock0/service
MAINTAINER Jon Chen <bsd@voltaire.sh>

EXPOSE 8080
VOLUME ["/var/run/docker.sock"]

ADD corepxe /usr/local/bin/corepxe

ADD run /service/corepxe/run
