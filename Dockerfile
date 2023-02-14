FROM alpine

RUN apk --update add redis
RUN apk add --no-cache python3 py3-pip
RUN apk add --no-cache git make musl-dev go


# for nvm
# RUN apk add -U curl bash ca-certificates openssl ncurses coreutils python3 make gcc g++ libgcc linux-headers grep util-linux binutils findutils
# RUN touch ~/.profile
# RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash

# for caddy
# RUN apk add caddy
# RUN caddy start


RUN wget https://github.com/rawdaGastan/ginit/releases/download/v0.1/ginit_0.1_Linux_x86_64.tar.gz
RUN tar xzf ginit_0.1_Linux_x86_64.tar.gz && mv ginit /bin

COPY ./init.sh .
RUN [ "chmod", "+x", "/init.sh"]
ENTRYPOINT [ "/init.sh" ]