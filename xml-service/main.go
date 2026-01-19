package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

/*
ENV esperadas (no docker-compose):
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
*/

type IngestResponse struct {
	IDRequisicao string `json:"idRequisicao"`
	Status       string `json:"status"`
}

type WebhookPayload struct {
	IDRequisicao string `json:"idRequisicao"`
	Status       string `json:"status"` // OK | ERRO_VALIDACAO | ERRO_PERSISTENCIA
	DocumentID   int64  `json:"documentId,omitempty"`
	Erro         string `json:"erro,omitempty"`
}

// ---------------- XML (modelo de domínio desacoplado) ----------------

type RelatorioVendas struct {
	XMLName      xml.Name   `xml:"RelatorioVendas"`
	DataGeracao  string     `xml:"DataGeracao,attr"`
	VersaoMapper string     `xml:"VersaoMapper,attr"`
	Encomendas   Encomendas `xml:"Encomendas"`
}

type Encomendas struct {
	Lista []Encomenda `xml:"Encomenda"`
}

type Encomenda struct {
	ID        string `xml:"Id,attr"`
	Data      string `xml:"Data,attr"`
	DataEnvio string `xml:"DataEnvio,attr"`
	ModoEnvio string `xml:"ModoEnvio,attr"`

	Cliente  Cliente `xml:"Cliente"`
	Vendedor string  `xml:"Vendedor"`
	Itens    Itens   `xml:"Itens"`
}

type Cliente struct {
	ID       string      `xml:"Id,attr"`
	Segmento string      `xml:"Segmento,attr"`
	Nome     string      `xml:"Nome"`
	Local    Localizacao `xml:"Localizacao"`
}

// ✅ NOVO: Imposto no XML
type Imposto struct {
	Taxa float64 `xml:"Taxa"`
}

type Localizacao struct {
	Pais         string `xml:"Pais,attr"`
	Cidade       string `xml:"Cidade,attr"`
	Estado       string `xml:"Estado,attr"`
	CodigoPostal string `xml:"CodigoPostal,attr"`
	Regiao       string `xml:"Regiao,attr"`

	// ✅ NOVO: bloco de imposto opcional
	Imposto *Imposto `xml:"Imposto,omitempty"`
}

type Itens struct {
	Lista []Item `xml:"Item"`
}

type Item struct {
	ProdutoID    string `xml:"ProdutoId,attr"`
	Categoria    string `xml:"Categoria,attr"`
	SubCategoria string `xml:"SubCategoria,attr"`

	NomeProduto string  `xml:"NomeProduto"`
	Devolvido   bool    `xml:"Devolvido"`
	Quantidade  int     `xml:"Quantidade"`
	ValorVenda  float64 `xml:"ValorVenda"`
	Desconto    float64 `xml:"Desconto"`
	Lucro       float64 `xml:"Lucro"`
}

// ---------------- util ----------------

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Missing env %s", key)
	}
	return v
}

func openDB() *sql.DB {
	host := mustEnv("DB_HOST")
	port := mustEnv("DB_PORT")
	user := mustEnv("DB_USER")
	pass := mustEnv("DB_PASSWORD")
	name := mustEnv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}

func postWebhook(url string, payload WebhookPayload) {
	b, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Printf("Webhook falhou: %v", err)
	}
}

func normalizeHeader(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "\ufeff")
	return s
}

func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("int vazio")
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("float vazio")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func parseBoolReturned(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "yes", "true", "returned":
		return true
	case "not", "no", "false", "":
		return false
	default:
		return false
	}
}

// ---------------- parsing CSV real (com grouping por Order ID) ----------------

type requiredColumns struct {
	orderID      int
	orderDate    int
	shipDate     int
	shipMode     int
	customerID   int
	customerName int
	segment      int
	country      int
	city         int
	state        int
	postalCode   int
	region       int
	salesPerson  int
	productID    int
	category     int
	subCategory  int
	productName  int
	returned     int
	sales        int
	quantity     int
	discount     int
	profit       int
}

