

# Dataform formatter

Format `.sqlx` files in your Dataform project using [sqlfluff](https://github.com/sqlfluff/sqlfluff)

[![Version](https://img.shields.io/github/v/release/ashish10alex/formatdataform)](https://github.com/ashish10alex/formatdataform/releases)
![Linux](https://img.shields.io/badge/Linux-supported-success)
![macOS](https://img.shields.io/badge/macOS-supported-success)


**To format a file or directory run**

```bash
formatdataform format <path_to_file_or_directory>
```

> [!NOTE]
When ran for the first time on a dataform workspace, the formatdataform cli will setup necessary files & default `.sqluff` config file in `.formatdataform` directory
to support formatting. You can alternatively manually setup the defaults by running `formatdataform setup`.

You can override the default sqlfluff config file by using the `--config` flag or `-c` (shorthand) as follows

```bash
formatdataform -c <path_to_sqlfluff_config_file> format <path_to_file_or_directory>
```

Or alternatively, you can directly edit the sqlfluff config file generated in `.formatdataform/.sqlfluff`

> [!IMPORTANT]
> Ensure that the config block has [following lines](https://github.com/ashish10alex/formatdataform/blob/b265412ed05fbca21620e520713ac4faf1d61e82/cmd/sqlfluff_config.go#L116-L117) which handles
the regex for `[sqlfluff:templater:placeholder]` to handle the parsing of `${ref("TABLE_NAME")}` blocks in `.sqlx` files


### Installation

1. **Prerequisite:** Install [sqlfluff](https://github.com/sqlfluff/sqlfluff)

```
pip install sqlfluff
```

2. **Install the Latest release of `formatdataform` binary**

```
curl -sSfL https://raw.githubusercontent.com/ashish10alex/formatdataform/main/install_latest.sh | bash
```

**OR**

```bash
go install github.com/ashish10alex/formatdataform@latest
```
This installs the binary `formatdataform` to `$GOBIN`, which defaults to `$GOPATH/bin`

**OR**

Manually clone the repository and build the cli and add the cli to your system path

```bash
git clone https://github.com/ashish10alex/formatdataform.git
go build -o formatdataform
mv formatdataform /usr/local/bin/formatdataform
```
