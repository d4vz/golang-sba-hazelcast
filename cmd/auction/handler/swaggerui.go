package handler

import (
    "net/http"
)

// SwaggerUI serves a minimal Swagger UI using CDN and the generated doc JSON at /swagger/doc.json.
func SwaggerUI(mux *http.ServeMux) {
    mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/swagger/" {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        _, _ = w.Write([]byte(`<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({ url: '/swagger/doc.json', dom_id: '#swagger-ui' });
      };
    </script>
  </body>
</html>`))
    })
}


