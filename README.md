<p align="center">
    <img src="https://github.com/user-attachments/assets/275b9ce4-3a4b-4965-82a3-5b6160ea76a5" align="center" width="30%">
</p>
<div align="center">
    <i>The new frontier of financial services for the creators' economy</i>
</div>
<div align="center">
<b>Collateralized tokenization of receivables for debt issuance through crowdfunding</b>
</div>
<br>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/tribes?style=default&logo=opensourceinitiative&logoColor=white&color=959CD0" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/tribes?style=default&logo=git&logoColor=white&color=D1DCCB" alt="last-commit">
</p>

##  Table of Contents

- [Overview](#-overview)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Testing](#testing)
- [Contributors](#contributors)

##  Overview

<div align="justify">
A crowdfunding platform designed for prominent content creators, enabling them to monetize their influence by issuing tokenized debt instruments collateralized exclusively by their tokenized future receivables. Based on Resolution No. 88 of the Brazilian Securities and Exchange Commission (CVM), the Brazilian SEC, the platform connects creators with a network of investors, offering a structured and transparent alternative to finance scalable businesses while leveraging the economic potential of their audiences, ensuring legal compliance and attractive returns for investors.
</div>

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
    └── website
```

##  Getting Started

###  Prerequisites
1. [Install Docker Desktop for your operating system](https://www.docker.com/products/docker-desktop/).

    To install Docker RISC-V support without using Docker Desktop, run the following command:
    
   ```shell
   ❯ docker run --privileged --rm tonistiigi/binfmt --install all
   ```

2. [Download and install the latest version of Node.js](https://nodejs.org/en/download).

3. Cartesi CLI is an easy-to-use tool to build and deploy your dApps. To install it, run:

   ```shell
   ❯ npm i -g @cartesi/cli
   ```

4. [Download and Install the latest version of Golang.](https://go.dev/doc/install)

5. Install development node:

   ```shell
   ❯ npm i -g nonodo
   ```

6. Install air ( hot reload tool ):

   ```shell
   ❯ go install github.com/air-verse/air@latest
   ```

###  Running

**Build rollup from image**

```sh
❯ docker pull ghcr.io/tribeshq/tribes-machine:latest
```

**Generate rollup filesystem**

```sh
❯ cartesi build --from-image ghcr.io/henriquemarlon/tribes-machine
```

**Run validator node**

```sh
❯ cartesi run
```

###  Testing

```sh
❯ make test
```

##  Contributors

<p align="left">
   <a href="https://github.com{/henriquemarlon/tribes/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=henriquemarlon/tribes">
   </a>
</p>
