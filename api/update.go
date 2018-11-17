package api

import (
	"log"
	"net/http"
	"os/exec"
)

func update(reboot chan struct{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("running git pull...")
		cmd := exec.Command("git", "pull")
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, "git pull failed", http.StatusInternalServerError)
			w.Write(output)
			log.Print(string(output))
			return
		}

		log.Println("running go build...")
		cmd = exec.Command("go", "build")
		output, err = cmd.CombinedOutput()
		if err != nil {
			http.Error(w, "go build failed", http.StatusInternalServerError)
			w.Write(output)
			log.Print(string(output))
			return
		}

		close(reboot)
	}
}
