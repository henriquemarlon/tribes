<p align="center">
    <img src="https://github.com/user-attachments/assets/275b9ce4-3a4b-4965-82a3-5b6160ea76a5" align="center" width="30%">
</p>
<br>
<div align="center">
    <i>An Typescript example using Cartesi Cli, Nonodo and Deroll as High-Level Framework</i>
</div>
<div align="center">
<b>This example aims to demonstrate the lifecycle of a DApp using Deroll. In addition it serves.</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/tribes?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/henriquemarlon/tribes?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/henriquemarlon/tribes?style=default&color=0080ff" alt="repo-language-count">
</p>
<p align="center"><!-- default option, no dependency badges. -->
</p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>
<br>

##  Table of Contents

- [ Overview](#-overview)
- [ Getting Started](#-getting-started)
  - [ Prerequisites](#-prerequisites)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Testing](#-testing)
- [ Project Roadmap](#-project-roadmap)
- [ Contributors](#-contributors)
- [ Acknowledgments](#-acknowledgments)

##  Overview

Here's a compelling overview of the Tribes project:

Tribes is an open-source project that solves the problem of decentralized data management for large-scale applications. Its key features include automated workflow orchestration through a Makefile, dependency management using Go modules and cryptographic hashes in the `go.sum` file, and support for popular libraries like Gorm, Cobra, and Wire. This project benefits developers by streamlining their workflow, ensuring code integrity, and providing a reliable foundation for building scalable software systems. The target audience includes data engineers, software architects, and developers working on large-scale applications that require efficient data processing and management.


##  Project Structure

```sh
└── tribes/
    ├── .github
    │   ├── dependabot.yaml
    │   └── workflows
    ├── LICENSE
    ├── Makefile
    ├── README.md
    ├── build
    │   └── Dockerfile.machine
    ├── cmd
    │   ├── tribes-auth
    │   ├── tribes-ramp
    │   ├── tribes-rollup
    │   └── tribes-web3-provider
    ├── configs
    │   └── sqlite.go
    ├── coverage.md
    ├── go.mod
    ├── go.sum
    ├── internal
    │   ├── domain
    │   ├── infra
    │   └── usecase
    ├── pkg
    │   ├── rollups_contracts
    │   └── router
    ├── test
    │   └── integration
    └── website
```

##  Getting Started

###  Prerequisites

Before getting started with tribes, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go
- **Package Manager:** Go modules, Cargo

###  Installation

Install tribes using one of the following methods:

**Build from source:**

1. Clone the tribes repository:
```sh
❯ git clone https://github.com/henriquemarlon/tribes
```

2. Navigate to the project directory:
```sh
❯ cd tribes
```

3. Install the project dependencies:

**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go build
```

**Using `cargo`** &nbsp; [<img align="center" src="" />]()

```sh
❯ echo 'INSERT-INSTALL-COMMAND-HERE'
```

###  Usage
Run tribes using the following command:
**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go run {entrypoint}
```

**Using `cargo`** &nbsp; [<img align="center" src="" />]()

```sh
❯ echo 'INSERT-RUN-COMMAND-HERE'
```

###  Testing
Run the test suite using the following command:
**Using `go modules`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Go-00ADD8.svg?style={badge_style}&logo=go&logoColor=white" />](https://golang.org/)

```sh
❯ go test ./...
```

**Using `cargo`** &nbsp; [<img align="center" src="" />]()

```sh
❯ echo 'INSERT-TEST-COMMAND-HERE'
```

##  Project Roadmap

- [X] **`Task 1`**: <strike>Implement feature one.</strike>
- [ ] **`Task 2`**: Implement feature two.
- [ ] **`Task 3`**: Implement feature three.

##  Contributors

<p align="left">
   <a href="https://github.com{/henriquemarlon/tribes/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=henriquemarlon/tribes">
   </a>
</p>

##  Acknowledgments

- List any resources, contributors, inspiration, etc. here.
