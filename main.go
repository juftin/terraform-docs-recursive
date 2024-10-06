package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/terraform-docs/terraform-docs/cmd"
)

const (
	exitCodeError = 1
)

var (
	configFileNames = []string{
		".terraform-docs.yml",
		".terraform-docs.yaml",
	}
	skipDirs = map[string]bool{
		".terragrunt-cache": true,
		".terraform":        true,
		".git":              true,
		"node_modules":      true,
		".venv":             true,
		"vendor":            true,
		"dist":              true,
		"build":             true,
		"target":            true,
		".idea":             true,
		".vscode":           true,
	}
)

// TFDocsCLI executes the terraform-docs command with the provided arguments.
func TFDocsCLI(args []string) error {
	tfDocs := cmd.NewCommand()
	tfDocs.SetArgs(args)
	return tfDocs.Execute()
}

// RelativePathsToAbsolute converts a list of relative file paths to absolute paths.
func RelativePathsToAbsolute(args []string) ([]string, error) {
	var paths []string
	for _, arg := range args {
		absPath, err := filepath.Abs(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for %s: %w", arg, err)
		}
		paths = append(paths, absPath)
	}
	return paths, nil
}

// PathsToUniqueDirs returns a list of unique directories from a list of file paths.
func PathsToUniqueDirs(paths []string) []string {
	uniqueDirs := make(map[string]bool)
	var dirs []string
	for _, path := range paths {
		dir := filepath.Dir(path)
		if !uniqueDirs[dir] {
			uniqueDirs[dir] = true
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

// ValidateConfigFileExists checks for a terraform-docs configuration file in the provided directory.
func ValidateConfigFileExists(dir string) string {
	for dir != filepath.Dir(dir) {
		for _, configFile := range configFileNames {
			configFilePath := filepath.Join(dir, configFile)
			if _, err := os.Stat(configFilePath); err == nil {
				return configFilePath
			}
		}
		dir = filepath.Dir(dir)
	}
	return ""
}

// ValidateDirIsTerraformModule checks if the provided directory is a Terraform module.
func ValidateDirIsTerraformModule(dir string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".tf" || ext == ".tf.json" {
			return true, nil
		}
	}
	return false, nil
}

// TFDocsPreCommitCommand processes the provided arguments and runs terraform-docs on valid modules.
func TFDocsPreCommitCommand(args []string) error {
	paths, err := RelativePathsToAbsolute(args)
	if err != nil {
		return err
	}
	dirs := PathsToUniqueDirs(paths)
	for _, path := range dirs {
		isTerraformModule, err := ValidateDirIsTerraformModule(path)
		if err != nil {
			return err
		}
		configFilePath := ValidateConfigFileExists(path)
		if isTerraformModule && configFilePath != "" {
			if err := TFDocsCLI([]string{"--config", configFilePath, path}); err != nil {
				return fmt.Errorf("failed to run terraform-docs on %s: %w", path, err)
			}
		}
	}
	return nil
}

// Module represents a Terraform module with its path and config file.
type Module struct {
	Path       string
	ConfigFile string
}

// FindAllModules traverses the provided directory and returns a list of all Terraform modules.
func FindAllModules(dir string) ([]Module, error) {
	var modules []Module
	var mu sync.Mutex
	var wg sync.WaitGroup

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if skipDirs[filepath.Base(path)] {
			return filepath.SkipDir
		}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			validTerraformModule, err := ValidateDirIsTerraformModule(path)
			if err != nil {
				log.Printf("Warning: %v", err)
			}
			if validTerraformModule {
				configFilePath := ValidateConfigFileExists(path)
				if configFilePath != "" {
					mu.Lock()
					modules = append(modules, Module{Path: path, ConfigFile: configFilePath})
					mu.Unlock()
				}
			}
		}(path)

		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory %s: %w", dir, err)
	}
	return modules, nil
}

// TFDocsFindAllCommand searches for all Terraform modules in the provided directory and runs terraform-docs on each module.
func TFDocsFindAllCommand(dir string) error {
	resolvedDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %s: %w", dir, err)
	}
	allModules, err := FindAllModules(resolvedDir)
	if err != nil {
		return err
	}
	for _, module := range allModules {
		if err := TFDocsCLI([]string{"--config", module.ConfigFile, module.Path}); err != nil {
			return fmt.Errorf("failed to run terraform-docs on %s: %w", module.Path, err)
		}
	}
	return nil
}

// TFDocsPreCommitCLI creates a new cobra.Command for the terraform-docs pre-commit CLI.
func TFDocsPreCommitCLI() *cobra.Command {
	cli := &cobra.Command{
		Use:           "terraform-docs-pre-commit",
		Short:         "A CLI for terraform-docs pre-commit hooks",
		Long:          "A CLI for terraform-docs pre-commit hooks",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	preCommitCmd := &cobra.Command{
		Use:   "pre-commit [files...]",
		Short: "Run pre-commit hook",
		Long:  "Run pre-commit hook for terraform-docs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return TFDocsPreCommitCommand(args)
		},
	}
	cli.AddCommand(preCommitCmd)

	findAllCmd := &cobra.Command{
		Use:   "docs [directory]",
		Short: "Generate terraform-docs for all modules",
		Long:  "Generate terraform-docs for all Terraform modules in the provided directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			resolvedDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to resolve absolute path for %s: %w", dir, err)
			}
			info, err := os.Stat(resolvedDir)
			if os.IsNotExist(err) || !info.IsDir() {
				return fmt.Errorf("provided path is not a valid directory: %s", resolvedDir)
			}
			return TFDocsFindAllCommand(resolvedDir)
		},
	}
	cli.AddCommand(findAllCmd)

	return cli
}

func main() {
	if err := TFDocsPreCommitCLI().Execute(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(exitCodeError)
	}
}
