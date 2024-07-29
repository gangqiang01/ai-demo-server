package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

// LogFunc used to log stdout and stderr.
type LogFunc func(string)
type Process interface {
	// Set timeout for process to exit after being stopped.
	Timeout(time.Duration) Process

	// Set function called on stdout line.
	StdoutLogger(LogFunc) Process

	// Set function called on stderr line.
	StderrLogger(LogFunc) Process

	// Start process with context.
	Start(ctx context.Context) error

	// Stop process.
	Stop()
}

// process manages subprocesses.
type process struct {
	timeout      time.Duration
	cmd          *exec.Cmd
	stdoutLogger LogFunc
	stderrLogger LogFunc

	done chan struct{}
}

// NewProcessFunc is used for mocking.
type NewProcessFunc func(*exec.Cmd) Process

// NewProcess return process.
func NewProcess(cmd *exec.Cmd) Process {
	return process{
		timeout: 1000 * time.Millisecond,
		cmd:     cmd,
	}
}

func (p process) Timeout(timeout time.Duration) Process {
	p.timeout = timeout
	return p
}

func (p process) StdoutLogger(l LogFunc) Process {
	p.stdoutLogger = l
	return p
}

func (p process) StderrLogger(l LogFunc) Process {
	p.stderrLogger = l
	return p
}

func (p process) Start(ctx context.Context) error {
	if p.stdoutLogger != nil {
		if p.cmd.Stdout == nil {
			pipe, err := p.cmd.StdoutPipe()
			if err != nil {
				klog.Errorf("StdoutPipe: %v", err)
			}
			p.attachLogger(p.stdoutLogger, "stdout", pipe)
		}

	}
	if p.stderrLogger != nil {
		if p.cmd.Stderr == nil {
			pipe, err := p.cmd.StderrPipe()
			if err != nil {
				klog.Errorf("StderrPipe: %v", err)
			}
			p.attachLogger(p.stderrLogger, "stderr", pipe)
		}

	}

	if err := p.cmd.Start(); err != nil {
		klog.Errorf("Process start error: %v", err.Error())
		return err
	}

	if err := p.cmd.Wait(); err != nil {
		klog.Errorf("Process exit error: %v", err.Error())
		return err
	}
	p.Stop()
	return nil
}

func (p process) attachLogger(logFunc LogFunc, label string, pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				klog.Errorf("attachLogger failed with %s", debug.Stack())
			}
		}()
		for scanner.Scan() {
			msg := fmt.Sprintf("%v: %v", label, scanner.Text())
			logFunc(msg)
		}
	}()
}

// Note, can't use CommandContext to Stop process as it would
// kill the process before it has a chance to exit on its own.
func (p process) Stop() {
	klog.Infof("Stop process client")
	p.cmd.Process.Signal(os.Interrupt) //nolint:errcheck
	time.Sleep(2 * time.Second)
	p.cmd.Process.Kill()
}

type Cmd struct {
	Name    string
	Command func(...string) *exec.Cmd
}

// New returns FFMPEG.
func NewCmd(name string) *Cmd {
	command := func(args ...string) *exec.Cmd {
		klog.Infof("command: %s, args: %v", name, args)
		return exec.Command(name, args...)
	}
	return &Cmd{Command: command}
}

// ParseArgs slices arguments.
func ParseArgs(args string) []string {
	return strings.Split(strings.TrimSpace(args), " ")
}
