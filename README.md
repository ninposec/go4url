# GoURL

Reads URLs from STID or url input flag, extracts all endpoints, both absolute and relative URL endpoints.

A common use case is extracting URL endpoints from Javascript files.


### Installation

### Usage

```bash
echo "https://x.y.z/js/main.js" | go4url
cat jsurls.txt | go4url
```

Can be combined with other tools, e.g. GetJS:

```bash
cat alive.txt | getJS --complete --insecure | go4url
```

Go4url help/usage:

```bash
go run go4url.go -h

		
██████   ██████   ██   ██ ██    ██ ██████  ██      
██       ██    ██ ██   ██ ██    ██ ██   ██ ██      
██   ███ ██    ██ ███████ ██    ██ ██████  ██      
██    ██ ██    ██      ██ ██    ██ ██   ██ ██      
 ██████   ██████       ██  ██████  ██   ██ ███████ 
												   
																   
			
			
go4url v.0.1
Author: ninposec

Usage: cat urls.txt | go4url
Extract URLs from Input e.g. JS Files

Options:
  -c int
    	Concurrency level (default 1)
  -nd
    	Ignore and suppress error messages
  -urls string
    	File containing URLs


```

