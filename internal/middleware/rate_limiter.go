package middleware

import (
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/handlers"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	visitors   map[string]*rate.Limiter
	mutex      sync.Mutex
	rate       rate.Limit
	burst      int
	timeToLive time.Duration
}

func newIPLimiter(
	rateLimit rate.Limit,
	burst int,
	timeToLive time.Duration,
) *ipLimiter {
	limiter := &ipLimiter{
		visitors:   make(map[string]*rate.Limiter),
		rate:       rateLimit,
		burst:      burst,
		timeToLive: timeToLive,
	}

	go func() {
		ticker := time.NewTicker(timeToLive)
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return limiter
}

func (ipLimiter *ipLimiter) get(ip string) *rate.Limiter {
	ipLimiter.mutex.Lock()
	defer ipLimiter.mutex.Unlock()

	if limiter, ok := ipLimiter.visitors[ip]; ok {
		return limiter
	}

	limiter := rate.NewLimiter(ipLimiter.rate, ipLimiter.burst)
	ipLimiter.visitors[ip] = limiter
	return limiter
}

func (ipLimiter *ipLimiter) cleanup() {
	ipLimiter.mutex.Lock()
	defer ipLimiter.mutex.Unlock()
	for ip, limiter := range ipLimiter.visitors {
		if limiter.Allow() {
			delete(ipLimiter.visitors, ip)
		}
	}
}

var (
	globalLimiter = newIPLimiter(
		rate.Every(time.Minute/100),
		100,
		10*time.Minute,
	)
)

func GlobalRateLimiter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			ip, _, err := net.SplitHostPort(request.RemoteAddr)
			if err != nil {
				handlers.WriteError(writer, apierrors.NewAPIError(
					"Invalid IP",
					apierrors.ErrBadRequest,
				), http.StatusBadRequest)
				return
			}

			if !globalLimiter.get(ip).Allow() {
				handlers.WriteError(writer, apierrors.NewAPIError(
					"Too many request in a short period",
					apierrors.ErrTooManyRequests,
				), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func RateLimitMiddleware(
	rateLimit rate.Limit,
	burst int,
	timeToLive time.Duration,
) func(http.Handler) http.Handler {
	limiter := newIPLimiter(rateLimit, burst, timeToLive)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			ip, _, err := net.SplitHostPort(request.RemoteAddr)
			if err != nil {
				handlers.WriteError(writer, apierrors.NewAPIError(
					"Invalid IP",
					apierrors.ErrBadRequest,
				), http.StatusBadRequest)
				return
			}

			if !limiter.get(ip).Allow() {
				handlers.WriteError(writer, apierrors.NewAPIError(
					"Too many request in a short period",
					apierrors.ErrTooManyRequests,
				), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
