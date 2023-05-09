package main

import (
	"fmt"
	"github.com/blockcypher/gobcy/v2"
	"log"
	"math/big"
	"time"
)

var bc = gobcy.API{"KEY", "btc", "main"}

func main() {
	for {
		fmt.Println("Was möchtest du tun?")
		fmt.Println("1. Geld einer Adresse überprüfen")
		fmt.Println("2. Transaktion verfolgen")
		fmt.Println("3. Transaktion überprüfen")
		fmt.Println("4. Programm beenden")
		var input int
		_, err := fmt.Scan(&input)
		if err != nil {
			log.Fatal(err)
		}

		switch input {
		case 1:
			fmt.Println("Gib die Adresse ein, von der du das Geld überprüfen möchtest:")
			var address string
			_, err := fmt.Scan(&address)
			if err != nil {
				log.Fatal(err)
			}
			balance, err := bc.GetAddrBal(address, nil)
			if err != nil {
				log.Fatal(err)
			}
			btcBalance := new(big.Float).Quo(new(big.Float).SetInt(&balance.FinalBalance), big.NewFloat(1e8))
			fmt.Printf("Der Kontostand beträgt %0.8f BTC\r\n", btcBalance)
		case 2:
			//nach der Empfängeraddresse Fragen und nach dem Betrag und dann warten bis diese Transaktion existiert und bestätigt wird
			fmt.Println("Gib die Empfängeradresse ein:")
			var address string
			_, err := fmt.Scan(&address)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Gib den Betrag (IN BTC) ein:")
			var amount int64
			_, err = fmt.Scan(&amount)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Warte auf Transaktion...\r\n")
			for {
				tx, err := bc.GetAddr(address, nil)
				if err != nil {
					log.Fatal(err)
				}
				for _, t := range tx.UnconfirmedTXRefs {
					val := new(big.Int).Set(&t.Value)
					if val.Int64() == amount {
						fmt.Println("Die Transaktion wurde am", t.Received.Format("02-01-2006"), "empfangen.")
						if t.Confirmations >= 1 {
							fmt.Println("Bestätigt ✓")
						} else {
							fmt.Println("Nicht bestätigt ✗")
						}
						break
					}
				}
				time.Sleep(5 * time.Second)
			}
		case 3:
			fmt.Println("Gib die Transaktions-ID ein, die du verfolgen möchtest:")
			var txid string
			_, err := fmt.Scan(&txid)
			if err != nil {
				log.Fatal(err)
			}
			tx, err := bc.GetTX(txid, nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Die Transaktion wurde am", tx.Received.Format("02-01-2006"), "empfangen.")
			if tx.Confirmations >= 1 {
				fmt.Println("Bestätigt ✓")
			} else {
				fmt.Println("Nicht bestätigt ✗")
			}
		case 4:
			return
		default:
			fmt.Println("Ungültige Eingabe")
		}
	}
}