func buildColumnIndexMap(header []string) (requiredColumns, error) {
	idx := map[string]int{}
	for i, h := range header {
		idx[normalizeHeader(h)] = i
	}

	need := []string{
		"Order ID", "Order Date", "Ship Date", "Ship Mode",
		"Customer ID", "Customer Name", "Segment",
		"Country", "City", "State", "Postal Code", "Region",
		"Retail Sales People",
		"Product ID", "Category", "Sub-Category", "Product Name",
		"Returned", "Sales", "Quantity", "Discount", "Profit",
	}

	for _, n := range need {
		if _, ok := idx[n]; !ok {
			return requiredColumns{}, fmt.Errorf("coluna obrigatoria ausente: %q", n)
		}
	}

	return requiredColumns{
		orderID:      idx["Order ID"],
		orderDate:    idx["Order Date"],
		shipDate:     idx["Ship Date"],
		shipMode:     idx["Ship Mode"],
		customerID:   idx["Customer ID"],
		customerName: idx["Customer Name"],
		segment:      idx["Segment"],
		country:      idx["Country"],
		city:         idx["City"],
		state:        idx["State"],
		postalCode:   idx["Postal Code"],
		region:       idx["Region"],
		salesPerson:  idx["Retail Sales People"],
		productID:    idx["Product ID"],
		category:     idx["Category"],
		subCategory:  idx["Sub-Category"],
		productName:  idx["Product Name"],
		returned:     idx["Returned"],
		sales:        idx["Sales"],
		quantity:     idx["Quantity"],
		discount:     idx["Discount"],
		profit:       idx["Profit"],
	}, nil
}

// ✅ ALTERADO: recebe taxRate opcional e grava no XML em Localizacao.Imposto
func parseOrdersFromCSV(r io.Reader, taxRate *float64) ([]Encomenda, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler header: %w", err)
	}

	col, err := buildColumnIndexMap(header)
	if err != nil {
		return nil, err
	}

	orders := map[string]*Encomenda{}

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha CSV: %w", err)
		}
		if len(row) == 0 {
			continue
		}

		orderID := strings.TrimSpace(row[col.orderID])
		if orderID == "" {
			return nil, fmt.Errorf("Order ID vazio")
		}

		if _, ok := orders[orderID]; !ok {
			loc := Localizacao{
				Pais:         strings.TrimSpace(row[col.country]),
				Cidade:       strings.TrimSpace(row[col.city]),
				Estado:       strings.TrimSpace(row[col.state]),
				CodigoPostal: strings.TrimSpace(row[col.postalCode]),
				Regiao:       strings.TrimSpace(row[col.region]),
			}

			// ✅ NOVO: aplica imposto (se vier taxRate)
			if taxRate != nil {
				loc.Imposto = &Imposto{Taxa: *taxRate}
			}

			orders[orderID] = &Encomenda{
				ID:        orderID,
				Data:      strings.TrimSpace(row[col.orderDate]),
				DataEnvio: strings.TrimSpace(row[col.shipDate]),
				ModoEnvio: strings.TrimSpace(row[col.shipMode]),
				Cliente: Cliente{
					ID:       strings.TrimSpace(row[col.customerID]),
					Segmento: strings.TrimSpace(row[col.segment]),
					Nome:     strings.TrimSpace(row[col.customerName]),
					Local:    loc,
				},
				Vendedor: strings.TrimSpace(row[col.salesPerson]),
				Itens:    Itens{Lista: []Item{}},
			}
		}

		qty, err := parseInt(row[col.quantity])
		if err != nil {
			return nil, fmt.Errorf("Quantity invalido (Order %s): %w", orderID, err)
		}
		sales, err := parseFloat(row[col.sales])
		if err != nil {
			return nil, fmt.Errorf("Sales invalido (Order %s): %w", orderID, err)
		}
		discount, err := parseFloat(row[col.discount])
		if err != nil {
			return nil, fmt.Errorf("Discount invalido (Order %s): %w", orderID, err)
		}
		profit, err := parseFloat(row[col.profit])
		if err != nil {
			return nil, fmt.Errorf("Profit invalido (Order %s): %w", orderID, err)
		}

		item := Item{
			ProdutoID:    strings.TrimSpace(row[col.productID]),
			Categoria:    strings.TrimSpace(row[col.category]),
			SubCategoria: strings.TrimSpace(row[col.subCategory]),
			NomeProduto:  strings.TrimSpace(row[col.productName]),
			Devolvido:    parseBoolReturned(row[col.returned]),
			Quantidade:   qty,
			ValorVenda:   sales,
			Desconto:     discount,
			Lucro:        profit,
		}

		orders[orderID].Itens.Lista = append(orders[orderID].Itens.Lista, item)
	}

	var result []Encomenda
	for _, o := range orders {
		result = append(result, *o)
	}
	return result, nil
}

