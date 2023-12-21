# Prompt Generator

## Description
This go program generates a prompt from provided prompt templates and specified files in the working directory.


## Configuration
### Prompt Templates
Inside the Main.prompt file, you can specify the prompt templates to use. The prompt templates are defined in the `example/prompts/templates` directory.
The Files.prompt and the Anonymize.prompt files are special prompt templates that are used to specify the files to use and the anonymization tokens to apply.


## Usage (Linux amd64)
### Install
```bash
./scripts/amd64/linux/install.sh
```

### Run
```bash
prompt-generator
```


## Development
### Build from source
```bash
./scripts/amd64/linux/build.sh
```

### Run from binary (after build)
```bash
./scripts/amd64/linux/run.sh
```

### Build source and run from binary (in one step)
```bash
./scripts/amd64/linux/build_and_run.sh
```
