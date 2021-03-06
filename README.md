# frangipanni
Program to convert lines of text into beautiful tree structures.

<img src="frangipanni.jpg" alt="A Tree" width="200" align="right">

The program reads each line on the standard input in turn. It breaks each line into tokens, then adds the sequence of tokens into a tree structure. Lines with the same leading tokens are placed in the same branch of the tree. The tree is printed as indented lines or JSON format. Alternatively the tree can be passed to a user-provided Lua script which can produce any output format.

Options control where the line is broken into tokens, and how it is analysed and output.

## Basic Operation

Here is a simple example. Given this  command `sudo find /etc -maxdepth 3 | tail -9 `, 

We get this data:

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

When we pipe this into the `frangipanni` program :
```
sudo find /etc -maxdepth 3 | tail -9 | frangipanni
```
we see this output:

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
By default, it reads each line and splits them into tokens when it finds a non-alphanumeric character. 

In this next example we're processing a list of files produced by `find` so we only want to break on directories. So we can specify `-breaks /`. 

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

# Download and Installation

## Pre-Compiled Binaries
You can download executables of `frangipanni` from the Github repository in the  [Releases area.](https://github.com/birchb1024/frangipanni/releases) You will find archive files containing the binary and Lua files:

```
frangipanni_darwin_amd64.zip
frangipanni_freebsd_amd64.zip
frangipanni_js_wasm.zip
frangipanni_linux_386.tgz
frangipanni_linux_amd64.tgz
frangipanni_linux_arm64.tgz
frangipanni_netbsd_amd64.zip
frangipanni_openbsd_amd64.zip
frangipanni_windows_386.zip
frangipanni_windows_amd64.zip
```

Download the file for your operating system and hardware, then decompress the archive in a directory. Ensure the binary file is executable, and you're ready to go. Your directory should look something like this, depending on operating system:
```
.
├── [-rwxrwxr-x]  frangipanni_linux_arm64
├── [-rw-rw-r--]  json.lua
├── [-rw-rw-r--]  markdown.lua
└── [-rw-rw-r--]  xml.lua
```

I have tested `frangipanni_linux_amd64`, the others are output from the Go cross-compiler and are provided 'as-is'. Please send a Pull Request if you find an issue. 

## Building From Source Code

If there is no pre-compiled binary for your platform, you can compile from the source. First you need [the 'Go' compiler](https://golang.org/doc/install), version 1.16.5 or greater. After cloning [the frangipanni repository](https://github.com/birchb1024/frangipanni) with git, it suffices to run `GO111MODULE=on go build` and the executable will be built as `frangipanni`. You can run the regression test suite with `test/confidence.sh`, but first install `jp` from [github.com/jmespath/jp](https://github.com/jmespath/jp). 

Read all about `Go` here: [golang.org](https://golang.org/)  

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
    	Slice line after each character.
  -counts
    	Print number of matches at the end of the line.
  -depth int
    	Maximum tree depth to print. (default 2147483647)
  -down
    	Sort branches in descending order. (default ascending)
  -format string
    	Format of output: indent|json (default "indent")
  -indent int
    	Number of spaces to indent per level. (default 4)
  -level int
    	Analyse down to this level (positive integer). (default 2147483647)
  -lua string
    	Lua Script to run
  -no-fold
    	Don't fold into one line.
  -separators
    	Print leading separators.
  -skip int
    	Number of leading fields to skip.
  -sort string
    	Sort by input|alpha|counts. Sort the branches either by input order,
        or via alphanumeric character ordering, or the branch frequency count.
        (default "input")
  -spacer string
    	Characters to indent lines with. (default " ")
  -version
    	Print frangipanni's version number and exit.
```

# Examples

## Log files

Operating systems and applications generate copious ASCII log files, frangipanni can help make sense of them. 
Here is a portion of a Linux system log file, `/var/log/syslog`: 

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

default frangipanni output is:

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

With the `-skip 5` option we can ignore the date and time at the beginning of each line. The output is

```
localhost
    systemd
        Removed slice User Slice of root
        Stopping User Slice of root
        Starting Docker Cleanup
        Started Docker Cleanup
    dockerd-current: time="2020-05-10T04:00:00
        629849861+10:00" level=debug msg="Calling GET /_ping
        629948000+10:00" level=debug msg="Unable to determine container for
        630103455+10:00" level=debug msg="{Action=_ping, LoginUID=12345678, PID=21075
        630684502+10:00" level=debug msg="Calling GET /v1.26/containers/json?all=1&filters=%7B%22status%22%3A%7B%22dead%22%3Atrue%7D%7D
        630704513+10:00" level=debug msg="Unable to determine container for containers
        630735545+10:00" level=debug msg="{Action=json, LoginUID=12345678, PID=21075
```

We can use the `-counts -sort counts -down` flags to list the most frequently occurring branches first. In this next example I 
* decompress all the historical log files (`zcat`)
* skip the month name and day-of-month `-skip 2`, this gives a report for each hour of the day
* just print the hours `depth 1 -no-fold`
* sort by the number of branches in descending order `-counts -sort counts -down`
 
```
$ zcat /var/log/syslog* | ./frangipanni -skip 2 -depth 1 -no-fold -counts -sort counts -down | head
10: 1247
02: 295
12: 106
06: 94
11: 91
00: 77
21: 76
19: 70
05: 69
13: 68
```
Clearly 10:00am is the busiest time for this machine. 

You can look for unusual events in security-related log files with the `-counts -sort counts` flags since less frequent items are pushed to the top. This command is quite interesting on my Linux system:

` zcat /var/log/auth* | ./frangipanni -skip 1 -depth 7 -counts -sort counts -breaks '0123456789- []():/'  | less`

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
  ./frangipanni  -breaks ' ' -sort alpha -format json -no-fold | \
  jq '."jupiter"'
```

```json
{
  "defaultAccount": "alice1",
  "hasUser": [
    "alice1",
    "birchb1",
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

## Lua Examples

### Accessing the Tree from Lua and output to JSON

First, we are going tell frangipanni to output via a Lua program called `json.lua`, and we will format the json with the 'jp' program.

```
$ <test/fixtures/simplechars.txt frangipanni -lua json.lua | jp @
```

The Lua script uses the `github.com/layeh/gopher-json` module which is imported in the Lua. The data
is made available in the variable `frangipanni` which has a table for each node, with fields

* `depth` - in the tree starting from 0
* `lineNumber` - the token was first detected
* `numMatched` - the number of times the token was seen
* `sep` - separation characters preceding the token
* `text` - the token itself
* `children` - a table containing the child nodes 

```Lua
local json = require("json")
print(json.encode(frangipanni))
```
The output shows that all the fields of the parsed nodes are passed to Lua in a Table.
The root node is empty except for it's children. The Lua script is therafore able to use
the fields intelligently.

```json
{
  "depth": 0,
  "lineNumber": -1,
  "numMatched": 1,
  "sep": "",
  "text": "",
  "children": {
    "1.2": {
      "children": [],
      "depth": 1,
      "lineNumber": 8,
      "numMatched": 1,
      "sep": "",
      "text": "1.2"
    },
    "A": {
      "children": [],
      "depth": 1,
      "lineNumber": 1,
      "numMatched": 1,
      "sep": "",
      "text": "A"
    }
  }
}
```

### Accessing frangipanni command-line arguments from Lua

frangipanni hands the command arguments to the Lua interpreter in the variable `frangipanni_args`. This holds a table keyed on the argument switch and holding the current value in a string. Command-line arguments found after switches are collected in the `args` array in this table. This example Lua program, `args.lua`, prints the table in JSON format:
```Lua
local json = require("json")
print(json.encode(frangipanni_args))
```
The output is:
```JSON
$ ./frangipanni -breaks / -counts -depth 3 -level 20 -lua args.lua arg1 arg2 </dev/null | jp '@'
{
  "args": [
    "arg1",
    "arg2"
  ],
  "breaks": "/",
  "chars": "false",
  "counts": "true",
  "depth": "3",
  "down": "false",
  "format": "indent",
  "indent": "4",
  "level": "20",
  "lua": "args.lua",
  "no-fold": "false",
  "separators": "false",
  "skip": "0",
  "sort": "input",
  "spacer": " ",
  "version": "false"
}
```

### Markdown

This example shows recursive scanning of the tree, with output format controlled by
an environment variable. 

```Lua
function indent(n)
    for i=1, n do
        io.write("   ")
    end
end

function markdown(node, bullet)
    if node.lineNumber > 0 then  -- don't write a root note
        indent(node.depth -1)
        io.write(bullet)
        print(node.text)
    end
    for k, v in pairs(node.children) do
        markdown(v, bullet)
    end
end

markdown(frangipanni, os.getenv("NUMBERED_LIST") and "1. " or "* ")
```

With `./frangipanni -lua markdown.lua <test/fixtures/simplechars.txt` The output looks like this:

```
* Z
* 1.2
* A
* C
   * 2
   * D
* x.a
   * 1
   * 2
```
or with `NUMBERED_LIST=1 ./frangipanni -lua markdown.lua <test/fixtures/simplechars.txt`
```
1. Z
1. 1.2
1. A
1. C
   1. 2
   1. D
1. x.a
   1. 1
   1. 2
```
Which renders like this:
1. Z
1. 1.2
1. A
1. C
   1. 2
   1. D
1. x.a
   1. 1
   1. 2

### XML

The xml.lua script provided in the release outputs very basic XML format which might suit simple inputs. Example

`find /proc | head -5  | ./frangipanni -breaks / -no-fold -lua xml.lua`

```XML
<root count="1" sep="">
   <proc count="5" sep="/">
      <fb count="1" sep="/"/>
      <fs count="3" sep="/">
         <aufs count="2" sep="/">
            <plink_maint count="1" sep="/"/>
         </aufs>
      </fs>
   </proc>
</root>
```






