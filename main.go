package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

var createBuffer = `
(let ((buf (generate-new-buffer "{{.}}")))
	(buffer-name buf))`
var updateBuffer = `
(with-current-buffer {{.Name}}
  (save-excursion (goto-char (point-max))
  (insert "{{js .Data}}\n")))`
var setMode = `
(with-current-buffer {{.Name}}
  ({{.Data}}-mode))`

const BUFFER_NAME = "*stdin*"

type Chunk struct {
	Name string
	Data string
}

func main() {
	bufferName, err := CreateBuffer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create buffer: %v\n", err)
	}
	ch := make(chan string)
	done := make(chan bool)
	// start update buffer routine
	go UpdateBuffer(bufferName, ch, done)
	scanner := bufio.NewScanner(os.Stdin)
	i := 0
	for scanner.Scan() {
		i++
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	//fmt.Fprintf(os.Stderr, "Read %d lines\n", i)
	done <- true
	<-done

	if len(os.Args) > 1 {
		modeName := os.Args[1]
		if modeName != "" {
			SetMode(bufferName, modeName, done)
			<-done
		}
	}
	os.Exit(0)
}

func SetMode(bufferName, modeName string, done chan bool) {
	tmpl, err := template.New("").Parse(setMode)
	if err != nil {
		done <- true
		return
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, Chunk{bufferName, modeName})
	c := exec.Command("emacsclient", "-n", "-u", "-e", buf.String())
	c.CombinedOutput()
	done <- true
}

func CreateBuffer() (string, error) {
	tmpl, err := template.New("").Parse(createBuffer)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	tmpl.Execute(&buf, BUFFER_NAME)
	c := exec.Command("emacsclient", "-n", "-e", buf.String())
	if out, err := c.CombinedOutput(); err == nil {
		return fmt.Sprintf("%s", out), nil
	} else {
		return "", err
	}
}
func UpdateBuffer(bufferName string, ch chan string, done chan bool) {
	tmpl, err := template.New("").Parse(updateBuffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse template: %v\n", err)
	}
	var buf bytes.Buffer
	i := 0
	for {
		select {
		case line := <-ch:
			tmpl.Execute(&buf, Chunk{bufferName, line})
			c := exec.Command("emacsclient", "-n", "-u", "-e", buf.String())
			if _, err := c.CombinedOutput(); err != nil {
				// probably the buffer was closed
				os.Exit(1)
			}
			i++
			buf.Reset()
		case <-done:
			//fmt.Fprintf(os.Stderr, "Sent %d lines\n", i)
			done <- true
		}
	}
}
