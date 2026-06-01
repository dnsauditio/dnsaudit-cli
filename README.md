# DNSAudit CLI

A fast, robust, and professional command-line interface for the [DNSAudit.io](https://dnsaudit.io) API. Built specifically for DevOps teams, system administrators, and penetration testers to seamlessly run DNS security scans from the terminal.

## Installation

As a Go binary, `dnsaudit` has zero dependencies. Simply compile the binary and move it to your path:

```bash
# Clone the repository
git clone https://github.com/dnsaudit/cli dnsaudit-cli
cd dnsaudit-cli

# Build for your current system
go build -o dnsaudit main.go

# --- Cross-Compilation ---
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o dnsaudit-windows.exe main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o dnsaudit-linux-amd64 main.go

# Build for Mac (Apple Silicon / ARM64)
GOOS=darwin GOARCH=arm64 go build -o dnsaudit-mac-arm64 main.go

# (Optional) Move to your path for global access
sudo mv dnsaudit /usr/local/bin/
```

## Authentication

Every API request requires a valid DNSAudit API key. You can authenticate in one of two ways:

1. **Environment Variable** (Recommended for CI/CD pipelines & Docker):
   ```bash
   export DNSAUDIT_API_KEY="your-api-key"
   ```
2. **Configuration File**:
   Run the setup command to save your key securely to `~/.config/dnsaudit/config.json`:
   ```bash
   dnsaudit configure
   ```

---

## Commands & Usage

### `scan`
Runs a comprehensive DNS security scan on the target domain. Automatically formats the output into a clean, readable summary containing security grades, top issues, and remediation recommendations.

```bash
dnsaudit scan -d example.com
```

### `history`
Fetches a list of your most recently scanned domains alongside their historical security scores.

```bash
dnsaudit history --limit 20
```

### `export`
Downloads a full security report for a specific domain to your local disk. 

```bash
# Download a detailed PDF report
dnsaudit export -d example.com -f pdf

# Download raw JSON results
dnsaudit export -d example.com -f json -o output_file.json
```

---

## Flags Reference

| Command | Flag | Short | Default | Description |
| :--- | :--- | :---: | :---: | :--- |
| **Global** | `--help` | `-h` | | Show help and usage for any command |
| **`scan`** | `--domain` | `-d` | | **(Required)** The target domain to scan |
| | `--json` | `-j` | `false` | Output the raw JSON response (ideal for piping into `jq`) |
| **`history`**| `--limit` | `-l` | `10` | The maximum number of historical results to return (max `100`) |
| | `--json` | `-j` | `false` | Output the raw JSON response |
| **`export`** | `--domain` | `-d` | | **(Required)** The target domain to export data for |
| | `--format` | `-f` | `pdf` | The export format (either `pdf` or `json`) |
| | `--output` | `-o` | `<domain>-report.<fmt>`| Custom output file path |

---

## Rate Limiting

The CLI is designed to be respectful of the DNSAudit API limits.
- **Burst Limit Handling**: If the API returns a `429 Too Many Requests` due to burst limits, the CLI will automatically parse the required wait time, pause execution silently in the background, and retry your request seamlessly without crashing.
- **Daily Limits**: If you exceed your account's daily scan limits, the CLI will exit gracefully with a clear error message.

## Output Examples

### Standard Terminal Output
```text
$ dnsaudit scan -d example.com
[*] Starting scan for example.com...
--------------------------------------------------------------------------------
Target: example.com
Status: success

[+] Security Grade: B+ (Score: 87)
[*] Good DNS security with minor improvements recommended

Issues Breakdown:
  Critical: 0
  Warning:  3
  Info:     8

Top Recommendations:
  - Address warning-level issues to improve security grade
  - Review and strengthen DNS security configuration

Key Findings:
  [warning] DMARC: Missing DMARC Reporting Addresses
  [warning] CAA: Missing CAA records
  [warning] SPF: Missing Wildcard SPF for Subdomains
--------------------------------------------------------------------------------
```

### Automation & Scripting (JSON)
The CLI plays perfectly with bash pipelines:
```bash
$ dnsaudit scan -d example.com -j | jq '.grade.score'
87
```
