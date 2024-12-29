# Caddy Darkvisitors Module

[![Go Reference](https://pkg.go.dev/badge/github.com/polykernel/caddy-darkvisitors.svg)](https://pkg.go.dev/github.com/polykernel/caddy-darkvisitors)

A super simple [Caddy](https://caddyserver.com/) module for interacting with the [Dark Visitors API](https://darkvisitors.com/docs/analytics).

## Building

To compile this Caddy module, follow the instructions from [Build from Source](https://caddyserver.com/docs/build) and import the `github.com/polykernel/caddy-darkvisitors` module.

## Configuration

### Syntax

```Caddyfile
darkvisitors {
  access_token <token>
  robots_txt {
    agent_types <types...>
    disallow <path>
  }
}
```

- **access_token** sets the OAuth authorization token used to communicate with the Dark Visitors API. Global placeholders are supported in the argument.
- **robots_txt** enables generation of robots.txt derived from agent analytics data using the Dark Visitors API.
  - **agent_types** specifies a list of [agent types](https://darkvisitors.com/agents) to be blocked by the generated robots.txt. The special token "\*" is supported as an argument which resolves to all documented agent types. Note: when "\*" is passed, there must be no further arguments.
  - **disallow** specifies the path to disallow for the specified agent types. Default: `/`.

If the `robots_txt` block is configured, then the special variable `http.vars.dv_robots_txt` in the HTTP request context will be set to the raw content of the robots.txt returned by the Dark Visitors API. Note: the robots.txt query is performed once during the provision phase of the module lifecycle and cached thereafter.

By default, the `darkvisitors` directive is ordered before [`header`](https://caddyserver.com/docs/caddyfile/directives#directive-header) in the Caddyfile. This ensures that the raw request content (sensitive data such as cookies are still stripped) is used to build a visit event. If this order does not fit your needs, you can change the order using the global [`order`](https://caddyserver.com/docs/caddyfile/directives#directive-order) directive. For example:

```Caddyfile
{
  order darkvisitors before handle
}
```

### Example

A basic Caddyfile configuration is provided below:

```Caddyfile
darkvisitors {
  access_token {env.DV_ACCESS_TOKEN}
  robots_txt {
    agent_types "AI Assistant" "AI Data Scraper"
    disallow /
  }
}
```

## License

Copyright (c) 2024 polykernel

The source code in this repository is made avaliable under the [MIT](https://opensource.org/license/mit) license.