# frangipanni
Program to convert lines of text into a tree structure.

<img src="frangipanni.jpg" alt="A Tree" width="200">

# Usage

```
cat file.txt | go run frangipanni.go
```

or from binary 

```
cat file.txt | frangipanni
```

# Example

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

Output is:

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