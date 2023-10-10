package setup

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func Run(port string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	services, err := Services()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	controllers := Controllers(services)

	r := Routes(controllers)

	fmt.Println("Starting app on port: %v", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