func buildXML(mapperVersion string, encomendas []Encomenda) ([]byte, error) {
	if len(encomendas) == 0 {
		return nil, fmt.Errorf("ERRO_VALIDACAO: nenhuma encomenda encontrada no CSV")
	}
	for _, e := range encomendas {
		if e.ID == "" || len(e.Itens.Lista) == 0 {
			return nil, fmt.Errorf("ERRO_VALIDACAO: encomenda invalida (id vazio ou sem itens)")
		}
	}

	doc := RelatorioVendas{
		DataGeracao:  time.Now().Format("2006-01-02"),
		VersaoMapper: mapperVersion,
		Encomendas:   Encomendas{Lista: encomendas},
	}

	out, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}

// ---------------- handlers ----------------

func main() {
	db := openDB()
	defer db.Close()

	mux := http.NewServeMux()

	// POST /ingest (multipart)
	// fields: idRequisicao, mapperVersion, webhookUrl, file (CSV), taxRate (opcional)
	mux.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(20 << 20); err != nil {
			http.Error(w, "multipart invalido", http.StatusBadRequest)
			return
		}

		idReq := r.FormValue("idRequisicao")
		mapperVersion := r.FormValue("mapperVersion")
		webhookURL := r.FormValue("webhookUrl")

		if strings.TrimSpace(idReq) == "" || strings.TrimSpace(mapperVersion) == "" || strings.TrimSpace(webhookURL) == "" {
			http.Error(w, "campos obrigatorios: idRequisicao, mapperVersion, webhookUrl", http.StatusBadRequest)
			return
		}

		// ✅ NOVO: taxRate opcional
		var taxRate *float64
		taxRateStr := strings.TrimSpace(r.FormValue("taxRate"))
		if taxRateStr != "" {
			v, err := strconv.ParseFloat(taxRateStr, 64)
			if err != nil {
				postWebhook(webhookURL, WebhookPayload{IDRequisicao: idReq, Status: "ERRO_VALIDACAO", Erro: "taxRate invalido"})
				http.Error(w, "taxRate invalido", http.StatusUnprocessableEntity)
				return
			}
			taxRate = &v
		}

		f, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "arquivo CSV (field 'file') obrigatorio", http.StatusBadRequest)
			return
		}
		defer f.Close()

		encomendas, err := parseOrdersFromCSV(f, taxRate)
		if err != nil {
			postWebhook(webhookURL, WebhookPayload{IDRequisicao: idReq, Status: "ERRO_VALIDACAO", Erro: err.Error()})
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		xmlBytes, err := buildXML(mapperVersion, encomendas)
		if err != nil {
			postWebhook(webhookURL, WebhookPayload{IDRequisicao: idReq, Status: "ERRO_VALIDACAO", Erro: err.Error()})
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		var documentID int64
		insert := `INSERT INTO relatorio (xml_documento, mapper_version) VALUES ($1::xml, $2) RETURNING id`
		if err := db.QueryRow(insert, string(xmlBytes), mapperVersion).Scan(&documentID); err != nil {
			postWebhook(webhookURL, WebhookPayload{IDRequisicao: idReq, Status: "ERRO_PERSISTENCIA", Erro: err.Error()})
			http.Error(w, "erro ao persistir XML", http.StatusInternalServerError)
			return
		}

		postWebhook(webhookURL, WebhookPayload{IDRequisicao: idReq, Status: "OK", DocumentID: documentID})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IngestResponse{IDRequisicao: idReq, Status: "RECEBIDO"})
	})

	// ---------------- XPath queries (no Postgres) ----------------

	mux.HandleFunc("/query/vendas-por-categoria", func(w http.ResponseWriter, r *http.Request) {
		categoria := strings.TrimSpace(r.URL.Query().Get("categoria"))
		if categoria == "" {
			http.Error(w, "categoria obrigatoria", http.StatusBadRequest)
			return
		}

		q := `
		WITH last_doc AS (
		  SELECT xml_documento
		  FROM relatorio
		  ORDER BY id DESC
		  LIMIT 1
		),
		vals AS (
		  SELECT unnest(xpath('//Item[@Categoria=' || quote_literal($1) || ']/ValorVenda/text()', xml_documento))::text AS v
		  FROM last_doc
		)
		SELECT COALESCE(SUM((v)::numeric), 0) AS total
		FROM vals;
		`

		var total float64
		if err := db.QueryRow(q, categoria).Scan(&total); err != nil {
			http.Error(w, "erro query: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"categoria":   categoria,
			"totalVendas": total,
		})
	})

	mux.HandleFunc("/query/lucro-por-regiao", func(w http.ResponseWriter, r *http.Request) {
		regiao := strings.TrimSpace(r.URL.Query().Get("regiao"))
		if regiao == "" {
			http.Error(w, "regiao obrigatoria", http.StatusBadRequest)
			return
		}

		q := `
		WITH last_doc AS (
		  SELECT xml_documento
		  FROM relatorio
		  ORDER BY id DESC
		  LIMIT 1
		),
		encomendas AS (
		  SELECT unnest(xpath('//Encomenda[Cliente/Localizacao[@Regiao=' || quote_literal($1) || ']]', xml_documento)) AS e
		  FROM last_doc
		),
		lucros AS (
		  SELECT unnest(xpath('.//Item/Lucro/text()', e))::text AS p
		  FROM encomendas
		)
		SELECT COALESCE(SUM((p)::numeric), 0) AS lucro_total
		FROM lucros;
		`

		var lucro float64
		if err := db.QueryRow(q, regiao).Scan(&lucro); err != nil {
			http.Error(w, "erro query: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"regiao":     regiao,
			"lucroTotal": lucro,
		})
	})

	mux.HandleFunc("/query/encomendas-prejuizo", func(w http.ResponseWriter, r *http.Request) {
		q := `
		WITH last_doc AS (
		  SELECT xml_documento
		  FROM relatorio
		  ORDER BY id DESC
		  LIMIT 1
		),
		encomendas AS (
		  SELECT unnest(xpath('//Encomenda', xml_documento)) AS e
		  FROM last_doc
		),
		por_encomenda AS (
		  SELECT
		    (xpath('string(@Id)', e))[1]::text AS order_id,
		    COALESCE((
		      SELECT SUM((x::text)::numeric)
		      FROM unnest(xpath('.//Item/Lucro/text()', e)) AS t(x)
		    ), 0) AS lucro_total
		  FROM encomendas
		)
		SELECT order_id, lucro_total
		FROM por_encomenda
		WHERE lucro_total < 0
		ORDER BY lucro_total ASC
		LIMIT 100;
		`

		rows, err := db.Query(q)
		if err != nil {
			http.Error(w, "erro query: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Row struct {
			OrderID    string  `json:"orderId"`
			LucroTotal float64 `json:"lucroTotal"`
		}
		var out []Row
		for rows.Next() {
			var id string
			var lucro float64
			if err := rows.Scan(&id, &lucro); err != nil {
				http.Error(w, "erro scan: "+err.Error(), http.StatusInternalServerError)
				return
			}
			out = append(out, Row{OrderID: id, LucroTotal: lucro})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"count":   len(out),
			"results": out,
		})
	})

	log.Println("xml-service up on :8081 (REST)")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
