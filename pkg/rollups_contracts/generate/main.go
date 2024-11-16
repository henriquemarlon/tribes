package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const (
	openzeppelin        = "https://registry.npmjs.org/@openzeppelin/contracts/-/contracts-5.0.2.tgz"
	rollupsContractsUrl = "https://registry.npmjs.org/@cartesi/rollups/-/rollups-2.0.0-rc.10.tgz"
	baseContractsPath   = "package/export/artifacts/contracts/"
	bindingPkg          = "rollups_contracts"
)

type contractBinding struct {
	jsonPath string
	typeName string
	outFile  string
}

var bindings = []contractBinding{
	{
		jsonPath: "package/build/contracts/EIP712.json",
		typeName: "EIP712",
		outFile:  "./pkg/rollups_contracts/eip712.go",
	},
	{
		jsonPath: baseContractsPath + "inputs/InputBox.sol/InputBox.json",
		typeName: "InputBox",
		outFile:  "./pkg/rollups_contracts/input_box.go",
	},
	{
		jsonPath: baseContractsPath + "dapp/Application.sol/Application.json",
		typeName: "Application",
		outFile:  "./pkg/rollups_contracts/application.go",
	},
	{
		jsonPath: baseContractsPath + "portals/ERC20Portal.sol/ERC20Portal.json",
		typeName: "ERC20Portal",
		outFile:  "./pkg/rollups_contracts/erc20_portal.go",
	},
}

func main() {
	// Configurar logs detalhados
	slog.Info("Starting contract bindings generation")

	// Baixar e descompactar pacotes de contratos
	contractsZip, err := downloadContracts(rollupsContractsUrl)
	checkErr("download contracts", err)
	defer contractsZip.Close()

	contractsTar, err := unzip(contractsZip)
	checkErr("unzip contracts", err)
	defer contractsTar.Close()

	contractsOpenZeppelin, err := downloadContracts(openzeppelin)
	checkErr("download OpenZeppelin contracts", err)
	defer contractsOpenZeppelin.Close()

	contractsTarOpenZeppelin, err := unzip(contractsOpenZeppelin)
	checkErr("unzip OpenZeppelin contracts", err)
	defer contractsTarOpenZeppelin.Close()

	// Mapear arquivos necessários
	files := make(map[string]bool)
	for _, b := range bindings {
		files[b.jsonPath] = true
	}

	contents, err := readFilesFromTar(contractsTar, files)
	checkErr("read files from tar (Rollups)", err)

	contentsZ, err := readFilesFromTar(contractsTarOpenZeppelin, files)
	checkErr("read files from tar (OpenZeppelin)", err)

	// Mesclar conteúdos de ambos os pacotes
	for key, content := range contentsZ {
		contents[key] = content
	}

	// Gerar bindings para cada contrato
	for _, b := range bindings {
		content := contents[b.jsonPath]
		if content == nil {
			log.Fatalf("missing contents for %s", b.jsonPath)
		}
		generateBinding(b, content)
	}

	slog.Info("Contract bindings generation completed successfully")
}

// Exit if there is any error.
func checkErr(context string, err any) {
	if err != nil {
		log.Fatalf("%s: %v", context, err)
	}
}

// Download the contracts from the provided URL.
func downloadContracts(url string) (io.ReadCloser, error) {
	slog.Info("Downloading contracts", slog.String("url", url))
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download contracts from %s: %w", url, err)
	}
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		return nil, fmt.Errorf("failed to download contracts from %s: status code %s", url, response.Status)
	}
	return response.Body, nil
}

// Decompress the buffer with the contracts.
func unzip(r io.Reader) (io.ReadCloser, error) {
	slog.Info("Unzipping contracts")
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip: %w", err)
	}
	return gzipReader, nil
}

// Read the required files from the tar archive.
func readFilesFromTar(r io.Reader, files map[string]bool) (map[string][]byte, error) {
	contents := make(map[string][]byte)
	tarReader := tar.NewReader(r)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, fmt.Errorf("error while reading tar: %w", err)
		}

		if files[header.Name] {
			slog.Info("Found file in tar", slog.String("name", header.Name))
			contents[header.Name], err = io.ReadAll(tarReader)
			if err != nil {
				return nil, fmt.Errorf("error while reading file inside tar: %w", err)
			}
		}
	}
	return contents, nil
}

// Extract the ABI from the contract JSON.
func getAbi(rawJson []byte) []byte {
	var contents struct {
		Abi json.RawMessage `json:"abi"`
	}
	err := json.Unmarshal(rawJson, &contents)
	checkErr("decode JSON", err)
	return contents.Abi
}

// Generate the Go bindings for the contracts.
func generateBinding(b contractBinding, content []byte) {
	// Ensure the output directory exists
	err := os.MkdirAll("./pkg/rollups_contracts", 0755)
	checkErr("create output directory", err)

	// Prepare binding parameters
	abi := getAbi(content)
	var (
		sigs    []map[string]string
		abis    = []string{string(abi)}
		bins    = []string{""}
		types   = []string{b.typeName}
		libs    = make(map[string]string)
		aliases = make(map[string]string)
	)

	// Generate bindings
	code, err := bind.Bind(types, abis, bins, sigs, bindingPkg, bind.LangGo, libs, aliases)
	checkErr("generate binding", err)

	// Write binding to file
	err = os.WriteFile(b.outFile, []byte(code), 0600)
	checkErr("write binding file", err)
	slog.Info("Generated binding", slog.String("file", b.outFile))
}
