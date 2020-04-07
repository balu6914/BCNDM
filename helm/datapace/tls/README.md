TLS Generation

In order to have ui.datapace.local work over tls:

```openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"```

```kubectl create secret tls ${CERT_NAME} --key ${KEY_FILE} --cert ${CERT_FILE}```

