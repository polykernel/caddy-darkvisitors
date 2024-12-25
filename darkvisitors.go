package caddydarkvisitors

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// The default address for the Dark Visitors agent analytics API endpoint.
const DefaultEndpoint = "https://api.darkvisitors.com/visits"

func init() {
	caddy.RegisterModule(Darkvisitors{})
	httpcaddyfile.RegisterHandlerDirective("darkvisitors", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("darkvisitors", "after", "route")
}

// Darkvisitors is a middleware which implements a HTTP handler that sends
// HTTP request information as visit events to the Dark Visitors API.
//
// Its API is still experimental and may be subject to change.
type Darkvisitors struct {
	// The address of the Dark Visitors agent analytics API endpoint.
	//
	// Defaults to `https://api.darkvisitors.com/visits` if unspecified.
	Endpoint string `json:"endpoint,omitempty"`

	// The access token used to authenticate to the Dark Visitors agent
	// analytics API endpoint.
	AccessToken string `json:"access_token"`

	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Darkvisitors) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.darkvisitors",
		New: func() caddy.Module { return new(Darkvisitors) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Darkvisitors) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()

	if m.Endpoint == "" {
		m.Endpoint = DefaultEndpoint
	} else {
		m.Endpoint = repl.ReplaceAll(m.Endpoint, "")
	}
	m.AccessToken = repl.ReplaceAll(m.AccessToken, "")
	m.logger = ctx.Logger()

	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (m Darkvisitors) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// run the next handler
	err := next.ServeHTTP(w, r)
	if err != nil {
		return err
	}

	go func() {
		sanitizedHeaders := r.Header.Clone()
		sanitizedHeaders.Del("Cookie")

		visitEvent := map[string]interface{}{
			"request_path":    r.URL.Path,
			"request_method":  r.Method,
			"request_headers": sanitizedHeaders,
		}

		body, err := json.Marshal(visitEvent)
		if err != nil {
			m.logger.Error("Error marshaling visitor event", zap.Error(err))
			return
		}

		m.logger.Debug("Visit event payload constructed", zap.ByteString("payload", body))

		client := &http.Client{}
		req, err := http.NewRequest("POST", m.Endpoint, bytes.NewBuffer(body))
		if err != nil {
			m.logger.Error("Error creating request", zap.Error(err))
			return
		}

		req.Header.Set("Authorization", "Bearer "+m.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			m.logger.Warn("Error sending visitor event", zap.Error(err))
		} else {
			m.logger.Debug("Visitor event sent", zap.Int("status", resp.StatusCode))
		}
		defer resp.Body.Close()
	}()

	return nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *Darkvisitors) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	for d.NextBlock(0) {
		switch d.Val() {
		case "endpoint":
			if !d.NextArg() {
				return d.ArgErr()
			}
			m.Endpoint = d.Val()
		case "access_token":
			if !d.NextArg() {
				return d.ArgErr()
			}
			m.AccessToken = d.Val()
		default:
			return d.Errf("unrecognized subdirective '%s'", d.Val())
		}
	}

	if d.NextArg() {
		return d.Errf("unexpected argument '%s'", d.Val())
	}

	if m.AccessToken == "" {
		return d.Err("missing access token")
	}

	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Darkvisitors middleware.
//
// Syntax:
//
//	darkvisitors {
//	  access_token <token>
//	}
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Darkvisitors
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Darkvisitors)(nil)
	_ caddyhttp.MiddlewareHandler = (*Darkvisitors)(nil)
	_ caddyfile.Unmarshaler       = (*Darkvisitors)(nil)
)
