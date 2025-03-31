package handlers

import (
	"net"
	"net/http"
)

// WithIPCheck is a middleware that checks if the request's IP address belongs to a trusted subnet.
// It extracts the IP address from the X-Real-IP header, parses it and checks if it is within the
// trusted subnet. If the IP address is trusted, it passes the request to the next handler in the
// chain. If not, it logs the information and responds with a Forbidden status.
func (h Handlers) WithIPCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from the X-Real-IP header
		ipStr := r.Header.Get("X-Real-IP")
		if ipStr == "" {
			// If there's no IP address, pass the request to the next handler
			next.ServeHTTP(w, r)
			return
		}

		// Parse the IP address
		ip := net.ParseIP(ipStr)

		// Check if the IP address is within the trusted subnet
		if !h.config.TrustedSubnetCIDR.Contains(ip) {
			// Log the information and respond with a Forbidden status if not
			h.lg.Sugar.Infow("IP is not in trusted subnet", "ip", ipStr)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// If the IP address is trusted, pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
