package ui

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/fatih/color"
)

// SpinnerFrames defines different animation styles for the spinner
var SpinnerFrames = map[string][]string{
	"dots": {"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
	"arrow": {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	"smooth": {"▰▱▱▱▱▱▱", "▰▰▱▱▱▱▱", "▰▰▰▱▱▱▱", "▰▰▰▰▱▱▱", "▰▰▰▰▰▱▱", "▰▰▰▰▰▰▱", "▰▰▰▰▰▰▰", "▱▰▰▰▰▰▰", "▱▱▰▰▰▰▰", "▱▱▱▰▰▰▰", "▱▱▱▱▰▰▰", "▱▱▱▱▱▰▰", "▱▱▱▱▱▱▰", "▱▱▱▱▱▱▱"},
	"bounce": {"⠁", "⠂", "⠄", "⠂"},
	"classic": {"|", "/", "-", "\\"},
}

// Spinner represents an animated terminal spinner
type Spinner struct {
	frames     []string
	message    string
	frameIndex int
	stopChan   chan struct{}
	done       bool
	doneMsg    string
	errorMsg   string
	hasError   bool
	color      *color.Color
	mu         sync.Mutex
	interval   time.Duration
}

// NewSpinner creates a new spinner with a default style
func NewSpinner(message string) *Spinner {
	return NewSpinnerWithStyle(message, "dots")
}

// NewSpinnerWithStyle creates a new spinner with a specific style
func NewSpinnerWithStyle(message string, style string) *Spinner {
	frames, ok := SpinnerFrames[style]
	if !ok {
		frames = SpinnerFrames["dots"] // Default to dots if style not found
	}
	
	return &Spinner{
		frames:   frames,
		message:  message,
		stopChan: make(chan struct{}),
		color:    color.New(color.FgHiCyan),
		interval: 100 * time.Millisecond,
	}
}

// SetColor sets the color of the spinner
func (s *Spinner) SetColor(c *color.Color) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.color = c
}

// SetInterval sets the frame change interval
func (s *Spinner) SetInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interval = interval
}

// Start starts the spinner animation
func (s *Spinner) Start() *Spinner {
	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.mu.Lock()
				if s.done {
					s.mu.Unlock()
					return
				}
				frame := s.frames[s.frameIndex%len(s.frames)]
				s.frameIndex++
				message := s.message
				s.mu.Unlock()
				
				// Clear the current line
				fmt.Print("\r\033[K")
				
				// Print the spinner frame and message
				if s.color != nil {
					s.color.Printf("%s %s", frame, message)
				} else {
					fmt.Printf("%s %s", frame, message)
				}
				
				time.Sleep(s.interval)
			}
		}
	}()
	
	return s
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	if s.done {
		s.mu.Unlock()
		return
	}
	s.done = true
	s.mu.Unlock()
	s.stopChan <- struct{}{}
	
	// Clear the current line
	fmt.Print("\r\033[K")
}

// Success stops the spinner and displays a success message
func (s *Spinner) Success(message string) {
	s.mu.Lock()
	s.doneMsg = message
	s.hasError = false
	s.mu.Unlock()
	
	s.Stop()
	successColor := color.New(color.FgHiGreen)
	successColor.Printf("✓ %s\n", message)
}

// Error stops the spinner and displays an error message
func (s *Spinner) Error(message string) {
	s.mu.Lock()
	s.errorMsg = message
	s.hasError = true
	s.mu.Unlock()
	
	s.Stop()
	errorColor := color.New(color.FgHiRed)
	errorColor.Printf("✗ %s\n", message)
}

// UpdateMessage updates the spinner message while it's running
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}
