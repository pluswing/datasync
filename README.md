# DataSync

Streamline Your Development, Synchronize Your Success

## Overview
DataSync is a powerful CLI tool designed for software developers. It simplifies the process of dumping, applying, and managing the history of various data types, such as model files and databases. This tool is particularly beneficial for projects involving MySQL databases and specific file directories.

## Features
DataSync includes several key features:

- Data Dump (datasync dump): Dumps data into a zip file stored in the .datasync folder. The dump timestamp and data ID are saved in the .datasync_version file.
- Data Apply (datasync apply): Restores data from a zip file specified by the .datasync_version ID or an argument-provided ID.
- List Dumps (datasync ls): Displays a list of all dumped data.
- Push to Storage (datasync push): Uploads dumped files from the .datasync folder to the specified Google Cloud Storage.
- Pull from Storage (datasync pull): Downloads dump data from the remote storage and places it in the .datasync folder.

## Configuration
DataSync utilizes a YAML configuration file for setting up target data sources and storage options. Below is an example configuration:

```yaml
Copy code
targets:
  - kind: mysql
    config:
      host: localhost
      port: 3306
      user: root
      password: root
      database: sample
  - kind: file
    config:
      path: "./cmd"
  - kind: file
    config:
      path: "./dump/dump_mysql"
  - kind: file
    config:
      path: "./.vscode"
  - kind: file
    config:
      path: "./LICENSE"

storage:
  kind: gcs
  config:
    bucket: datasync000001
    dir: "test01/test02"
```

## Usage Scenarios

For New Developers
New team members can quickly set up their development environment with just a few commands:

```bash
git clone ... <project dir>
cd <project dir>
datasync pull  # Download data
datasync apply  # Apply data
For Solo Development
Solo developers can use DataSync locally without remote storage:
```

```bash
datasync dump -m "test"
# Perform model modifications and performance evaluation
datasync dump -m "test02"  # Dump data after successful changes
# datasync apply  # Revert to previous state if needed
```

## Installation and Setup

Download binary from Github Release page.


## Contributing
[Guidelines for contributing to the development of DataSync, including how to submit issues and pull requests.]

## License

MIT License

## Contact
[Contact information or links for support or further inquiries about DataSync.]


DataSync is developed to streamline the workflow of developers, ensuring efficient and error-free data management in software development projects.
