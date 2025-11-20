#!/bin/bash

set -e

export CAROOT=`pwd`

# Remove all old leaf certs
rm *alexsci*.pem || true
rm -Rf ./invalid-certs/

create_cert() {
    DOMAIN=$1
    DN_PREFIX=$2
    CERT_PREFIX=$3
    ~/code/mkcert/mkcert "${DN_PREFIX}${DOMAIN}"
    mkdir -p ./invalid-certs/live/$DOMAIN
    mv "${CERT_PREFIX}${DOMAIN}-key.pem" ./invalid-certs/live/$DOMAIN/privkey.pem
    mv "${CERT_PREFIX}${DOMAIN}.pem"     ./invalid-certs/live/$DOMAIN/fullchain.pem
}

create_cert "tlsrpt.alexsci.com" "" ""

