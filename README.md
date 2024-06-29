

# Dataform formatter

Format `.sqlx` files in your Dataform project

Setup necessary files & default `.sqluff` config file in `.formatdataform` directory to support formatting *( needs to be done just once for a Dataform project )*
Setup if automatically called if user tries to format files without calling it explicitly

```bash
formatdataform setup
```

**Format a file or directory**

```bash
formatdataform format <path_to_file_or_directory>
```

To override the default config file

```bash
formatdataform -c <path_to_sqlfluff_config_file> format <path_to_file_or_directory>
```


### [Installation](#installation)

Install [sqlfluff](https://github.com/sqlfluff/sqlfluff)

```
pip install sqlfluff
```

**Latest release**

```
curl -sSfL https://raw.githubusercontent.com/ashish10alex/formatdataform/main/install_latest.sh | bash
```
**OR**

```bash
go install github.com/ashish10alex/formatdataform@latest
```
This installs the binary `formatdataform` to $GOBIN, which defaults to $GOPATH/bin.

**OR**

Manually clone the repository and build the cli and add the cli to your system path

```bash
git clone https://github.com/ashish10alex/formatdataform.git
go build -o formatdataform
mv formatdataform /usr/local/bin/formatdataform

```

### To run tests

1. Install [gotestsum](https://github.com/gotestyourself/gotestsum?tab=readme-ov-file#output-format) for prettier test visualization
2. Run `gotestsum --format testname`
