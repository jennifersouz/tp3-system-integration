package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/rpc"
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

func (s *XMLRPCServer) ProfitByRegion(args *ProfitByRegionArgs, reply *ProfitByRegionReply) error {
	if args == nil || args.Region == "" {
		return fmt.Errorf("region obrigatoria")
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
	if err := s.db.QueryRow(q, args.Region).Scan(&lucro); err != nil {
		return err
	}

	reply.Region = args.Region
	reply.LucroTotal = lucro
	return nil
}

func startXMLRPC(db *sql.DB) error {
	srv := rpc.NewServer()
	if err := srv.Register(&XMLRPCServer{db: db}); err != nil {
		return err
	}

	log.Println("XML-RPC server escutando em :8099")
	lis, err := net.Listen("tcp", ":8099")
	if err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("XML-RPC accept error: %v", err)
			continue
		}
		go srv.ServeConn(conn)
	}
}
