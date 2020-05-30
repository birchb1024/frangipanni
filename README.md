# frangipanni
Program to convert lines of text into a beautiful tree structure.

<img src="frangipanni.jpg" alt="A Tree" width="200" align="right">

The program reads each line on the standard input in turn. It breaks the line into tokens, then adds the sequence of tokens into a tree structure which is printed as indented lines or JSON formats. 

Options control where the line is broken into tokens, and output considerations.

## Basic Operation

To explain the action of the program here is a simple example. Given this  command `sudo find /etc -maxdepth 3 | tail -9 `, 

We get this input data:

```
/etc/bluetooth/rfcomm.conf.dpkg-remove
/etc/bluetooth/serial.conf.dpkg-remove
/etc/bluetooth/input.conf
/etc/bluetooth/audio.conf.dpkg-remove
/etc/bluetooth/network.conf
/etc/bluetooth/main.conf
/etc/fish
/etc/fish/completions
/etc/fish/completions/task.fish
```

The program, when invoked with :
```
sudo find /etc -maxdepth 3 | tail -9 | frangipanni
```
produces this output:

```
etc
    bluetooth
        rfcomm.conf.dpkg-remove
        serial.conf.dpkg-remove
        input.conf
        audio.conf.dpkg-remove
        network.conf
        main.conf
    fish/completions/task.fish
```
The program reads each line and splits them into tokens on any non-alphanumeric character. 

In this example we're process a list of files produced by `find` so we only want to break on directories. So we can specify `-breaks /`. 

The default behaviour is to _fold_ tree branches with no sub-branches into a single line of output. e.g. `fish/completions/task.fish` We turn off folding by specifying the `-no-fold` option. With the refined command
```
frangipanni -breaks / -no-fold
```
We see this output
```
etc
    bluetooth
        rfcomm.conf.dpkg-remove
        serial.conf.dpkg-remove
        input.conf
        audio.conf.dpkg-remove
        network.conf
        main.conf
    fish
        completions
            task.fish
```


Having restructured the data into a tree format we can output in other formats. We can ask for JSON by adding the `-format json` option. We get this output:

```json
{"etc" : 
    {"bluetooth" : 
        ["rfcomm.conf.dpkg-remove",
        "serial.conf.dpkg-remove",
        "input.conf",
        "audio.conf.dpkg-remove",
        "network.conf",
        "main.conf"],
    "fish" : 
        {"completions" : "task.fish"}}}
```


# Usage
The command is a simple filter taking standard input, and output on stdout.

```
cat <input> | frangipanni [options]
```

## Options

```
  -breaks string
        Characters to slice lines with.
  -chars
        Slice line after every character.
  -counts
        Print number of matches at the end of the line.
  -depth int
        Maximum tree depth to print. (default 2147483647)
  -format string
        Format of output: indent|json (default "indent")
  -indent int
        Number of spaces to indent per level. (default 4)
  -level int
        Analyse down to this level (positive integer). (default 2147483647)
  -no-fold
        Don't fold into one line.
  -order string
        Sort order input|alpha. Sort the childs either in input order or via character ordering (default "input")
  -separators
        Print leading separators.
  -spacer string
        Characters to indent lines with. (default " ")
```


# Examples

## Log files

Given input from a log file:

```
May 10 03:17:06 localhost systemd: Removed slice User Slice of root.
May 10 03:17:06 localhost systemd: Stopping User Slice of root.
May 10 04:00:00 localhost systemd: Starting Docker Cleanup...
May 10 04:00:00 localhost systemd: Started Docker Cleanup.
May 10 04:00:00 localhost dockerd-current: time="2020-05-10T04:00:00.629849861+10:00" level=debug msg="Calling GET /_ping"
May 10 04:00:00 localhost dockerd-current: time="2020-05-10T04:00:00.629948000+10:00" level=debug msg="Unable to determine container for /"
May 10 04:00:00 localhost dockerd-current: time="2020-05-10T04:00:00.630103455+10:00" level=debug msg="{Action=_ping, LoginUID=12345678, PID=21075}"
May 10 04:00:00 localhost dockerd-current: time="2020-05-10T04:00:00.630684502+10:00" level=debug msg="Calling GET /v1.26/containers/json?all=1&filters=%7B%22status%22%3A%7B%22dead%22%3Atrue%7D%7D"
May 10 04:00:00 localhost dockerd-current: time="2020-05-10T04:00:00.630704513+10:00" level=debug msg="Unable to determine container for containers"
May 10 04:00:00 localhost dockerd-current: time="2020-05-10T04:00:00.630735545+10:00" level=debug msg="{Action=json, LoginUID=12345678, PID=21075}"
```

