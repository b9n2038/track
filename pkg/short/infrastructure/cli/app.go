// pkg/infrastructure/cli/app.go
package cli

import (
	"act/pkg/short/domain/model"
	"act/pkg/short/domain/ports/primary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct {
	service primary.ListService
	rootCmd *cobra.Command
}

type Profile struct {
	Name    string
	BaseDir string
	// other profile-specific fields
}

type Config struct {
	CurrentProfile string
	Profiles       map[string]Profile
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".short")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.short")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Create default config
			return &Config{
				Profiles: make(map[string]Profile),
			}, nil
		}
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func NewApp(service primary.ListService) *App {
	var pathFlag string

	rootCmd := &cobra.Command{
		Use:   "short",
		Short: "Short List - Helping you separate the critictal few from the trivial many. ",
		Long:  "Life is short, lists are long... Maybe this will a finite list tool can help? ",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			path, _ := cmd.Flags().GetString("path")
			fmt.Println("PersistentPreRun", path)
			pathFlag = path
		},
	}
	rootCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", "", "Override default storage location")

	app := &App{
		service: service,
		rootCmd: rootCmd,
	}

	// config, err := LoadConfig()
	// if config == nil || err != nil {
	// 	fmt.Println("no config")
	// } else {
	// 	fmt.Println("found config", config.CurrentProfile)
	// }
	//
	fmt.Println("pathFlag", pathFlag)
	app.setupCommands()
	return app
}

func (a *App) Run() error {
	return a.rootCmd.Execute()
}

func (a *App) setupCommands() {
	a.rootCmd.AddCommand(
		a.newAddListCommand(),
		a.newAddCommand(),
		a.newListCommand(),
		a.newCloseCommand(),
		a.newOpenCommand(),
		a.newConfigCommand(),
	)
}

func (a *App) GetStoragePath() string {

	path1, _ := a.rootCmd.PersistentFlags().GetString("path")
	path2, _ := a.rootCmd.Flags().GetString("path")
	path3, _ := a.rootCmd.InheritedFlags().GetString("path")

	fmt.Printf("Debug - PersistentFlags path: %q\n", path1)
	fmt.Printf("Debug - Flags path: %q\n", path2)
	fmt.Printf("Debug - InheritedFlags path: %q\n", path3)

	path, _ := a.rootCmd.Flags().GetString("path")
	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}
	// Convert to absolute path
	if !filepath.IsAbs(path) {
		currentDir, err := os.Getwd()
		if err == nil {
			path = filepath.Join(currentDir, path)
		}
	}

	// Clean the path to remove any .. or . elements
	return filepath.Clean(path)
}

func (a *App) newAddListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-list <name>",
		Short: "Create a new list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			maxCount, _ := cmd.Flags().GetInt("max-count")
			limitHandling, _ := cmd.Flags().GetString("limit-handling")

			config := model.DefaultConfig()
			if maxCount > 0 {
				config.MaxCount = maxCount
			}
			if limitHandling != "" {
				if limitHandling != string(model.MoveLastToClosed) && limitHandling != string(model.PushFront) {
					return fmt.Errorf("invalid limitHandling: must be moveLastToClosed or pushFront")
				}
				config.LimitHandling = model.LimitHandling(limitHandling)
			}

			if err := a.service.CreateList(args[0], config); err != nil {
				return fmt.Errorf("failed to create list: %w", err)
			}

			fmt.Printf("Created new list '%s' with maxCount=%d and limitHandling=%s\n",
				args[0], config.MaxCount, config.LimitHandling)
			return nil
		},
	}

	cmd.Flags().Int("max-count", 3, "Maximum number of items in open list")
	cmd.Flags().String("limit-handling", string(model.MoveLastToClosed),
		"How to handle list limits (moveLastToClosed or pushFront)")

	return cmd
}

func (a *App) newAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add <list> <item>",
		Short: "Add an item to the open list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.service.AddItem(args[0], args[1]); err != nil {
				return fmt.Errorf("failed to add item: %w", err)
			}

			fmt.Printf("Added item to list '%s'\n", args[0])
			return nil
		},
	}
}

func (a *App) newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list <name>",
		Short: "Show all items in a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := a.service.GetList(args[0])
			if err != nil {
				return fmt.Errorf("failed to get list: %w", err)
			}

			fmt.Printf("List: %s (max: %d, handling: %s)\n\n",
				list.Name, list.Config.MaxCount, list.Config.LimitHandling)

			fmt.Println("Open items:")
			if len(list.Open) == 0 {
				fmt.Println("  (empty)")
			}
			for i, item := range list.Open {
				fmt.Printf("  %d: %s\n", i, item)
			}

			fmt.Println("\nClosed items:")
			if len(list.Closed) == 0 {
				fmt.Println("  (empty)")
			}
			for i, item := range list.Closed {
				fmt.Printf("  %d: %s\n", i, item)
			}

			return nil
		},
	}
}

func (a *App) newCloseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "close <list> <index>",
		Short: "Move an item from open to closed list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			index, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid index: %w", err)
			}

			if err := a.service.MoveToClosed(args[0], index); err != nil {
				return fmt.Errorf("failed to close item: %w", err)
			}

			fmt.Printf("Moved item at index %d to closed list\n", index)
			return nil
		},
	}
}

func (a *App) newOpenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "open <list> <index>",
		Short: "Move an item from closed to open list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			index, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid index: %w", err)
			}

			if err := a.service.MoveToOpen(args[0], index); err != nil {
				return fmt.Errorf("failed to open item: %w", err)
			}

			fmt.Printf("Moved item at index %d to open list\n", index)
			return nil
		},
	}
}

func (a *App) newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config <list> <setting> <value>",
		Short: "Configure list settings",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := a.service.GetList(args[0])
			if err != nil {
				return fmt.Errorf("failed to get list: %w", err)
			}

			config := list.Config
			setting := args[1]
			value := args[2]

			switch setting {
			case "max-count":
				count, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("invalid maxCount: %w", err)
				}
				config.MaxCount = count
			case "limit-handling":
				if value != string(model.MoveLastToClosed) && value != string(model.PushFront) {
					return fmt.Errorf("invalid limitHandling: must be moveLastToClosed or pushFront")
				}
				config.LimitHandling = model.LimitHandling(value)
			default:
				return fmt.Errorf("unknown setting: %s", setting)
			}

			if err := a.service.UpdateConfig(args[0], config); err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}

			fmt.Printf("Updated %s setting for list '%s'\n", setting, args[0])
			return nil
		},
	}

	return cmd
}
