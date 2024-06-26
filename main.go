package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
)

// The address it makes requests to
var address string = retrieveAddress()

// Allows it to make https requests to localhost mocking https
var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var client = &http.Client{Transport: transport}

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "setKey",
				Aliases: []string{"sk"},
				Usage:   "Setting the event key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "TBA Event Key",
					},
				},
				Action: func(cCtx *cli.Context) error {
					performChecks()
					KeyChangeRequest(cCtx.String("key"))
					return nil
				},
			},
			{
				Name:    "login",
				Aliases: []string{"L"},
				Usage:   "Login to the app",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u"},
						Usage:   "Username",
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
						Usage:   "Password",
					},
				},
				Action: func(cCtx *cli.Context) error {
					checkForAddress()
					sendPassword(cCtx.String("username"), cCtx.String("password"))
					return nil
				},
			},
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "Validates the server is on",
				Action: func(cCtx *cli.Context) error {
					checkForAddress()
					validateOn()
					return nil
				},
			},
			{
				Name:  "getSchedule",
				Usage: "Gets the schedule for the current selected event",
				Action: func(cCtx *cli.Context) error {
					checkForAddress()
					requestSchedule()
					return nil
				},
			},
			{
				Name:  "update-address",
				Usage: "Updates the address the CLI will go to.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"a"},
						Usage:   "The new address",
					},
				},
				Action: func(cCtx *cli.Context) error {
					updateAddress(cCtx.String("address"))
					return nil
				},
			},
			{
				Name:  "update-sheet",
				Usage: "Updates the sheet the backend will use.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "sheet",
						Aliases: []string{"s"},
						Usage:   "The new sheet ID",
					},
				},
				Action: func(cCtx *cli.Context) error {
					performChecks()
					updateSheet(cCtx.String("sheet"))
					return nil
				},
			},
			{
				Name:  "getScouterSchedule",
				Usage: "Gets the schedule of a single scouter.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "scouter",
						Usage: "The scouter's name",
					},
				},
				Action: func(cCtx *cli.Context) error {
					performChecks()
					getScouterSchedule(cCtx.String("scouter"))
					return nil
				},
			},
			{
				Name:  "getLeaderboard",
				Usage: "Gets the leaderboard from the backend",
				Action: func(cCtx *cli.Context) error {
					performChecks()
					getLeaderboard()
					return nil
				},
			},
			{
				Name:  "getAddress",
				Usage: "Gets the cached address",
				Action: func(cCtx *cli.Context) error {
					println(retrieveAddress())
					return nil
				},
			},
			{
				Name:  "getUsers",
				Usage: "Gets the list of all users",
				Action: func(cCtx *cli.Context) error {
					performChecks()
					getUsers()
					return nil
				},
			},
			{
				Name:  "genPassword",
				Usage: "Generates a new password using bcrypt",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p"},
						Usage:   "Name of scouter to modify",
					},
				},
				Action: func(cCtx *cli.Context) error {
					encrypted, _ := bcrypt.GenerateFromPassword([]byte(cCtx.String("password")), 6)
					fmt.Println(string(encrypted))
					return nil
				},
			},

			{
				Name:  "modify-leaderboard",
				Usage: "Modify the leaderboard",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Name of scouter to modify",
					},
					&cli.StringFlag{
						Name:    "Modification",
						Aliases: []string{"m"},
						Usage:   "The type of modification: Increase, Decrease, or Set",
					},
					&cli.IntFlag{
						Name:    "By",
						Aliases: []string{"b"},
						Usage:   "How much to modify by",
					},
				},
				Action: func(cCtx *cli.Context) error {
					performChecks()
					modifyLeaderboard(cCtx.String("name"), Modification(cCtx.String("Modification")), cCtx.Int("By"))
					return nil
				},
			},
			{
				Name:  "addBadge",
				Usage: "Add a badge",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Name of scouter to modify",
					},
					&cli.StringFlag{
						Name:    "badge",
						Aliases: []string{"b"},
						Usage:   "The Badge name to add",
					},
					&cli.StringFlag{
						Name:    "description",
						Aliases: []string{"d"},
						Usage:   "The badge description",
					},
				},
				Action: func(cCtx *cli.Context) error {
					performChecks()
					addBadge(cCtx.String("name"), cCtx.String("badge"), cCtx.String("description"))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

// Checks for a valid address and certificate
func performChecks() {
	checkForAddress()
	checkForValidCert()
}

// Sends a request to change the event key
func KeyChangeRequest(newKey string) {
	keyBytes := []byte(newKey)

	body := bytes.NewReader(keyBytes)
	request, err := http.NewRequest("GET", address+"/keyChange", body)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Add("Certificate", retrieveCredentials().Certificate)

	resp, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}
	newBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(newBody)
	log.Print(sb)
}

