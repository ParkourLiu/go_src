FROM ubuntu:16.04

ENV TZ "Asia/Shanghai"
RUN echo $TZ > /etc/timezone && \
    apt-get update && apt-get install -y tzdata && \
    rm /etc/localtime && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

ADD muser /
EXPOSE 8888
ENTRYPOINT /muser -prop=/etc/prop.properties