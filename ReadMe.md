# Maya CLI

**Maya** is a modular command-line tool designed to simplify media management and downloading workflows. It provides intelligent download handling for movies, TV series, and other multimedia content with metadata extraction and organization.

---

## Features

- Download movies and series from supported sources
- Extract metadata in Jellyfin-friendly format
- Organize episodes for series with proper folder structure
- Resume interrupted downloads using temporary files
- Handle M3U8 playlists and direct links
- Parallel downloads for faster series downloading
- Extensible architecture for future tools and modules

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/ajaysinghnp/maya-cli.git
cd maya-cli
````

1. Build the CLI:

```bash
go build -o maya main.go
```

1. Optionally, move `maya` to a directory in your `$PATH`:

```bash
sudo mv maya /usr/local/bin/
```

---

## Usage

```bash
maya [command] [flags]
```

### Download a movie

```bash
maya download <url>
```

### Download a series with concurrency

```bash
maya download <url> --concurrency 5
```

### Help

```bash
maya --help
maya download --help
```

---

## Flags for `download`

- `-o, --output` : Specify output directory or filename (default: auto-generated)
- `-r, --resume` : Resume interrupted download if cached files exist
- `-c, --concurrency` : Number of simultaneous downloads for series episodes
- `-v, --verbose` : Enable verbose logging to terminal

---

## Extensibility

Maya is modular, so you can add new commands or tools easily. All subcommands follow the same pattern:

```bash
maya <tool> [options]
```

---

## Logging

Maya supports:

- Colorful terminal logging (info, success, debug, warn, error)
- Optional file logging (can be added via logger configuration)

---

## License

MIT License. See [LICENSE](LICENSE) for details.
