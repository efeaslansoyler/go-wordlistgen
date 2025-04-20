# Go-Wordlistgen 🔑

[![en](https://img.shields.io/badge/lang-en-red.svg)](README.md)
[![tr](https://img.shields.io/badge/lang-tr-green.svg)](README.tr.md)

A powerful wordlist generator for password cracking written in Go, that creates customized wordlists based on personal information.

## Features

- 🖥️ Dual Interface: Choose between an interactive Terminal User Interface (TUI) or Command Line Interface (CLI)
- 👤 Personal Info Based: Generate wordlists using:
  - First name and last name
  - Birthday
  - Related words
- 🔄 Advanced Variations:
  - Leet speak (1337) transformations
  - Capitalization variations
  - Length constraints
- 💾 Customizable output file location

## Installation

```bash
go install github.com/efeaslansoyler/go-wordlistgen@latest
```

## Usage

### TUI Mode (Default)

Simply run:
```bash
go-wordlistgen
```

### CLI Mode

```bash
go-wordlistgen --cli [options]

Options:
  -c, --cli                Run in CLI mode
  -f, --firstname string   First name (and middle name if needed)
  -l, --lastname string    Last name
  -b, --birthday string    Birthday in format DD/MM/YYYY
  -w, --words string       Related words separated by commas
      --min string        Minimum password length (default "6")
      --max string        Maximum password length (default "12")
  -o, --output string     Output file path (default "wordlist.txt")
      --leet             Enable leet speak variations
      --caps             Enable capitalization variations
```

Example:
```bash
go-wordlistgen --cli -f "John" -l "Doe" -b "01/01/1990" -w "hobby,pet,city" --leet --caps
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Efe Aslan Söyler (efeaslan1703@gmail.com)