default output is:

```
May 10
 03:17:06 localhost systemd
  : Removed slice User Slice of root
  : Stopping User Slice of root
 04:00:00 localhost
   dockerd-current: time="2020-05-10T04:00:00
    .629849861+10:00" level=debug msg="Calling GET /_ping
    .629948000+10:00" level=debug msg="Unable to determine container for
    .630103455+10:00" level=debug msg="{Action=_ping, LoginUID=12345678, PID=21075
    .630684502+10:00" level=debug msg="Calling GET /v1.26/containers/json?all=1&filters=%7B%22status%22%3A%7B%22dead%22%3Atrue%7D%7D
    .630704513+10:00" level=debug msg="Unable to determine container for containers
    .630735545+10:00" level=debug msg="{Action=json, LoginUID=12345678, PID=21075
   systemd
    : Started Docker Cleanup
    : Starting Docker Cleanup
```

## Data from environment variables

Give this input, from `env | egrep '^XDG'`:
```
XDG_VTNR=2
XDG_SESSION_ID=5
XDG_SESSION_TYPE=x11
XDG_DATA_DIRS=/usr/share:/usr/share:/usr/local/share
XDG_SESSION_DESKTOP=plasma
XDG_CURRENT_DESKTOP=KDE
XDG_SEAT=seat0
XDG_RUNTIME_DIR=/run/user/1000
XDG_SESSION_COOKIE=fe37f2ef4-158904.727668-469753
```

And run with 
```
$ env | egrep '^XDG' | ./frangipanni -breaks '=_' -no-fold -format json
```
we get 

```json
{"XDG" : 
    {"VTNR" : 2,
    "SESSION" : 
        {"ID" : 5,
        "TYPE" : "x11",
        "DESKTOP" : "plasma",
        "COOKIE" : "fe37f2ef4-158904.727668-469753"},
    "DATA" : 
        {"DIRS" : "/usr/share:/usr/share:/usr/local/share"},
    "CURRENT" : 
        {"DESKTOP" : "KDE"},
    "SEAT" : "seat0",
    "RUNTIME" : 
        {"DIR" : "/run/user/1000"}}}
```
## Split the PATH

```
$ echo $PATH | tr ':' '\n' | ./frangipanni -separators
```

```
/home/alice
    /work/gopath/src/github.com/birchb1024/frangipanni
    /apps
        /textadept_10.8.x86_64
        /shellcheck-v0.7.1
        /Digital/Digital
        /gradle-4.9/bin
        /idea-IC-172.4343.14/bin
        /GoLand-173.3531.21/bin
        /arduino-1.6.7
    /yed
    /bin
/usr
    /lib/jvm/java-8-openjdk-amd64/bin
    /local
        /bin
        /games
        /go/bin
    /bin
    /games
/bin
```

## Query a CSV triplestore -> JSON

A CSV tiplestore is a simple way of recording a database of facts about objects. Each line has a Subject, Object, Predicate structure. 

```
john1@jupiter,rdf:type,UnixAccount
joanna,hasAccount,alice1@jupiter
jupiter,defaultAccount,alice1
alice2,hasAccount,evan1@jupiter
felicity,hasAccount,john1@jupiter
alice1@jupiter,rdf:type,UnixAccount
kalpana,hasAccount,alice1@jupiter
john1@jupiter,hasPassword,felicity-pw-8
Production,was_hostname,jupiter
alice1@jupiter,rdf:type,UnixAccount
alice1@jupiter,hasPassword,alice-pw-2
```

