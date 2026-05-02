FROM deloo/markdown-docs AS builder
#RUN apk add util-linux
RUN apk add perl-file-rename --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing/

COPY . /home/src
WORKDIR /home/src
RUN file-rename "s/README.MD/index.md/" *
RUN file-rename "s/README.MD/index.md/" */*
RUN find . -name README.MD 
ENV WORKSPACE=/home
ENV TITLE="RabbitMQ курс"
ENV LANGUAGE=ru
ENV ICON=rabbit-variant-outline
RUN makedocs ./src ./dst

FROM nginx:alpine

COPY --from=builder /home/dst /usr/share/nginx/html/
WORKDIR /usr/share/nginx/html/
