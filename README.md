# terraform-docs-recursive

Recursively generate terraform documentation for all modules in a directory.

[![Homebrew](https://img.shields.io/github/v/release/juftin/terraform-docs-recursive?label=brew&color=blue&logo=homebrew)](https://brew.sh/)
[![Go Reference](https://pkg.go.dev/badge/github.com/juftin/terraform-docs-recursive.svg)](https://pkg.go.dev/github.com/juftin/terraform-docs-recursive)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-grey?logo=github)](https://github.com/features/actions)

This project is a wrapper around [terraform-docs](https://github.com/terraform-docs/terraform-docs/)
with logic to recursively find and generate documentation for all modules in a directory. In order for
a module to be considered for documentation, its directory (or any of its parent directories)
must contain a [terraform-docs configuration file](https://terraform-docs.io/user-guide/configuration/).
`terraform-docs-recursive` will search use whichever configuration file is found first in the directory tree.

## Usage

### Run as a CLI

```bash
terraform-docs-recursive docs .
```

#### Install with homebrew

```bash
brew tap juftin/terraform-docs-recursive https://github.com/juftin/terraform-docs-recursive
brew install terraform-docs-recursive
```

#### Install with Go

```bash
go install github.com/juftin/terraform-docs-recursive@v1.6.0
```

### Pre-Commit Hook

To use this project as a [pre-commit hook](https://pre-commit.com/)
you can add the following to your `.pre-commit-config.yaml`:

```yaml
repos:
    - repo: https://github.com/juftin/terraform-docs-recursive
      rev: v1.6.0
      hooks:
          - id: terraform-docs-recursive
```

### GitHub Actions

```yaml
name: Generate Terraform Docs
on:
    push:
runs-on: ubuntu-latest
steps:
    - name: Checkout Repository
      uses: actions/checkout@v4
    - name: Generate Terraform Docs
      uses: juftin/terraform-docs-recursive@v1
    - name: Commit Changes
      uses: EndBug/add-and-commit@v9
      with:
          message: 📝 update terraform documentation
          default_author: github_actions
```
