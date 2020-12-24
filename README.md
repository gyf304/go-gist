gist(1) -- upload code to https://gist.github.com
=================================================

Inspired by https://github.com/defunkt/gist

## Synopsis

The gist gem provides a `gist` command that you can use from your terminal to
upload content to https://gist.github.com/.

## Installation

Go to Releases to view available releases.
Download Archive for your OS and architecture and uncompress to your PATH.

## Command

‌To upload the contents of `a.go` just:

    gist a.go

‌Upload multiple files:

    gist a.go b.go c.go
    gist *.go

‌Use `-p` to make the gist private:

    gist -p a.go

‌Use `-d` to add a description:

    gist -d "description" a.go

To read a gist and print it to STDOUT

    gist -r GIST_ID
    gist -r 374130

You can also use `-o` to specify a output file.

‌See `gist --help` for more options.

## Login

You will be prompted to login if `gist` does not have your login information already.
