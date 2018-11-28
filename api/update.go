package api

import (
	"log"
	"net/http"
	"os/exec"
)

func update(reboot chan struct{}) HttpErrorHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("running git pull...")
		cmd := exec.Command("git", "pull")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Print(string(output))
			return err
		}

		log.Println("running go build...")
		cmd = exec.Command("go", "build")
		output, err = cmd.CombinedOutput()
		if err != nil {
			log.Print(string(output))
			return err
		}

		log.Println("running go test...")
		cmd = exec.Command("go", "test")
		output, err = cmd.CombinedOutput()
		if err != nil {
			log.Print(string(output))
			return err
		}

		close(reboot)
		return nil
	}
}
