package echo

import (
	"WebApplication/env"
	"fmt"
)

func echoAddress() string {
	return fmt.Sprintf("%s:%s", env.Get(env.EchoAddress, "127.0.0.1"), env.Get(env.EchoPort, "8080"))
}
func swaggerAddress() string {
	return fmt.Sprintf("%s:%s",
		env.Get(env.EchoSwaggerAddress, "nil"),
		env.Get(env.EchoSwaggerPort, "nil"))
}