// Validates the certificate
func validateInternally() bool {
	request, _ := http.NewRequest("GET", address+"/", bytes.NewBufferString(""))
	request.Header.Add("Certificate", retrieveCredentials().Certificate)

	resp, _ := client.Do(request)

	return resp != nil
}

// Validates the server is on
func validateOn() {
	if validateInternally() {
		println("Server validated to be on!")
	} else {
		fmt.Printf("Server offline. Please make sure %v is the right address.", address)
	}
}

// Sends a request for the event schedule
func requestSchedule() {
	response, _ := client.Get(address + "/schedule")

	if response == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(response.Body)

		sb := string(newBody)
		log.Print(sb)
	}

}

// Logs in
func sendPassword(username string, password string) {
	pub := getPublicKey()

	encryptedPassword, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(password))
	request := LoginRequest{
		Username:          username,
		EncryptedPassword: encryptedPassword,
	}

	jsonLogin, _ := json.Marshal(request)

	response, err := client.Post(address+"/login", "application/json", bytes.NewBuffer(jsonLogin))

	if err != nil {
		log.Fatalln(err)
	}

	newBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(newBody)
	log.Println(sb)

	saveCredentials(response.Header.Get("uuid"), response.Header.Get("certificate"))
}

// Gets the RSA public key from the server
func getPublicKey() *rsa.PublicKey {
	response, _ := client.Get(address + "/pub")

	newBody, _ := io.ReadAll(response.Body)

	block, _ := pem.Decode(newBody)

	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)

	return key
}

// A request to log in
type LoginRequest struct {
	Username          string // The username
	EncryptedPassword []byte // The password, encrypted with the RSA public key
}

// User credentials
type Credentials struct {
	UUID        string // The uuid
	Certificate string // The certificate
}

// The user configuration directory
var appdata, _ = os.UserConfigDir()

// The path to the GreenScoutCLI config directory
var configDir = filepath.Join(appdata, "GreenScoutCLI")

// Saves credentials to the file system
func saveCredentials(uuid string, certificate string) {
	os.Mkdir(configDir, os.ModePerm) // Not checking error here bcs the only real error is it alr existing

	credentialPath := filepath.Join(configDir, "credentials.json")

	file, _ := os.Create(credentialPath)
	defer file.Close()

	myCreds := Credentials{UUID: uuid, Certificate: certificate}

	json.NewEncoder(file).Encode(myCreds)
}

// Gets the credentials from the file system
func retrieveCredentials() Credentials {
	file, _ := os.Open(filepath.Join(configDir, "credentials.json"))
	defer file.Close()

	var credentials Credentials

	json.NewDecoder(file).Decode(&credentials)

	return credentials
}

// Updates the address it attempts to make requests to
func updateAddress(newAddress string) {
	os.Mkdir(configDir, os.ModePerm) // Not checking error here bcs the only real error is it alr existing

	addressPath := filepath.Join(configDir, "address.txt")

	file, _ := os.Create(addressPath)
	defer file.Close()

	file.WriteString(newAddress)

	address = newAddress
}

