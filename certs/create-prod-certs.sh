#!/bin/bash

certbot certonly \
      --dns-digitalocean \
      --dns-digitalocean-credentials ~/certbot-creds.ini \
      -d 'tlsrpt.alexsci.com'

