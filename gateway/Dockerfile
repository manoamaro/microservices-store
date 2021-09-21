FROM nginx

RUN apt-get update && apt-get install unzip

RUN rm -v /etc/nginx/conf.d/*

ENV CT_URL https://releases.hashicorp.com/consul-template/0.27.0/consul-template_0.27.0_linux_amd64.zip
RUN curl -O $CT_URL
RUN unzip consul-template_0.27.0_linux_amd64.zip -d /usr/local/bin

ADD nginx.conf /etc/consul-templates/nginx.conf
ADD entrypoint.sh /bin/entrypoint.sh
RUN chmod +x /bin/entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/bin/entrypoint.sh"]