// Retrieves the address to make requests to
func retrieveAddress() string {
	file, _ := os.Open(filepath.Join(configDir, "address.txt"))
	defer file.Close()

	dataBytes, _ := io.ReadAll(file)

	return string(dataBytes)
}

// Sends a sheet update request
func updateSheet(newSheet string) {
	response, _ := client.Post(address+"/sheetChange", "text/plain", bytes.NewBufferString(newSheet))

	if response == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(response.Body)

		sb := string(newBody)
		log.Println("Server Returned: " + sb)
	}
}

// Gets the schedule of a specific scouter
func getScouterSchedule(name string) {
	request, _ := http.NewRequest("GET", address+"/singleSchedule", bytes.NewBufferString(""))
	request.Header.Add("Certificate", retrieveCredentials().Certificate)
	request.Header.Add("userInput", name)

	resp, _ := client.Do(request)

	if resp == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(resp.Body)

		sb := string(newBody)
		log.Println("Server Returned: " + sb)
	}
}

// A leaderboard modification request
type ModRequest struct {
	Name string
	By   int
	Mod  Modification
}

// A badge
type Badge struct {
	ID          string
	Description string
}

type Modification string

const ( // Modification type enum
	Increase Modification = "Increase"
	Decrease Modification = "Decrease"
	Set      Modification = "Set"
)

// Gets the leaderboard from the server
func getLeaderboard() {
	response, _ := client.Get(address + "/leaderboard")

	if response == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(response.Body)

		sb := string(newBody)
		log.Print(sb)
	}
}

// Sends a request to modify the leaderboard
func modifyLeaderboard(name string, mod Modification, by int) {
	jsonBytes, _ := json.Marshal(ModRequest{Name: name, Mod: mod, By: by})
	request, _ := http.NewRequest("POST", address+"/modScore", bytes.NewBuffer(jsonBytes))
	request.Header.Add("Certificate", retrieveCredentials().Certificate)

	resp, _ := client.Do(request)

	if resp == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(resp.Body)

		sb := string(newBody)
		log.Println("Server Returned: " + sb)
	}
}

// Adds a badge to the leaderboard
func addBadge(name string, badge string, description string) {
	jsonBytes, _ := json.Marshal(Badge{ID: badge, Description: description})

	request, _ := http.NewRequest("POST", address+"/addBadge", bytes.NewBuffer(jsonBytes))
	request.Header.Add("Certificate", retrieveCredentials().Certificate)
	request.Header.Add("Username", name)

	resp, _ := client.Do(request)

	if resp == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(resp.Body)

		sb := string(newBody)
		log.Println("Server Returned: " + sb)
	}
}

// Gets all users from the server
func getUsers() {
	request, err := http.NewRequest("GET", address+"/allUsers", bytes.NewBufferString(""))
	request.Header.Add("Certificate", retrieveCredentials().Certificate)
	if err != nil {
		log.Fatalln(err)
	}

	response, _ := client.Do(request)

	if response == nil {
		log.Println("Server did not return a response.")
	} else {
		newBody, _ := io.ReadAll(response.Body)

		sb := string(newBody)
		log.Print(sb)
	}
}

// Ensures the address is not null
func checkForAddress() {
	if address == "" {
		log.Fatalln("Please enter an address for the server using ./GreenScoutCLI update-address", address)
	}
}

// Checks the certificate is valid
func checkForValidCert() {
	request, err := http.NewRequest("GET", address+"/certificateValid", bytes.NewBufferString(""))
	request.Header.Add("Certificate", retrieveCredentials().Certificate)
	if err != nil {
		log.Fatalln(err)
	}

	response, _ := client.Do(request)

	if response.StatusCode != 200 {
		if validateInternally() {
			log.Fatalln("Certificate Invalid. Please log in with ./GreenScoutCLI login")
		} else {
			log.Fatalln("Server offline. Please make sure " + address + " is the right address.")
		}
	}
}
