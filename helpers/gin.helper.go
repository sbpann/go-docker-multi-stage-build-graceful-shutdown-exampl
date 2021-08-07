
package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func HttpServerFromGinEngine(r *gin.Engine, addr ...string)  *http.Server{
	return &http.Server{
		Addr:    resolveAddress(addr...),
		Handler: r,
	}
}
func resolveAddress(addr ...string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
		debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		debugPrint("PORT is programmatically set to Using port \"%s\"", addr[0])
		return addr[0]
	default:
		panic("too many parameters")
	}
}

func debugPrint(format string, values ...interface{}) {
	if gin.IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] "+format, values...)
	}
}
