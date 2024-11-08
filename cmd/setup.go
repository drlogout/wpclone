package cmd

import (
	"bufio"
	"croox/wpclone/config"
	"croox/wpclone/pkg/message"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var Setup = &cli.Command{
	Name:  "setup",
	Usage: "Setup wpclone (add binary to PATH)",
	Action: func(ctx *cli.Context) error {
		binPath := filepath.Join(config.ConfigDir(), "bin")

		if err := addPathToShells(binPath); err != nil {
			return fmt.Errorf("failed to add path to shells: %v", err)
		}

		message.Info("For bash, you may need to restart your shell or run `source ~/.bashrc`")
		message.Info("For zsh, you may need to restart your shell or run `source ~/.zshrc`")
		message.Info("For fish, you may need to restart your shell or run `source ~/.config/fish/config.fish`")

		message.Successf("Successfully set up wpclone")
		return nil
	},
}

func addPathToShells(binPath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Define the paths to the shell configuration files
	bashConfig := filepath.Join(homeDir, ".bashrc")
	zshConfig := filepath.Join(homeDir, ".zshrc")
	fishConfig := filepath.Join(homeDir, ".config", "fish", "config.fish")

	// Define the export command for each shell
	bashCommand := fmt.Sprintf("export PATH=\"$PATH:%s\"\n", binPath)
	zshCommand := fmt.Sprintf("export PATH=\"$PATH:%s\"\n", binPath)
	fishCommand := fmt.Sprintf("set -x PATH $PATH %s\n", binPath)

	// Append the command to each shell configuration file
	if err := appendToFile(bashConfig, bashCommand); err != nil {
		return fmt.Errorf("failed to update bash config: %v", err)
	}

	if err := appendToFile(zshConfig, zshCommand); err != nil {
		return fmt.Errorf("failed to update zsh config: %v", err)
	}

	if err := appendToFile(fishConfig, fishCommand); err != nil {
		return fmt.Errorf("failed to update fish config: %v", err)
	}

	return nil
}

// appendToFile appends text to a file, creating the file if it doesn't exist
func appendToFile(filePath, text string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		return nil
	}

	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(text) {
			// The line already exists, no need to append
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error %v", err)
	}

	// Append the line since it doesn't exist
	if _, err := file.WriteString("\n# Add wpclone to PATH\n" + text + "\n"); err != nil {
		return err
	}

	return nil
}
