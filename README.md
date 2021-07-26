# `tbd`

[![Go Report Card](https://goreportcard.com/badge/github.com/lucasepe/tbd?style=flat-square)](https://goreportcard.com/report/github.com/lucasepe/tbd) &nbsp;&nbsp;&nbsp;
[![Release](https://img.shields.io/github/release/lucasepe/tbd.svg?style=flat-square)](https://github.com/lucasepe/tbd/releases/latest) &nbsp;&nbsp;&nbsp;
[![codecov](https://codecov.io/gh/lucasepe/tbd/branch/main/graph/badge.svg?style=flat-square)](https://codecov.io/gh/lucasepe/tbd)

_"to be defined"_

## A really simple way to create text templates with placeholders.

This tool is deliberately simple and trivial, no advanced features. 

> If you need advanced templates rendering which supports complex syntax and a huge list of datasources (JSON, YAML,  AWS EC2 metadata, BoltDB, Hashicorp > Consul and Hashicorp Vault secrets), I recommend you use one of these:
>
> - [gotemplate](https://github.com/hairyhenderson/gomplate)
> - [pongo2](https://github.com/flosch/pongo2)
> - [quicktemplate](https://github.com/valyala/quicktemplate)

## Built-in Variables

When executed inside a Git repository, `tbd` automatically exports some variables related to the Git repository which may be useful in the build phase.

These variables are: `ARCH`, `OS`, `REPO_COMMIT`, `REPO_HOST`, `REPO_NAME`, `REPO_ROOT`, `REPO_TAG`, `REPO_TAG_CLEAN`, `REPO_URL`, `TIMESTAMP`.

Try it! With `tbd` in your `PATH`, go in a Git folder and type:

```sh
$ tbd vars
+----------------+------------------------------------------+
| ARCH           | amd64                                    |
| OS             | linux                                    |
| REPO_COMMIT    | a3193274112d3a6f5c2a0277e2ca07ec238d622f |
| REPO_HOST      | github.com                               |
| REPO_NAME      | tbd                                      |
| REPO_ROOT      | lucasepe                                 |
| REPO_TAG       | v0.1.1                                   |
| REPO_TAG_CLEAN | 0.1.1                                    |
| REPO_URL       | https://github.com/lucasepe/tbd          |
| TIMESTAMP      | 2021-07-26T14:22:36Z                     |
+----------------+------------------------------------------+
```

> Obviously in your case the values ​​will be different.

## How does a template looks like ?

A template is a text document in which you can insert placeholders for the text you want to make dynamic.

- a placeholder is delimited by `{{` and `}}` - (i.e. `{{ FULL_NAME }}`)
- all text outside placeholders is copied to the output unchanged

Example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: {{ metadata.name }}
  labels:
    app: {{ metadata.labels.app }}
 spec:
  containers:
    - name: {{ container.1.name }}
      image: {{ container.1.image }}
      ports:
        - containerPort: {{ container.1.port }}
    - name: {{ container.2.name }}
      image: {{ container.2.image }}
      ports:
        - containerPort: {{ container.2.port }}
```

Another example:

```txt
{{ greeting }}

I will be out of the office from {{ start.date }} until {{ return.date }}. 
If you need immediate assistance while I’m away, please email {{ contact.email }}.

Best,
{{ name }}
```

## How can I define placeholders values?

Create a text file in which you enter the values for the placeholders.

- define a placeholder value using `KEY = value` (or `KEY: value`)
- empty lines are skipped
- lines beginning with `#` are treated as comments

Example:

```sh
# metadata values
metadata.name = rss-site
metadata.labels.app = web

# containers values
container.1.name = front-end
container.1.image = nginx
container.1.port = 80

container.2.name = rss-reader
container.2.image: nickchase/rss-php-nginx:v1
container.2.port: 88
```

Another example...

```sh
greeting: Greetings
start.date: August, 9 
return.date: August 23
contact.email: pinco.pallo@gmail.com
name: Pinco Pallo 
```

## How fill in the template?

> Use the `merge` command

```sh
$ tbd merge /path/to/your/template /path/to/your/envfile
```

Example:

```sh
$ tbd merge testdata/sample.tbd testdata/sample.vars
```

👉 you can also specify an HTTP url to fetch your template and/or placeholders values.

Example:

```sh
$ tbd merge https://raw.githubusercontent.com/lucasepe/tbd/main/testdata/sample.tbd \
    https://raw.githubusercontent.com/lucasepe/tbd/main/testdata/sample.vars
```

and the output is...

```txt
Greetings

I will be out of the office from August, 9 until August 23. 
If you need immediate assistance while I’m away, please email pinco.pallo@gmail.com.

Best,
Pinco Pallo
```

## How to list all template placeholders?

> Use the `marks` command.

```sh
$ tbd marks /path/to/your/template
```

Example:

```sh
$ tbd marks testdata/sample.tbd
greeting
start.date
return.date
contact.email
name
```

## How to list all variables?

> Use the `vars` command.

```sh
$ tbd vars /path/to/your/envfile
```

Example:

```sh
$ tbd vars testdata/sample.vars
+----------------+------------------------------------------+
| Label          | Value                                    |
+----------------+------------------------------------------+
| ARCH           | amd64                                    |
| OS             | linux                                    |
| REPO_COMMIT    | a3193274112d3a6f5c2a0277e2ca07ec238d622f |
| REPO_HOST      | github.com                               |
| REPO_NAME      | tbd                                      |
| REPO_ROOT      | lucasepe                                 |
| REPO_TAG       | v0.1.1                                   |
| REPO_TAG_CLEAN | 0.1.1                                    |
| REPO_URL       | https://github.com/lucasepe/tbd          |
| TIMESTAMP      | 2021-07-26T14:17:49Z                     |
| contact.email  | pinco.pallo@gmail.com                    |
| greeting       | Greetings                                |
| name           | Pinco Pallo                              |
| return.date    | August 23                                |
| start.date     | August, 9                                |
+----------------+------------------------------------------+
```

> As you can see, since I ran the command in a Git repository, there are also relative variables.

# How to install?

If you have [golang](https://golang.org/dl/) installed:

```sh
$ go install github.com/lucasepe/tbd@latest
```

This will create the executable under your `$GOPATH/bin` directory.

## Ready-To-Use Releases 

If you don't want to compile the sourcecode yourself, [here you can find the tool already compiled](https://github.com/lucasepe/tbd/releases/latest) for:

- MacOS
- Linux
- Windows

<br/>

#### Credits

Thanks to [@valyala](https://github.com/valyala/) for the [fasttemplate](https://github.com/valyala/fasttemplate) library - which I have modified by adding and removing some functions for the `tbd` purpose.
