# 🚗 drive-janitor

Drive Janitor is a file detection and deletion tool, highly configurable via a `config.yaml` file.
It can detect files based on MIME type, age, filename (regex), or YARA rules, recursively starting from one or more source directories.

*Suggested use case:*
-> Detect PDF files older than 30 days in the Trash.

Runs on Windows, Linux, and macOS.

### Important Notes

On macOS, in order to scan the Trash directory, you must grant Full Disk Access to drive-janitor.

## Installation

You can download precompiled binairies in the release tab.

Alternativily, clone the repository:

```bash
	git clone https://github.com/secatscale/drive-janitor.git
	cd drive-janitor
```

Drive janitor requires pkg-config and yara installation : 


Linux installation 
```
	apt-get install libyara-dev yara pkg-config
```

MacOS installation
```
	brew install pkg-config yara
```

Windows installation
```
	vcpkg install yara
```

## Usage

1. Configure your settings in `config.yaml`.
2. Run the janitor:

```bash
go run drive-janitor -config config.yaml
```

or

```
./drive-janitor -config config.yaml
```



## Configuration

Edit `config.yaml` to specify:

- **Detections**:  
    Define which files to match in the `detections` section. Each detection rule uses a unique rule `name`, and a detection parameters :
    - `mimetype` (e.g., `"image/png"`) to filter files by MIME type. 
    - `max_age` (in days) to filter files by age. Setting `max_age` to `-1`, or leaving it unspecified disable age filtering.
    - `yara_rules_dir` to specify a folder where yara rules files are located. Files matching those rules are detected. Only one directory can be specified. 

- **Recursions**:  
    Specify which directories to scan in the `recursions` section. Each recursion includes a rule `name`, a start `path` (e.g., `"./samples"`), a `max_depth` for how deep to scan, and an optional `path_to_ignore` list to exclude subdirectories.

- **Logs**:  
    Configure where logs are stored in the `logs` section by setting a rule `name`, a `log_format` (text, json or csv) and a `log_repository` path (e.g., `"./var/.log"`).

- **Actions**:  
    Define what happens to matched files in the `actions` section. Each action has a rule `name`, a `delete` flag (e.g., `false` to only log), and a `log` reference.

- **Rules**:  
    Connect detections, recursions, and actions in the `rules` section. Each rule specifies which detections, recursion, and action to use for processing files.

⚡ All paths in the configuration supports `$trash`, `$download` ou `$home` 

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
