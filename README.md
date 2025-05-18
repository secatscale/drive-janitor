# 🚗 drive-janitor

A simple tool to help you detect and delete files on your disk based on their extension, age, or name — even those in the trash.

## Features

- 🧹 Detect and delete unwanted files and folders by extension, age, or name with regex
- 🗑️ Clean files even from the trash
- 🔒 Safe: dry-run mode to preview changes
- ⚡ Fast and easy to use

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/drive-janitor.git
cd drive-janitor
```

## Usage

1. Configure your settings in `config.yaml`.
2. Run the janitor:

```bash
go run drive-janitor -config config.yaml
```

## Configuration

Edit `config.yaml` to specify:

- **Detections**:  
    Define which files to match in the `detections` section. Each detection uses a unique `name`, a `mimetype` (e.g., `"image/png"`), and optionally a `max_age` (in days) to filter files by age. Setting `max_age` to `-1` disables age filtering.

- **Recursions**:  
    Specify which directories to scan in the `recursions` section. Each recursion includes a `name`, a `path` (e.g., `"./samples"`), a `max_depth` for how deep to scan, and an optional `path_to_ignore` list to exclude subdirectories.

- **Logs**:  
    Configure where logs are stored in the `logs` section by setting a `name` and a `log_repository` path (e.g., `"./var/.log"`).

- **Actions**:  
    Define what happens to matched files in the `actions` section. Each action has a `name`, a `delete` flag (e.g., `false` to only log), and a `log` reference.

- **Rules**:  
    Connect detections, recursions, and actions in the `rules` section. Each rule specifies which detections, recursion, and action to use for processing files.

Example:
```yaml
name: TestMain
version: 1.0.0

detections:
    - name: "detection1"
        mimetype: "image/png"
        max_age: 5

    - name: "detection2"
        mimetype: "application/pdf"

    - name: "detection3"
        mimetype: "video/mp4"

recursions:
    - name: "recursion1"
        path: "./samples"
        max_depth: 4
        path_to_ignore: ["./samples/ignore"]

logs:
    - name: "logging1"
        log_repository: "./var/.log"

actions:
    - name: "action1"
        delete: false 
        log: "logging1"

rules:
    - name: "FINALRULE1"
        action: "action1"
        detection: ["detection1", "detection2", "detection3"]
        recursion: "recursion1"
```

## Contributing

Pull requests are welcome! 

## License

[GPLv3](LICENSE)

---

*Happy cleaning! 🧹*
