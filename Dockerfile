FROM alpine:latest

MAINTAINER Edward Muller <a.mihai.nicolae@gmail.com>

WORKDIR "/opt"

ADD .docker_build/ytsync /opt/bin/ytsync
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/ytsync"]

