package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	dataDir  = "data"
	dataFile = filepath.Join(dataDir, "workouts.json")
)

// Workout represents one training entry

type Workout struct {
	Date     string  `json:"date"`     // ISO date: 2026-01-10
	Exercise string  `json:"exercise"` // e.g., "Squat, "Deadlift", "bench"
	Weight   float64 `json:"weight"`   // in kg
	Reps     int     `json:"reps"`
}

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "athletelog-cli",
	Short: "Track your workouts like a pro",
	Long: `Personal Training Log CLI
Built by a former athlete learning Rust, Go, Python, and Typescript in 2026.
Starting simple -> adding power later with Rust calculations, Python reports, and TS dashboard.`,
}

// addCmd - subcommand to add a workout
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new workout entry",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println("Usage: athletelog-cli add <date> <exercise> <weight> <reps>")
			fmt.Println("Example: athletelog-cli add 2026-01-10 Squat 100 5")
			os.Exit(1)
		}

		date := args[0]
		exercise := args[1]
		weightStr := args[2]
		repsStr := args[3]

		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			fmt.Printf("Invalid weight: %v\n", err)
			os.Exit(1)
		}

		reps, err := strconv.Atoi(repsStr)
		if err != nil {
			fmt.Printf("Invalid reps: %v\n", err)
			os.Exit(1)
		}

		// Validate date roughly
		if _, err := time.Parse("2006-01-02", date); err != nil {
			fmt.Println("Date must be in YYYY-MM-DD format")
			os.Exit(1)
		}

		workout := Workout{
			Date:     date,
			Exercise: exercise,
			Weight:   weight,
			Reps:     reps,
		}

		// Ensure data dir exists
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			fmt.Printf("Failed to create data directory: %v\n", err)
			os.Exit(1)
		}

		var workouts []Workout

		// Read existing data if file exists
		if _, err := os.Stat(dataFile); err == nil {
			data, err := os.ReadFile(dataFile)
			if err != nil {
				fmt.Printf("Failed to read data file: %v\n", err)
				os.Exit(1)
			}
			if len(data) > 0 {
				if err := json.Unmarshal(data, &workouts); err != nil {
					fmt.Printf("Failed to unmarshal JSON: %v\n", err)
					os.Exit(1)
				}
			}
		}

		// Append new workout
		workouts = append(workouts, workout)

		// Write back
		jsonData, err := json.MarshalIndent(workouts, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal JSON: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(dataFile, jsonData, 0644); err != nil {
			fmt.Printf("Failed to write data file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added: %s - %s @ %.1f lb x %d reps\n", date, exercise, weight, reps)
	},
}

// viewCmd - show all workouts
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View all workout entries",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(dataFile); os.IsNotExist(err) {
			fmt.Println("No workouts logged yet. Add some first!")
			return
		}

		data, err := os.ReadFile(dataFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		var workouts []Workout
		if err := json.Unmarshal(data, &workouts); err != nil {
			fmt.Printf("Error unmarshalling JSON: %v\n", err)
			return
		}

		if len(workouts) == 0 {
			fmt.Println("No workouts recorded yet.")
			return
		}

		fmt.Println("Your Training Log:")
		for i, w := range workouts {
			fmt.Printf("%3d | %s | %-12s | %.1f lb x %d reps\n", i+1, w.Date, w.Exercise, w.Weight, w.Reps)
		}

	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Calculate advanced stats (e.g. estimated 1RM) using Rust",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(dataFile); os.IsNotExist(err) {
			fmt.Println("No workouts yet. Add some first!")
			return
		}

		// Path to Rust binary
		projectRoot, _ := os.Getwd()
		rustBin := filepath.Join(projectRoot, "rust-stats", "stats", "target", "release", "stats")

		cmdExec := exec.Command(rustBin, dataFile)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			fmt.Printf("Failed to run Rust stats: %v\n", err)
			return
		}

		fmt.Println("\nRust-powered stats complete!")
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate progress report & charts using Python",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(dataFile); os.IsNotExist(err) {
			fmt.Println("No workouts yet. Add some first!")
			return
		}

		projectRoot, _ := os.Getwd()
		pyBin := filepath.Join(projectRoot, "python-report", "report.py")

		cmdExec := exec.Command("python3", pyBin, dataFile)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr

		if err := cmdExec.Run(); err != nil {
			fmt.Printf("Failed to run Python report: %v\n", err)
			return
		}

		fmt.Println("\nPython-powered report complete!")
	},
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch the interactive Typescript dashboard in your browser",
	Run: func(cmd *cobra.Command, args []string) {
		const port = ":3000"
		const dashboardPath = "/ts-dashboard/"
		const fullURL = "http://localhost" + port + dashboardPath

		// Start embedded static server in a goroutine
		go func() {
			rootDir, err := os.Getwd()
			if err != nil {
				fmt.Println("Failed to get working directory:", err)
				return
			}
			fmt.Printf("Serving files from: %s\n", rootDir)

			fs := http.FileServer(http.Dir(rootDir)) // serve from project root
			mux := http.NewServeMux()

			// Custom handler to support clean URLs and force index.html
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" {
					http.Redirect(w, r, dashboardPath, http.StatusFound)
					return
				}

				fs.ServeHTTP(w, r)

			})

			server := &http.Server{
				Addr:    "0.0.0.0" + port,
				Handler: mux,
			}
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("Server failed: %v\n", err)
			}
		}()

		// Give the server a moment to start
		time.Sleep(800 * time.Millisecond)

		//Cross-platform open browser
		var openCmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin": // macOS
			openCmd = exec.Command("open", fullURL)
		case "windows":
			openCmd = exec.Command("powershell", "-NoProfile", "-Command", "Start-Process", fullURL)
		default: // Linux, WSL, and others
			openCmd = exec.Command("xdg-open", fullURL)
		}

		if err := openCmd.Start(); err != nil {
			fmt.Printf("Auto-open failed: %v\n", err)
			fmt.Println("Dashboard available at:", fullURL)
			return
		}

		fmt.Println("Opening dashboard in your browser...")
		fmt.Println("If it doesn't open automatically, visit:", fullURL)
		fmt.Println("Press Ctrl+C to stop the server when done.")

		// Keep the server running
		select {}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(dashboardCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
