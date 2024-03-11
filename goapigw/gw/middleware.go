package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(verifier TokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(&CustomError{
				Inner:      err,
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			}))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(&CustomError{
				Inner:      err,
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			}))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(&CustomError{
				Inner:      err,
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			}))
			return
		}

		accessToken := fields[1]
		payload, err := verifier.Verify(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(&CustomError{
				Inner:      err,
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			}))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

/*func logMiddleware(logDB *sqlx.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now().UTC()
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		c.Next()

		latency := time.Since(t).Milliseconds()
		statusCode := c.Writer.Status()
		url := c.Request.RequestURI
		remoteAddr := c.Request.RemoteAddr
		remoteSplit := strings.Split(remoteAddr, ":")
		if len(remoteSplit) != 0 {
			remoteAddr = ""
		}
		remoteAddr = remoteSplit[0]

		method := c.Request.Method
		userAgent := c.Request.Header.Get("User-Agent")
		requestTime := t.Format("2006/01/02 15:04:05.000000")
		errs := c.Errors
		stacktrace := errs.String()
		errorResponse := blw.body.String()

		var userID string
		payload, exists := c.Get(authorizationPayloadKey)
		if exists {
			userID = payload.(*token.Payload).UserID
		}

		level := ""
		switch {
		case statusCode >= 500:
			// if logger.Fatal(...)
			// shutdown server
			//stacktrace := zap.Stack("stacktrace")
			//log.Println(stacktrace)
			level = "ERROR"
			logger.Error("",
				zap.String("request_time", requestTime),
				zap.Int("status_code", statusCode),
				zap.String("url", url),
				zap.String("method", method),
				zap.String("user_id", userID),
				zap.String("remote_address", remoteAddr),
				zap.String("user_agent", userAgent),
				zap.Int64("latency", latency),
				zap.String("error_response", errorResponse),
				//zap.Stack("stacktrace"),
			)
		case statusCode >= 400:
			level = "WARN"
			logger.Warn("",
				zap.String("request_time", requestTime),
				zap.Int("status_code", statusCode),
				zap.String("url", url),
				zap.String("method", method),
				zap.String("user_id", userID),
				zap.String("remote_address", remoteAddr),
				zap.String("user_agent", userAgent),
				zap.Int64("latency", latency),
				zap.String("error_response", errorResponse),
			)
		case statusCode >= 200:
			level = "INFO"
			errorResponse = ""
			logger.Info("",
				zap.String("request_time", requestTime),
				zap.Int("status_code", statusCode),
				zap.String("url", url),
				zap.String("method", method),
				zap.String("user_id", userID),
				zap.String("remote_address", remoteAddr),
				zap.String("user_agent", userAgent),
				zap.Int64("latency", latency),
			)
		}

		func() {
			query := `
				INSERT INTO access_log(request_time,
				                       status_code,
				                       url,
				                       method,
				                       user_id,
				                       remote_address,
				                       user_agent,
				                       latency,
				                       error_response,
				                       error,
				                       level)
				VALUES(?, ?, ?, ?, ?, INET_ATON(?), ?, ?, ?, ?, ?)`
			_, _ = logDB.Exec(query, requestTime, statusCode, url, method, userID, remoteAddr, userAgent, latency, errorResponse, stacktrace, level)
		}()

	}
}*/
