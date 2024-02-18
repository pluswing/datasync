# DataSync
[English](https://github.com/pluswing/datasync/blob/develop/README.md) [日本語](https://github.com/pluswing/datasync/blob/develop/README_ja.md)

A database sharing tool for developers

## What is this?
Have you ever needed to share the contents of a database with your team members during development work?
DataSync solves a common problem faced by developers.
Creating backups of databases, sharing them, and applying them can be cumbersome, time-consuming, and often complex.
Using DataSync, these processes become simple and efficient.
For example, you can easily pass on data needed for testing new features or reproducing bugs to your team members.
This speeds up the development cycle and allows you to focus on more productive activities.

## Overview
DataSync is a tool that enables backup, history management, and quick application of any backup for MySQL and files.

Key features include:

- Backup of databases and files
- History management of backups
- Application of any backup
- Easy sharing and retrieval of backups through integration with cloud storage

## Examples of Use

### Executing a Backup

```
$ datasync dump -m "feature_test"
✔️ mysql dump completed (database: sample)
✔︎ compress data completed.
Dump succeeded. Version ID = 35ca8d497d334891b2ff627174a2b88a
```

### Listing Backups
```
$ datasync ls -a
-- Remote versions --
224120fe68d14f6eaf2b4ea0533c497f 2024-01-30 13:53:21 test001
-- local versions --
35ca8d497d334891b2ff627174a2b88a 2024-02-10 09:33:27 test002
```

### Applying a Backup
```
$ datasync apply 35ca8d4
✔︎ decompress data completed.
✔︎ mysql import completed (database: sample)
Apply succeeded. Version ID = 35ca8d497d334891b2ff627174a2b88a
```

### Sending a Backup to the Cloud
```
$ datasync push
```

### Retrieving a Backup from the Cloud
```
$ datasync pull
```

## Installation
DataSync is a standalone single-file binary. Please follow the steps below:

1. Download the latest DataSync from the [Releases page](https://github.com/pluswing/datasync/releases).
2. Save the downloaded file to any location.
3. Move to the location where you saved the file from the command line, and run the following command to initialize DataSync:
```
datasync init
```
Now, you are ready to go!

## License
MIT License
