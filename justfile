

test:
    gotestsum --format testname

build:
    go build .
    sudo mv formatdataform /usr/local/bin/formatdataform
