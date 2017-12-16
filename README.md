# Git Repository Analyzer

Inspired by x0rz's [tweets analyzer](https://github.com/x0rz/tweets_analyzer), `gitstats` is a simple tool which goal is to perform commits activity analysis on git repositories and visualizing statistics such as as:

* Per user activity distribution.
* Daily commits distribution (per hour). 
* Weekly commits distribution (by day).
* Monthly commits distribution. 
* Yearly commits distribution.
* Words distribution in commits logs.

By default the statistics are generated for every author of the repository but they can be narrowed down filtering with the `-authors email.here@gmail.com,another.one@something.hell` command line parameter.

## Usage

You can download precompiled releases of gitstats [here](https://github.com/evilsocket/gitstats/releases), if instead you want to build it from source, make sure you have Go >= 1.8 installed and configured, then clone this repository, install the dependencies and compile:

    git clone https://github.com/evilsocket/gitstats $GOPATH/src/github.com/evilsocket/gitstats
    cd $GOPATH/src/github.com/evilsocket/gitstats
    make vendor_get
    make

To use `gitstats` simply:

    /path/to/gitstats -repo /path/to/repo

Example output:

![output](https://i.imgur.com/e4kGoAn.png)

## License

gitstats was made with ♥  by [Simone Margaritelli](https://www.evilsocket.net/) and it is released under the GPL 3 license.

