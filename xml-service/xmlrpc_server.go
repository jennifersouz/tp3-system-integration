package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type XMLRPCServer struct {
	db *sql.DB
}

type ProfitByRegionArgs struct {
	Region string `json:"region"`
}

type ProfitByRegionReply struct {
	Region     string  `json:"region"`
	LucroTotal float64 `json:"lucroTotal"`
}

// Manual XML-RPC handler
func (s *XMLRPCServer) handleXMLRPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse XML-RPC request
	bodyStr := string(body)
	log.Printf("[XML-RPC] Request: %s", bodyStr)

	if strings.Contains(bodyStr, "XMLRPCServer.ProfitByRegion") {
		// Extract region parameter
		startIdx := strings.Index(bodyStr, "<string>")
		endIdx := strings.Index(bodyStr, "</string>")

		if startIdx == -1 || endIdx == -1 {
			http.Error(w, "Invalid parameter", http.StatusBadRequest)
			return
		}

		region := bodyStr[startIdx+8 : endIdx]
		log.Printf("[XML-RPC] Recebendo ProfitByRegion para região: %s", region)

		// Query database
		q := `
		WITH encomendas AS (
		  SELECT unnest(xpath('//Encomenda[Cliente/Localizacao[@Regiao=''' || $1 || ''']]', xml_documento)) AS e
		  FROM relatorio
		),
		items AS (
		  SELECT unnest(xpath('.//Item', e))::xml AS item_elem
		  FROM encomendas
		),
		lucros AS (
		  SELECT (xpath('.//Lucro/text()', item_elem))[1]::text::numeric AS lucro
		  FROM items
		)
		SELECT COALESCE(SUM(lucro), 0) AS lucro_total
		FROM lucros;
		`

		var lucro float64
		if err := s.db.QueryRow(q, region).Scan(&lucro); err != nil {
			log.Printf("[XML-RPC] Erro na query: %v", err)
			responseError(w, "Database error")
			return
		}

		log.Printf("[XML-RPC] Lucro total para %s: %.2f", region, lucro)

		// Build XML-RPC response with struct
		buf := new(bytes.Buffer)
		xml.EscapeText(buf, []byte(region))
		escapedRegion := buf.String()

		response := fmt.Sprintf(`<?xml version="1.0"?>
<methodResponse>
<params>
<param>
<value>
<struct>
<member>
<name>region</name>
<value><string>%s</string></value>
</member>
<member>
<name>lucroTotal</name>
<value><double>%f</double></value>
</member>
</struct>
</value>
</param>
</params>
</methodResponse>`, escapedRegion, lucro)

		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, response)
		return
	}

	http.Error(w, "Unknown method", http.StatusBadRequest)
}

func responseError(w http.ResponseWriter, msg string) {
	response := fmt.Sprintf(`<?xml version="1.0"?>
<methodResponse>
<fault>
<value>
<struct>
<member>
<name>faultCode</name>
<value><int>1</int></value>
</member>
<member>
<name>faultString</name>
<value><string>%s</string></value>
</member>
</struct>
</value>
</fault>
</methodResponse>`, msg)

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func startXMLRPC(db *sql.DB) error {
	server := &XMLRPCServer{db: db}
	http.HandleFunc("/RPC2", server.handleXMLRPC)

	log.Println("XML-RPC server escutando em :8099/RPC2")
	return http.ListenAndServe(":8099", nil)
}
