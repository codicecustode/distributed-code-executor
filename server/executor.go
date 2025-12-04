package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"fmt"

	"distributed-code-executor/shared"
)

func executeCode(ctx context.Context, req *shared.CodeRequest) *shared.CodeResponse {

	startTime := time.Now()

	tempDir, err := os.MkdirTemp("", "code-exec-*")
	if err != nil {
		log.Fatal("Failed to create temp dir:", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after execution

	// Prepare execution environment
	execPath, args, err := prepareExecution(req, tempDir)
	if err != nil {
		return createErrorResponse(req.ID, "Failed to prepare execution", err)
	}

	output, err := runRemoteCodeWithLimit(ctx, execPath, args, tempDir)

	duration := time.Since(startTime).Milliseconds()

	resp := &shared.CodeResponse{
		ID:        req.ID,
		Output:    output,
		ExitCode:  -1,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		resp.Error = err.Error()
	}

	return resp

}

func prepareExecution(req *shared.CodeRequest, tempDir string) (string, []string, error) {

	switch req.Language {
	case "javascript":
		return prepareJavaScriptExecution(req, tempDir)
	}

	return "", []string{}, nil

}

func prepareJavaScriptExecution(req *shared.CodeRequest, tempDir string) (string, []string, error) {
	mainFile := filepath.Join(tempDir, "main.js")

	err := os.WriteFile(mainFile, []byte(req.Code), 0644)
	if err != nil {
		return "", nil, err
	}
	return "node", []string{mainFile}, nil
}

func runRemoteCodeWithLimit(ctx context.Context, execPath string, args []string, tempDir string) (string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cmdCtx, execPath, args...)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func createErrorResponse(id, msg string, err error) *shared.CodeResponse {
    errorMsg := msg
    if err != nil {
        errorMsg = fmt.Sprintf("%s: %v", msg, err)
    }
    
    return &shared.CodeResponse{
        ID:        id,
        Error:     errorMsg,
        ExitCode:  -1,
        Timestamp: time.Now(),
    }
}