FROM alpine

RUN apk --update add redis
RUN apk add --no-cache python3 py3-pip
RUN apk add --no-cache git make musl-dev go



RUN wget https://github.com/rawdaGastan/ginit/releases/download/v0.1/ginit_0.1_Linux_x86_64.tar.gz
RUN tar xzf ginit_0.1_Linux_x86_64.tar.gz && mv ginit /bin

COPY ./init.sh .
RUN [ "chmod", "+x", "/init.sh"]
ENTRYPOINT [ "/init.sh" ]