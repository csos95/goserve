# GoServe
A very basic file server

Install with `go get -u gitlab.csos95.com/csos95/goserve`

flag|description|default
-|-|-
d|directory to serve|.
e|comma separated list of files to exclude|
m|max number of ports to try|10
o|open the default browser|false
p|port to use (tries random ports [8000,9000) if in use)|8080