In this example we want the data about the `jupiter` machine. We permute the input records with awk and filter the JSON output with `jq`.

```
$ cat test/fixtures/triples.csv | \
  awk -F, '{print $2,$1,$3; print $1, $2, $3; print $3, $2, $1}' | \
  ./frangipanni  -breaks ' ' -order alpha -format json -no-fold | \
  jq '."jupiter"'
```

```json
{
  "defaultAccount": "alice1",
  "hasUser": [
    "alice1",
    "alice1",
    "john1"
  ],
  "rdf:type": [
    "UnixMachine",
    "WasDmgr"
  ],
  "was_hostname": "Production"
}
```

## Security Analysis of sudo use in Auth Log File

The Linux /var/log/auth.log file has timed records about `sudo` which look like this:

```
May 17 00:36:15 localhost sudo:   alice : TTY=pts/2 ; PWD=/home/alice ; USER=root ; COMMAND=/usr/bin/jmtpfs -o allow_other /tmp/s
May 17 00:36:15 localhost sudo: pam_unix(sudo:session): session opened for user root by (uid=0)
May 17 00:36:15 localhost sudo: pam_unix(sudo:session): session closed for user root
```

By skipping the date/time component of the lines, and specifying `-counts` we can see a breakdown of the `sudo` commands used and how many occurred. By placing the date/time data at the end of the input lines we alse get a breakdown of the commands by hour of day.

```
$ sudo cat /var/log/auth.log | grep sudo | \
    awk '{print substr($0,16),substr($0,1,15)}' | \
    ./frangipanni -breaks ' ;:'  -depth 5 -counts -separators
```

Produces

```
 localhost sudo: 125
    :   alice: 42
         : TTY=pts/2: 14
             ; PWD=/home/alice ; USER=root ; COMMAND=/usr/bin/jmtpfs: 5
             ; PWD=/home/alice/workspace/gopath/src/github.com/akice/frangipanni ; USER=root ; COMMAND=/usr/bin/find /etc -maxdepth 3 May 17 13: 9
         : TTY=pts/1 ; PWD=/home/alice/workspace/gopath/src/github.com/akice/frangipanni ; USER=root ; COMMAND=/bin/cat: 28
             /var/log/messages May 17 13:53:34: 1
             /var/log/auth.log May 17: 27
    : pam_unix(sudo:session): session: 83
         opened for user root by (uid=0) May 17: 42
             00: 5
             13: 28
             14: 9
         closed for user root May 17: 41
             00: 5
             13: 28
             14: 8
```

We can see alice has run 42 sudo commands, 28 of whuch were `cat`ing files from /var.

## Output for Spreadsheets

Inevitably you will need to output reports from frangipanni into a spreadsheet. You can use the `-spacer` option to specify the character(s) to use for indentation and before the counts. So with the file list example from above and this command

```
sudo find /etc -maxdepth 3 | tail -9 | frangipanni -no-fold -counts -indent 1 -spacer $'\t'
```

You will have a tab-separated output which can be imported to your spreadsheet.


||||
|---|---|---|
|etc|9
|bluetooth|6
||rfcomm.conf.dpkg-remove|1
||serial.conf.dpkg-remove|1
||input.conf|1
||audio.conf.dpkg-remove|1
||network.conf|1
||main.conf|1
|fish/completions/task.fish|3


## Output for Markdown

To use the output with markdown or other text-based tools, sepecify the `-separator` option. This can be used by tools like `sed` to convert the leading separator into the markup required. example to get a leading minus sign 
for an un-numbered Markdown list, use `sed` to 

```
sudo find /etc -maxdepth 3 | tail -9 | frangipanni -separators | sed 's;/; - ;'
```

Which results in an indented bullet list:

>- etc
>    - bluetooth
>        - rfcomm.conf.dpkg-remove
>        - serial.conf.dpkg-remove
>        - input.conf
>        - audio.conf.dpkg-remove
>        - network.conf
>        - main.conf
>    - fish/completions/task.fish











