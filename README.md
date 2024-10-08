# caddyja3

A caddy plugin to get the JA3 fingerprint from requests as a header

## Building with xcaddy

```shell
xcaddy build \
  --with github.com/antoniomika/caddyja3
```

## Sample Caddyfile

Note that this enforces HTTPS (TLS).\
You can add a http_redirect to automatically redirect `http` -> `https` like shown below.

```
{
	debug
	log {
		output stdout
		format console
	}
	order ja3 before respond
	servers {
		listener_wrappers {
			http_redirect
			ja3
			tls
		}
	}
}

localhost:2020 {
	ja3
	tls internal
	respond <<EOF
          JA3 Hash: {header.ja3-hash}
          JA3 String: {header.ja3-string}
          JA3 SNI: {header.ja3-sni}
          EOF 200
}
```

## Acknowledgements

This repository is based on the code from [caddy-ja3](https://github.com/rushiiMachine/caddy-ja3)