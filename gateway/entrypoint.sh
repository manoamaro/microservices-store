#!/bin/bash
service nginx start
consul-template -consul=$CONSUL_URL -template="/etc/consul-templates/nginx.conf:/etc/nginx/conf.d/service.conf:service nginx reload"