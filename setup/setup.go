package setup

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/DeanWard/erugo/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)

type state int

const (
	enterUsername state = iota
	enterPassword
	confirmPassword
	confirm
	done
)

type model struct {
	state           state
	username        string
	password        string
	confirmPassword string
	input           string
	err             string
	db              *sql.DB
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m.handleEnter()
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		default:
			// Handle confirmation step separately
			if m.state == confirm {
				if msg.String() == "y" || msg.String() == "Y" {
					err := createUser(m.db, m.username, m.password)
					if err != nil {
						m.err = err.Error()
						m.state = enterUsername
					} else {
						m.state = done
						return m, tea.Quit
					}
				} else if msg.String() == "n" || msg.String() == "N" {
					// Reset input and restart setup
					m.username = ""
					m.password = ""
					m.confirmPassword = ""
					m.input = ""
					m.state = enterUsername
				}
				return m, nil
			}

			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	var output string

	switch m.state {
	case enterUsername:
		output = fmt.Sprintf("Enter username: %s\n", m.input)
	case enterPassword:
		output = fmt.Sprintf("Enter password: %s\n", hidePassword(m.input))
	case confirmPassword:
		output = fmt.Sprintf("Confirm password: %s\n", hidePassword(m.input))
	case confirm:
		output = fmt.Sprintf("Confirm create user: %s?\n(Y to confirm, N to restart)\n", m.username)
	case done:
		output = "\nUser created successfully!\n"
	}

	if m.err != "" {
		output += fmt.Sprintf("\nError: %s\n", m.err) // Show errors in the UI
	}

	return output
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case enterUsername:
		m.username = m.input
		m.input = ""
		m.state = enterPassword
	case enterPassword:
		m.password = m.input
		m.input = ""
		m.state = confirmPassword
	case confirmPassword:
		m.confirmPassword = m.input
		if m.confirmPassword != m.password {
			m.err = "Passwords do not match! Try again."
			m.state = enterPassword
		} else {
			m.input = ""
			m.state = confirm
		}
	case confirm:
		if m.input == "y" || m.input == "Y" {
			err := createUser(m.db, m.username, m.password)
			if err != nil {
				m.err = err.Error()
				m.state = enterUsername
			} else {
				m.state = done
				return m, tea.Quit
			}
		} else if m.input == "n" || m.input == "N" {
			// Reset and restart setup
			m.username = ""
			m.password = ""
			m.confirmPassword = ""
			m.input = ""
			m.state = enterUsername
		}
	}
	m.input = "" // Reset input after handling
	return m, nil
}

func createUser(db *sql.DB, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err) // Log error
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		log.Println("Error inserting user into database:", err) // Log error
		return err
	}

	log.Println("User created successfully:", username) // Log success
	//clear the screen
	fmt.Println("\033[H\033[2J")

	fmt.Println("\n")
	color.New(color.FgYellow).Println("█████████████████████████████████████████████████████████████████████████████████")
	fmt.Println("\n")
	return nil
}

func hidePassword(s string) string {
	return strings.Repeat("*", len(s))
}

func RunSetup(database *sql.DB) {
	//clear the screen
	fmt.Println("\033[H\033[2J")
	displayLogo()
	//show the user a message explaining what the setup is for
	fmt.Println("Looks like you're setting up erugo for the first time! Welcome to erugo!")
	fmt.Println("This setup will help you create a user so that you can start using erugo.")
	fmt.Println("\n")
	//FgYellow line
	color.New(color.FgYellow).Println("█████████████████████████████████████████████████████████████████████████████████")
	fmt.Println("\n")
	// Check DB connection before running setup
	if err := database.Ping(); err != nil {
		log.Fatal("Database connection error:", err)
	}

	p := tea.NewProgram(model{state: enterUsername, db: database})
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running setup:", err)
	}
}

func displayLogo() {
	// Read the file
	data := utils.GetLogo()
	// Convert to string and split into lines
	lines := strings.Split(string(data), "\n")

	// Define rainbow colors (ANSI escape codes)
	colors := []func(string, ...interface{}){
		func(format string, a ...interface{}) { color.New(color.FgYellow).Printf(format, a...) },
		func(format string, a ...interface{}) { color.New(color.FgYellow).Printf(format, a...) },
		func(format string, a ...interface{}) { color.New(color.FgYellow).Printf(format, a...) },
		func(format string, a ...interface{}) { color.New(color.FgYellow).Printf(format, a...) },
		func(format string, a ...interface{}) { color.New(color.FgYellow).Printf(format, a...) },
		func(format string, a ...interface{}) { color.New(color.FgYellow).Printf(format, a...) },
	}

	// Print each line in a different color
	for i, line := range lines {
		colors[i%len(colors)]("%s\n", line)
	}
}
