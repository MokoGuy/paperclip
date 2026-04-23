package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Set up paperclip configuration",
		Long:  "Create the config file with your Paperless-NGX URL and API token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := configFilePath()
			if err != nil {
				return err
			}

			if _, err := os.Stat(configPath); err == nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Config already exists at %s\n", configPath)
				fmt.Fprint(cmd.ErrOrStderr(), "Overwrite? [y/N] ")
				reader := bufio.NewReader(os.Stdin)
				answer, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(answer)) != "y" {
					fmt.Fprintln(cmd.ErrOrStderr(), "Aborted.")
					return nil
				}
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Fprint(cmd.ErrOrStderr(), "Paperless-NGX URL (e.g. https://paperless.example.com): ")
			url, _ := reader.ReadString('\n')
			url = strings.TrimSpace(url)
			if url == "" {
				return fmt.Errorf("URL is required")
			}
			url = strings.TrimRight(url, "/")

			fmt.Fprint(cmd.ErrOrStderr(), "API token: ")
			token, _ := reader.ReadString('\n')
			token = strings.TrimSpace(token)
			if token == "" {
				return fmt.Errorf("token is required")
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "Testing connection to %s... ", url)
			if err := testConnection(url, token); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "FAILED")
				return fmt.Errorf("connection test failed: %w", err)
			}
			fmt.Fprintln(cmd.ErrOrStderr(), "OK")

			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			content := fmt.Sprintf("url = %q\ntoken = %q\n", url, token)
			if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "Config saved to %s\n", configPath)
			fmt.Fprintln(cmd.ErrOrStderr(), "Run 'paperclip sync' to populate the local cache.")
			return nil
		},
	}
}

func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "paperclip", "config.toml"), nil
}

func testConnection(url, token string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url+"/api/tags/?page_size=1", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned %d (check URL and token)", resp.StatusCode)
	}
	return nil
}
