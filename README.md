# CSTV-Go (Refactored)

## Overview
CSTV-Go is a WebSocket-based server using Fiber for handling real-time messaging.

## Features
- WebSocket support with client registration
- Structured logging with Logrus
- Modular design with configuration and services

## Installation
```sh
git clone https://github.com/oldmagic/cstv-go.git
cd cstv-go
go mod tidy


## About
This is [GOTV+](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Broadcast) broadcast server interface for Go(Fiber and Gin).  
  
GOTV+ is an extension of GOTV where you use HTTP(S) to distribute instead of connecting to a regular GOTV. This makes it easy to serve many more clients around the world with high quality GOTV as you can distribute the content with CDN's.  
Using `tv_broadcast` cvars you will enable GOTV+ on your CS:GO Server which will send fragmented data to the GOTV+ ingest (this application) which then serves them to clients which connects to it. The viewer will then watch the feed the same way you would when connecting directly to a GOTV instance.  
  
## Usage
- Launch GOTV+ Server(Fiber example)  
Directly using Go: `go run ./examples/inmemory/cmd/main.go -port 8080 -auth gopher`  
Using precompiled binary: `./gotv_plus -port :8080 -auth gopher`  

- Enable GOTV broadcast in CS:GO Server by adding this to server.cfg  
`MATCH_ID` can be what ever you want!  
```
tv_enable 1
tv_broadcast_url "http://<IP-ADDRESS>:8080/gotv"
tv_broadcast_origin_auth "gopher"
tv_broadcast 1
```

- Connect to GOTV+ from CS:GO Client  
In console: `playcast "http://<IP-ADDRESS>:8080/gotv/MATCH_ID"`  

### Tips
- You should configure `tv_broadcast_...` cvars **before** enable broadcast (`tv_broadcast 1`).  
And you may need to restart broadcast if you changed cvars. (`tv_broadcast 0;tv_broadcast 1`).
- `playcast` and `tv_broadcast_url` URL needs to be wrapped by double quotes (e.g. `playcast "http://<IP-ADDRESS>:8080/match/MATCH_ID"`) because of `http:` colon(`:`) is recognized as command seprator for CS:GO.
- `playcast` and `tv_broadcast_url` does not need slash(`/`) in last character of URL.  
For example `playcast "http://<IP-ADDRESS>:8080/match/id/MATCH_ID"` works properly,  
while `playcast "http://<IP-ADDRESS>:8080/match/id/MATCH_ID/"` **WILL NOT** work properly.  
Because requested URL will be ending with `//sync` (Double-slash).

## Hidden Options
There are several hidden options that are not documented. I haven't covered all of the options, but some of them are supported.   
- "F" (e.g. `playcast "http://<IP-ADDRESS>:8080/match/id/MATCH_ID" f500`) will play match from `500` fragment. CS:GO client will send request to `http://<IP-ADDRESS>:8080/match/id/MATCH_ID/sync?fragment=500`
- "A" (e.g. `playcast "http://<IP-ADDRESS>:8080/match/id/MATCH_ID" a`) will play match from `1` fragment. CS:GO client will send request to `http://<IP-ADDRESS>:8080/match/id/MATCH_ID/sync?fragment=1`
- "C" (e.g. `playcast "http://<IP-ADDRESS>:8080/match/id/MATCH_ID" c`) will play match from `1` fragment. CS:GO client **WILL NOT** send request to `http://<IP-ADDRESS>:8080/match/id/MATCH_ID/sync`, CS:GO client will get 1st fragment directly.
- "B" Option / Skips dem_stop control frame(?)

## Example with Public CDN
In this example we run the application with user `gotv`, and then use NGINX to proxy TCP/80 (HTTP) and TCP/443 (HTTPS) traffic to the application. 
We advice that you limit who can send POST requests (CS:GO servers external/internal IP address) directly to the service with local firewall (iptables, nftables etc), this is the reason why we limit in NGINX all requests to only GET. 

We use two page rules on Cloudflare. We bypass Cache for all requests to /sync as we want that to be served directly from the application, and then cache everything on rest of the URL's. 

```
gotv.example.com/*/sync
Cache Level: Bypass

gotv.example.com/*
Cache Level: Cache Everything
```

nginx config
```
upstream ingest {
        server <IP-ADDRESS>:8080; # Address to the gotv-plus-go application.
}

server {
	listen 80 default_server;
	listen [::]:80 default_server;

	server_name _;

	# Never cache /sync (required for Google CDN and others where you cant configure excluded URL's)
	location ~* \/sync$ {
		add_header Cache-Control "no-store";
		proxy_pass http://ingest;
	}

	# Only allow GET requests
	location / {
		limit_except GET {
			deny all;
		}
		proxy_pass http://ingest;
	}

}

server {
	listen 443 ssl;
	listen [::]:443 ssl;

	ssl_certificate /etc/ssl/cloudflare/cloudflare-cert.pem;
	ssl_certificate_key /etc/ssl/cloudflare/cloudflare-key.key;

	server_name gotv.example.com;

	# Never cache /sync (required for Google CDN and others where you cant configure excluded URL's)
	location ~* \/sync$ {
		add_header Cache-Control "no-store";
		proxy_pass http://ingest;
	}

	# Only allow GET requests
	location / {
		limit_except GET {
			deny all;
		}
		proxy_pass http://ingest;
	}
}
```

simple systemd service
```
[Unit]
Description=GOTV+ service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=gotv
WorkingDirectory=/path/to/gotv # Required as the application needs template
ExecStart=/path/to/gotv_plus -addr <IP-ADDRESS>:8080 -auth gopher

[Install]
WantedBy=multi-user.target
```
