<p align="center">
    <img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" align="center" width="30%">
</p>
<p align="center"><h1 align="center">TRIBES.GIT</h1></p>
<p align="center">
	<em><code>❯ REPLACE-ME</code></em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/henriquemarlon/tribes.git?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/henriquemarlon/tribes.git?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/henriquemarlon/tribes.git?style=default&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/henriquemarlon/tribes.git?style=default&color=0080ff" alt="repo-language-count">
</p>
<p align="center"><!-- default option, no dependency badges. -->
</p>
<p align="center">
	<!-- default option, no dependency badges. -->
</p>
<br>

##  Table of Contents

- [ Overview](#-overview)
- [ Features](#-features)
- [ Project Structure](#-project-structure)
  - [ Project Index](#-project-index)
- [ Getting Started](#-getting-started)
  - [ Prerequisites](#-prerequisites)
  - [ Installation](#-installation)
  - [ Usage](#-usage)
  - [ Testing](#-testing)
- [ Project Roadmap](#-project-roadmap)
- [ Contributing](#-contributing)
- [ License](#-license)
- [ Acknowledgments](#-acknowledgments)

---

##  Overview

<code>❯ REPLACE-ME</code>

---

##  Features

<code>❯ REPLACE-ME</code>

---

##  Project Structure

```sh
└── tribes.git/
    ├── .github
    │   └── workflows
    ├── LICENSE
    ├── Makefile
    ├── README.md
    ├── build
    │   ├── Dockerfile.machine
    │   └── Dockerfile.validator
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
        └── .gitkeep
```


###  Project Index
<details open>
	<summary><b><code>TRIBES.GIT/</code></b></summary>
	<details> <!-- __root__ Submodule -->
		<summary><b>__root__</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/go.mod'>go.mod</a></b></td>
				<td>- The `go.mod` file defines the project's module, `github.com/tribeshq/tribes`, and specifies its Go version and dependencies<br>- It lists required packages, including Ethereum libraries, dependency injection tools, cryptographic libraries, and database drivers, indicating a project likely involving blockchain technology and potentially a decentralized application<br>- Indirect dependencies highlight a complex project with numerous supporting libraries.</td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/Makefile'>Makefile</a></b></td>
				<td>- The Makefile orchestrates the project's build, testing, and documentation processes<br>- It defines targets for environment setup, building the Cartesi machine, running the application locally, generating data for rollups and contracts, executing tests with coverage reports, and launching documentation<br>- These targets streamline the development workflow.</td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/go.sum'>go.sum</a></b></td>
				<td>- The `go.sum` file serves as a checksum database for the project's dependencies<br>- It ensures the integrity and authenticity of the `github.com/DataDog/zstd` library (version 1.4.5), preventing accidental or malicious modification of downloaded packages<br>- This contributes to the overall security and reliability of the Go application.</td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- configs Submodule -->
		<summary><b>configs</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/configs/sqlite.go'>sqlite.go</a></b></td>
				<td>- The `sqlite.go` config file sets up a SQLite database connection within the Tribes application<br>- It initializes the database, runs migrations to create necessary tables for users, orders, contracts, and crowdfunding entities, and populates it with initial user data<br>- This ensures the application has a persistent data store for its core functionalities<br>- The configuration allows for both in-memory and persistent database usage.</td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- build Submodule -->
		<summary><b>build</b></summary>
		<blockquote>
			<table>
			<tr>
				<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/build/Dockerfile.machine'>Dockerfile.machine</a></b></td>
				<td>- The Dockerfile constructs a runtime environment for a Cartesi rollup application<br>- It cross-compiles a Rust library and a Go application for the RISC-V architecture, leveraging caching for efficiency<br>- The resulting image includes necessary tools and sets up a user for the application, ultimately preparing it for execution within the Cartesi machine emulator.</td>
			</tr>
			<tr>
				<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/build/Dockerfile.validator'>Dockerfile.validator</a></b></td>
				<td>- The Dockerfile configures a Cartesi rollups node image<br>- It leverages a base image and sets environment variables for snapshot directory and HTTP address<br>- Crucially, it copies a custom snapshot image into the node's designated snapshot directory, ensuring the node starts with a pre-configured state<br>- This facilitates consistent and reproducible deployments within the broader project's infrastructure.</td>
			</tr>
			</table>
		</blockquote>
	</details>
	<details> <!-- pkg Submodule -->
		<summary><b>pkg</b></summary>
		<blockquote>
			<details>
				<summary><b>rollups_contracts</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/rollups_contracts/input_box.go'>input_box.go</a></b></td>
						<td>- The `input_box.go` file within the `pkg/rollups_contracts` directory is a generated file containing the Go bindings for interacting with a smart contract named "InputBox" on the Ethereum blockchain<br>- It provides the necessary functions to interact with the contract's methods and events, enabling the broader project to seamlessly integrate with and utilize the InputBox contract's functionality within its rollup system<br>- The file's purpose is to abstract away the low-level details of contract interaction, simplifying the development process.</td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/rollups_contracts/erc20_portal.go'>erc20_portal.go</a></b></td>
						<td>- `erc20_portal.go` generates bindings for interacting with an ERC20Portal smart contract<br>- It provides functions to deposit ERC20 tokens, specifying the token, recipient application contract, amount, and execution layer data<br>- The generated code facilitates communication between Go applications and the Ethereum blockchain, enabling interaction with the ERC20Portal contract's functionalities within the larger rollups contracts package.</td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/rollups_contracts/application.go'>application.go</a></b></td>
						<td>- The `application.go` file within the `pkg/rollups_contracts` directory is a generated file containing Go bindings for interacting with the `Application` smart contract<br>- It provides the necessary structures and functions to interact with the contract from the Go application, enabling the broader project to call and monitor the contract's functionality<br>- The file's purpose is to facilitate communication between the Go application and the Ethereum blockchain, specifically the `Application` contract within the rollups system.</td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/rollups_contracts/eip712.go'>eip712.go</a></b></td>
						<td>- The `eip712.go` file within the `pkg/rollups_contracts` directory is a generated file containing the ABI (Application Binary Interface) binding for an EIP-712 compliant contract<br>- This binding facilitates interaction with the contract from the Go codebase, enabling the project to utilize the contract's functionality for structured data signing and verification, likely crucial for secure off-chain data handling within the rollup system<br>- The file's purpose is to provide a safe and efficient interface to the EIP-712 contract without requiring manual ABI handling.</td>
					</tr>
					</table>
					<details>
						<summary><b>generate</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/rollups_contracts/generate/main.go'>main.go</a></b></td>
								<td>- The `main.go` file generates Go language bindings for Cartesi Rollups contracts<br>- It downloads contract artifacts from npm, extracts relevant JSON ABI files, and uses the `bind` package to create Go code for interacting with these contracts<br>- These generated bindings are crucial for the project, enabling interaction with smart contracts from the Go application code.</td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>router</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/router/helper.go'>helper.go</a></b></td>
						<td>- `helper.go` provides a utility function for retrieving string values from a request context<br>- It facilitates accessing data stored within the context during request processing, acting as a central access point for context-based information within the routing layer of the application<br>- Failure to find a value results in an error message<br>- This function simplifies context value retrieval throughout the router package.</td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/router/router_test.go'>router_test.go</a></b></td>
						<td>- RouterSuite tests the routing functionality within the Rollmelette framework<br>- It verifies the correct handling of advance and inspect requests<br>- Tests cover successful message processing and payload extraction, ensuring the router properly routes and processes requests based on defined handlers and path parameters<br>- The suite uses a testing framework to assert expected behavior.</td>
					</tr>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/pkg/router/router.go'>router.go</a></b></td>
						<td>- `router.go` implements a routing mechanism for handling requests within the rollmelette application<br>- It registers and dispatches both "advance" and "inspect" handlers based on request paths<br>- Advance handlers process requests with payloads, while inspect handlers examine requests, extracting parameters from path patterns for context enrichment before execution<br>- The router uses regular expressions for flexible path matching.</td>
					</tr>
					</table>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- cmd Submodule -->
		<summary><b>cmd</b></summary>
		<blockquote>
			<details>
				<summary><b>tribes-rollup</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/cmd/tribes-rollup/main.go'>main.go</a></b></td>
						<td>- `main.go` serves as the entry point for the tribes-rollup command-line interface (CLI)<br>- It initializes and executes the root command, defined elsewhere in the `tribes-rollup` package<br>- The CLI's purpose within the larger Tribes project is likely to manage or interact with a rollup component, handling tasks related to data aggregation or processing<br>- Successful execution exits with code 0; otherwise, it exits with code 1, indicating an error.</td>
					</tr>
					</table>
					<details>
						<summary><b>root</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/cmd/tribes-rollup/root/wire.go'>wire.go</a></b></td>
								<td>- `wire.go` configures dependency injection within the Tribes Rollup application using the Wire library<br>- It defines sets of dependencies for repositories (user, order, crowdfunding, contract) and handlers (advance and inspect) related to Cartesi infrastructure<br>- The code facilitates the creation of middleware and handler structs, wiring together these components for application initialization.</td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/cmd/tribes-rollup/root/wire_gen.go'>wire_gen.go</a></b></td>
								<td>- Wire generates dependency injection code for the tribes-rollup command<br>- It configures middleware components (TLSN and RBAC) and handlers for advance and inspect operations across various entities (orders, users, crowdfunding, contracts)<br>- The generated code facilitates the instantiation and wiring of these components, leveraging a dependency injection framework for maintainability and testability within the larger application.</td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/cmd/tribes-rollup/root/root.go'>root.go</a></b></td>
								<td>- `root.go` constitutes the primary entry point for the tribes-rollup command-line application<br>- It initializes a database (SQLite, optionally in-memory), sets up routing for advance and inspect handlers, and runs a Rollmelette-based server<br>- The application manages various handlers for contracts, orders, crowdfunding, and users, incorporating middleware for security and access control.</td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>lib</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/cmd/tribes-rollup/lib/Cargo.toml'>Cargo.toml</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
							<details>
								<summary><b>src</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/cmd/tribes-rollup/lib/src/lib.rs'>lib.rs</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
						</blockquote>
					</details>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- test Submodule -->
		<summary><b>test</b></summary>
		<blockquote>
			<details>
				<summary><b>integration</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/test/integration/dapp_test.go'>dapp_test.go</a></b></td>
						<td><code>❯ REPLACE-ME</code></td>
					</tr>
					</table>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- internal Submodule -->
		<summary><b>internal</b></summary>
		<blockquote>
			<details>
				<summary><b>usecase</b></summary>
				<blockquote>
					<details>
						<summary><b>crowdfunding_usecase</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/delete_crowdfunding.go'>delete_crowdfunding.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/find_crowdfunding_by_id.go'>find_crowdfunding_by_id.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/create_crowdfunding.go'>create_crowdfunding.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/update_crowdfunding.go'>update_crowdfunding.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/find_crowdfundings_by_investor.go'>find_crowdfundings_by_investor.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/find_all_crowdfundings.go'>find_all_crowdfundings.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/find_crowdfunding_by_creator.go'>find_crowdfunding_by_creator.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/close_crowdfunding.go'>close_crowdfunding.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/general_dto.go'>general_dto.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/crowdfunding_usecase/settle_crowdfunding.go'>settle_crowdfunding.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>order_usecase</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/find_order_by_id.go'>find_order_by_id.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/find_orders_by_auction_id.go'>find_orders_by_auction_id.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/find_orders_by_state.go'>find_orders_by_state.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/find_orders_by_investor.go'>find_orders_by_investor.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/general_dto.go'>general_dto.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/create_order.go'>create_order.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/delete_order.go'>delete_order.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/order_usecase/find_all_orders.go'>find_all_orders.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>contract_usecase</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/contract_usecase/create_contract.go'>create_contract.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/contract_usecase/delete_contract.go'>delete_contract.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/contract_usecase/update_contract.go'>update_contract.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/contract_usecase/general_dto.go'>general_dto.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/contract_usecase/find_contract_by_symbol.go'>find_contract_by_symbol.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/contract_usecase/find_all_contracts.go'>find_all_contracts.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>user_usecase</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/find_user_by_role.go'>find_user_by_role.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/update_user.go'>update_user.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/delete_user_by_address.go'>delete_user_by_address.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/find_user_by_address.go'>find_user_by_address.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/create_user.go'>create_user.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/withdraw.go'>withdraw.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/general_dto.go'>general_dto.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/usecase/user_usecase/find_all_users.go'>find_all_users.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>domain</b></summary>
				<blockquote>
					<details>
						<summary><b>entity</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/user.go'>user.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/user_test.go'>user_test.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/order.go'>order.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/crowdfunding_test.go'>crowdfunding_test.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/crowdfunding.go'>crowdfunding.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/contract_test.go'>contract_test.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/contracts.go'>contracts.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/domain/entity/order_test.go'>order_test.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
				</blockquote>
			</details>
			<details>
				<summary><b>infra</b></summary>
				<blockquote>
					<details>
						<summary><b>repository</b></summary>
						<blockquote>
							<table>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/repository/contract_repository_sqlite.go'>contract_repository_sqlite.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/repository/order_repository_sqlite.go'>order_repository_sqlite.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/repository/user_respository_sqlite.go'>user_respository_sqlite.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							<tr>
								<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/repository/crowdfunding_repository_sqlite.go'>crowdfunding_repository_sqlite.go</a></b></td>
								<td><code>❯ REPLACE-ME</code></td>
							</tr>
							</table>
						</blockquote>
					</details>
					<details>
						<summary><b>cartesi</b></summary>
						<blockquote>
							<details>
								<summary><b>middleware</b></summary>
								<blockquote>
									<table>
									<tr>
										<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/middleware/rbac.go'>rbac.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									<tr>
										<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/middleware/tlsn.go'>tlsn.go</a></b></td>
										<td><code>❯ REPLACE-ME</code></td>
									</tr>
									</table>
								</blockquote>
							</details>
							<details>
								<summary><b>handler</b></summary>
								<blockquote>
									<details>
										<summary><b>advance_handler</b></summary>
										<blockquote>
											<table>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/advance_handler/contract_advance_handler.go'>contract_advance_handler.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/advance_handler/order_advance_handlers.go'>order_advance_handlers.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/advance_handler/user_advance_handler.go'>user_advance_handler.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/advance_handler/crowdfunding_advance_handlers.go'>crowdfunding_advance_handlers.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											</table>
										</blockquote>
									</details>
									<details>
										<summary><b>inspect_handler</b></summary>
										<blockquote>
											<table>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/inspect_handler/user_inspect_handlers.go'>user_inspect_handlers.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/inspect_handler/crowdfunding_inspect_handlers.go'>crowdfunding_inspect_handlers.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/inspect_handler/order_inspect_handlers.go'>order_inspect_handlers.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											<tr>
												<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/internal/infra/cartesi/handler/inspect_handler/contract_inspect_handler.go'>contract_inspect_handler.go</a></b></td>
												<td><code>❯ REPLACE-ME</code></td>
											</tr>
											</table>
										</blockquote>
									</details>
								</blockquote>
							</details>
						</blockquote>
					</details>
				</blockquote>
			</details>
		</blockquote>
	</details>
	<details> <!-- .github Submodule -->
		<summary><b>.github</b></summary>
		<blockquote>
			<details>
				<summary><b>workflows</b></summary>
				<blockquote>
					<table>
					<tr>
						<td><b><a href='https://github.com/henriquemarlon/tribes.git/blob/master/.github/workflows/ci.yaml'>ci.yaml</a></b></td>
						<td><code>❯ REPLACE-ME</code></td>
					</tr>
					</table>
				</blockquote>
			</details>
		</blockquote>
	</details>
</details>

---
##  Getting Started

###  Prerequisites

Before getting started with tribes.git, ensure your runtime environment meets the following requirements:

- **Programming Language:** Go
- **Package Manager:** Go modules, Cargo


###  Installation

Install tribes.git using one of the following methods:

**Build from source:**

1. Clone the tribes.git repository:
```sh
❯ git clone https://github.com/henriquemarlon/tribes.git
```

2. Navigate to the project directory:
```sh
❯ cd tribes.git
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
Run tribes.git using the following command:
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


---
##  Project Roadmap

- [X] **`Task 1`**: <strike>Implement feature one.</strike>
- [ ] **`Task 2`**: Implement feature two.
- [ ] **`Task 3`**: Implement feature three.

---

##  Contributing

- **💬 [Join the Discussions](https://github.com/henriquemarlon/tribes.git/discussions)**: Share your insights, provide feedback, or ask questions.
- **🐛 [Report Issues](https://github.com/henriquemarlon/tribes.git/issues)**: Submit bugs found or log feature requests for the `tribes.git` project.
- **💡 [Submit Pull Requests](https://github.com/henriquemarlon/tribes.git/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.

<details closed>
<summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your github account.
2. **Clone Locally**: Clone the forked repository to your local machine using a git client.
   ```sh
   git clone https://github.com/henriquemarlon/tribes.git
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m 'Implemented new feature x.'
   ```
6. **Push to github**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.
8. **Review**: Once your PR is reviewed and approved, it will be merged into the main branch. Congratulations on your contribution!
</details>

<details closed>
<summary>Contributor Graph</summary>
<br>
<p align="left">
   <a href="https://github.com{/henriquemarlon/tribes.git/}graphs/contributors">
      <img src="https://contrib.rocks/image?repo=henriquemarlon/tribes.git">
   </a>
</p>
</details>

---

##  License

This project is protected under the [SELECT-A-LICENSE](https://choosealicense.com/licenses) License. For more details, refer to the [LICENSE](https://choosealicense.com/licenses/) file.

---

##  Acknowledgments

- List any resources, contributors, inspiration, etc. here.

---
