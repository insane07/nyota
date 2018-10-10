package requestinterceptor

import (
	"nyota/backend/model"
	"net/http"
	"time"
)

/*
Below code is added refering to following github code:
   https://github.com/zenazn/goji/blob/master/web/middleware/nocache.go
*/
// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

// AddNoCacheHeader Adds no cache header to all requests...
func AddNoCacheHeader() Interceptor {

	// Create a new Middleware
	return func(f PrizmHandler) PrizmHandler {

		// Define the http.HandlerFunc
		return func(s *model.SessionContext, w http.ResponseWriter, r *http.Request) {

			// Delete any ETag headers that may have been set
			for _, v := range etagHeaders {
				if r.Header.Get(v) != "" {
					r.Header.Del(v)
				}
			}

			// Set our NoCache headers
			for k, v := range noCacheHeaders {
				w.Header().Set(k, v)
			}

			// Call the next middleware/handler in chain
			f(s, w, r)
		}
	}
}
