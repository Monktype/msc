package callback

import (
	"fmt"
	"net/http"
	"net/url"
)

const redirectHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>OAuth2 Redirect</title>
    <script type="text/javascript">
        window.onload = function() {
            // Get the fragment identifier
            const fragment = window.location.hash.substring(1);
            const params = new URLSearchParams(fragment);

            // Convert fragment parameters to query string
            const queryString = Array.from(params.entries()).map(([key, value]) => ` + "`${encodeURIComponent(key)}=${encodeURIComponent(value)}`" + `).join('&');

            // Redirect to the server with the query string
            window.location.href = ` + "`http://localhost:3024/redirect?${queryString}`" + `;
        };
    </script>
</head>
<body>
    <h1>Processing OAuth2 Callback...</h1>
</body>
</html>`

// StartCallbackServer starts an HTTP server to handle OAuth2 callbacks.
func StartCallbackServer(state string, responseChan chan<- CallbackResponse, shutdownChan <-chan bool) error {
	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Invalid query parameters", http.StatusBadRequest)
			return
		}

		// Combine query and fragment parameters
		params := make(url.Values)
		for k, v := range queryParams {
			params[k] = v
		}

		// Check if the request contains a fragment identifier
		if len(params) == 0 {
			// Serve the HTML page with embedded JavaScript
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintln(w, redirectHTML)
			return
		}

		receivedState := params.Get("state")
		if receivedState != state {
			http.Error(w, fmt.Sprintf("Invalid state %s", receivedState), http.StatusForbidden)
			return
		}

		if errParam := params.Get("error"); errParam != "" {
			responseChan <- CallbackResponse{
				Error: errParam + ": " + params.Get("error_description"),
			}
		} else {
			accessToken := params.Get("access_token")
			if accessToken == "" { // if it's a code flow and not a token flow.
				accessToken = params.Get("code")
			}
			responseChan <- CallbackResponse{
				AccessToken: accessToken,
			}
		}

		// Send a simple response back to the client
		fmt.Fprintln(w, "Callback received successfully.\nYou may close this tab or window.")
	})

	server := &http.Server{Addr: ":3024"}

	go func() {
		<-shutdownChan
		server.Shutdown(nil)
	}()

	return server.ListenAndServe()
}
