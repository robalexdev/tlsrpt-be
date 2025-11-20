#!/usr/bin/bash

set -e

docker compose --env-file dev.env build

# Full rebuild
docker compose --env-file dev.env down -t 0 || true
# Purge old DB
docker volume rm tlsrpt-be_tlsrpt-db || true

docker compose --env-file dev.env up -d --wait --remove-orphans

EMAIL="$(uuidgen | tr -d '-')@example.com"
PASSWORD=$(uuidgen)
TEST_DOMAIN=göpher.test
TEST_DOMAIN_NORMALIZED=xn--gpher-jua.test
echo "User ${USERNAME} ${PASSWORD}"

COOKIE_JAR=cookies.test

rm ${COOKIE_JAR} || true

echo "Create test user" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie-jar ${COOKIE_JAR} \
  https://localhost:9443/signup \
  -d "email=${EMAIL}&password1=${PASSWORD}&password2=${PASSWORD}" > run.log

echo "Ensure logged in" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/ | tee -a run.log | grep "${EMAIL}"

echo "Log out" | tee -a run.log
curl \
  -k \
  -X POST \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/signout >> run.log

echo "Ensure logged out" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/ | tee -a run.log | egrep -v "${EMAIL}"

echo "Try the wrong password" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie-jar ${COOKIE_JAR} \
  https://localhost:9443/signin \
  -d "email=${EMAIL}&password=123456789" | tee -a run.log | grep "Invalid username or password"

echo "Log back in" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie-jar ${COOKIE_JAR} \
  https://localhost:9443/signin \
  -d "email=${EMAIL}&password=${PASSWORD}" >> run.log

echo "Ensure logged in (again)" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/ | tee -a run.log | grep "${EMAIL}"


echo "Add a domain" | tee -a run.log
curl \
  -k \
  -L \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/domain/add \
  -d "domain=${TEST_DOMAIN}" | tee -a run.log | grep "${TEST_DOMAIN_NORMALIZED}"

echo "Ensure not validated yet" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/domain/1/ | tee -a run.log | grep "Check Domain"

echo "Validate it (test domain auto validates)" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  -X POST \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/domain/1/ >> run.log

echo "Ensure validated now" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/domain/1/ | tee -a run.log | grep "Enabled"

echo "Upload a report (json)" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/uploadReport \
  --form "file=@report.json" >> run.log

echo "Ensure the uploaded TLS report is visible" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/domain/1/ | tee -a run.log | grep "report-id-from-report.json"

echo "Ensure not an open relay" | tee -a run.log
./test-open-relay.exp 127.0.0.1 | tee -a run.log | grep "Relay access denied"

echo "Email a TLSRPT" | tee -a run.log
./test-send-email.exp 127.0.0.1 d-1 tlsrpt.alexsci.com | tee -a run.log | grep "Ok: queued as"

# Wait for postfix queue to process the mail
sleep 2

echo "Ensure the emailed TLS report is visible" | tee -a run.log
curl \
  -k \
  -H "Host: tlsrpt.alexsci.com" \
  --cookie ${COOKIE_JAR} \
  https://localhost:9443/domain/1/ | tee -a run.log | grep "report-id-of-report-over-email-1234"

echo "SUCCESS"

