# Caddy Darkvisitors Module

A super simple [Caddy](https://caddyserver.com/) module to send visit events to the [Dark Visitors API](https://darkvisitors.com/docs/analytics).

## Building

To compile this Caddy module, follow the instructions from [Build from Source](https://caddyserver.com/docs/build) and import the `github.com/polykernel/caddy-darkvisitors` module.

## Configuration

A basic Caddyfile configuration is provided below:

```Caddyfile
darkvisitors {
  # endpoint https://api.darkvisitors.com/visits
  access_token {env.DV_ACCESS_TOKEN}
}
```

By default, the `darkvisitors` directive is ordered after [`route`](https://caddyserver.com/docs/caddyfile/directives#directive-route) in the Caddyfile. If this order does not fit your needs, you can change the order using the global [`order`](https://caddyserver.com/docs/caddyfile/directives#directive-order) directive. For example:

```Caddyfile
{
  order darkvisitors before handle
}
```

Global placeholders are supported in the `darkvisitors` block.