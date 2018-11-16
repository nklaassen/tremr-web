package api

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

func update(w http.ResponseWriter, r *http.Request) {
	log.Println("running git pull...")
	cmd := exec.Command("git", "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		log.Print(string(output))
		return
	}

	log.Println("running go build...")
	cmd = exec.Command("go", "build")
	output, err = cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(output)
		log.Print(string(output))
		return
	}

	// execve the updated server. This should not return.
	err = syscall.Exec(os.Args[0], os.Args, os.Environ())
	if err != nil {
		http.Error(w, "failed to restart server", http.StatusInternalServerError)
	}
}
