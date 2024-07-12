

test:
    gotestsum --format testname

watch_test:
    gotestsum --format testname --watch

build:
    go build .
    sudo mv formatdataform /usr/local/bin/formatdataform
