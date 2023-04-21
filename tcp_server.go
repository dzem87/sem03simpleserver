package main

import (
	"io"
	"log"
	"net"
	"sync"
	"github.com/dzem87/is105sem03/mycrypt"
	"github.com/dzem87/minyr/yr"
	
)

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}

					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekrypter melding: ", string(dekryptertMelding))
					switch msg := string(dekryptertMelding); msg {

					case "Kjevik;SN39040;18.03.2022 01:50;6":
                                                fahrMelding, err := yr.CelsiusToFahrenheitLine(msg)
                                                if err!= nil {
                                                        log.Fatal(err)
                                                }

                                                kryptertFahrMelding := mycrypt.Krypter([]rune(fahrMelding), mycrypt.ALF_SEM03, 4)
                                                _, err = c.Write([]byte(string(kryptertFahrMelding)))
						if err != nil {
							log.Fatal(err)
						}

						log.Println("Kryptert til:", string(kryptertFahrMelding))

  				        case msg:
						kryptertMelding := mycrypt.Krypter([]rune(dekryptertMelding), mycrypt.ALF_SEM03, 4)
						_, err = c.Write([]byte(string(kryptertMelding)))
						if err != nil {
							log.Fatal(err)
						}

						log.Println("Kryptert til:", kryptertMelding)

					default:
						_, err = c.Write(buf[:n])
						if err!= nil {
							log.Fatal(err)
						}
					}
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
				}
			}(conn)
		}
	}()
	wg.Wait()
}
