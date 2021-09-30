# xrealip

A [Traefik](https://traefik.io) middleware plugins developed using the [Go language](https://golang.org).

The goal is to match the behavior of `X-Real-Ip` of nginx.

## Usage

### Configuration

For each plugin, the Traefik static configuration must define the module name (as is usual for Go packages).

The following declaration (given here in YAML) defines a plugin:

```yaml
# Static configuration
pilot:
  token: xxxxx

experimental:
  plugins:
    xrealip:
      moduleName: github.com/tommoulard/xrealip
      version: v0.0.1
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - xrealip

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000

  middlewares:
    xrealip:
      plugin:
        xrealip:
          from:
             - 127.0.0.1
          header: X-Forwarded-For
          resursive: true
```

The configuration matches the one given by nginx itself:

 - `from` is equivalent as `set_real_ip_from`
 - `header` is equivalent as `real_ip_header`
 - `recursive` is equivalent as `real_ip_recursive`

You can find a more in depth configuration description on the [nginx doc](https://nginx.org/en/docs/http/ngx_http_realip_module.html) itself.

[![Build Status](https://github.com/tommoulard/xrealip/workflows/Main/badge.svg?branch=master)](https://github.com/tommoulard/xrealip/actions)
