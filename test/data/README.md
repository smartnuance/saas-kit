Create new RS256 public/private key-pair for signing tokens in auth service and validating them in any service.

> ssh-keygen -t rsa -b 4096 -m PEM -f ./test/data/jwtRS256.key

> openssl rsa -in ./test/data/jwtRS256.key -pubout -outform PEM -out ./test/data/jwtRS256.key.pub
