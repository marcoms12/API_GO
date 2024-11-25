package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	

	"github.com/joho/godotenv"
)

type Customer struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	CpfCnpj       string `json:"cpfCnpj"`
	Phone         string `json:"phone"`
	MobilePhone   string `json:"mobilePhone"`
	Address       string `json:"address"`
	AddressNumber string `json:"addressNumber"`
	Complement    string `json:"complement"`
	Province      string `json:"province"`
	PostalCode    string `json:"postalCode"`
}

var userCustomer = Customer{
	Name:          "Marco Antônio Martins da Silva",
	Email:         "marcodrag00@gmail.com",
	CpfCnpj:       "11398917664",
	Phone:         "",
	MobilePhone:   "31992542059",
	Address:       "Rua Martinez Camelo de Souza",
	AddressNumber: "155",
	Complement:    "",
	Province:      "Novo Tupi",
	PostalCode:    "31846515",
}

var accessToken string

func main() {
	// Lê o arquivo .env e obtém as variáveis em um mapa
	envMap, err := godotenv.Read()
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo .env: %v", err)
	}

	// Define o token de acesso como uma variável global
	var exists bool
	accessToken, exists = envMap["ACCESS_TOKEN"]
	if !exists || accessToken == "" {
		log.Fatalf("ACCESS_TOKEN não encontrado ou está vazio no arquivo .env")
	} else {
		log.Println("Token de acesso carregado com sucesso.")
	}

	// handlers
	http.HandleFunc("/check_customer", checkCustomerHandler)
	http.HandleFunc("/create_customer", createCustomerHandler)

	fmt.Println("Servidor iniciado na porta 8080...")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

// verifica se o cliente já existe
func checkCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if accessToken == "" {
		log.Println("ACCESS_TOKEN não configurado no handler")
		http.Error(w, "Token de acesso não configurado", http.StatusInternalServerError)
		return
	}

	exists, err := checkCustomerExists(accessToken, userCustomer)
	if err != nil {
		http.Error(w, "Erro ao verificar cliente", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Cliente já existe", http.StatusConflict)
		return
	}

	fmt.Fprintln(w, "Cliente não encontrado.")
}

func createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if accessToken == "" {
		http.Error(w, "Token de acesso não configurado", http.StatusInternalServerError)
		return
	}

	exists, err := checkCustomerExists(accessToken, userCustomer)
	if err != nil {
		http.Error(w, "Erro ao verificar cliente", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Cliente já existe", http.StatusConflict)
		return
	}

	err = createCustomer(accessToken, userCustomer)
	if err != nil {
		http.Error(w, "Erro ao criar cliente", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Cliente criado com sucesso.")
}

func checkCustomerExists(accessToken string, customer Customer) (bool, error) {
	url := fmt.Sprintf("https://sandbox.asaas.com/api/v3/customers?cpfCnpj=%s", customer.CpfCnpj)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("access_token", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
		return true, nil // Cliente encontrado
	}

	return false, nil // Cliente não encontrado
}

func createCustomer(accessToken string, customer Customer) error {
	url := "https://sandbox.asaas.com/api/v3/customers"
	jsonData, err := json.Marshal(customer)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("access_token", accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("Resposta da criação do cliente:", string(body))
	return nil
}