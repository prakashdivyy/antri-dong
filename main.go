package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ebfe/scard"
)

func die(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func waitUntilCardPresent(ctx *scard.Context, readers []string) (int, error) {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return i, nil
			}
			rs[i].CurrentState = rs[i].EventState
		}
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
	}
}

func main() {

	// Establish a context
	ctx, err := scard.EstablishContext()
	if err != nil {
		die(err)
	}
	defer ctx.Release()

	// List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		die(err)
	}

	if len(readers) > 0 {

		index, err := waitUntilCardPresent(ctx, readers)
		if err != nil {
			die(err)
		}

		card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
		if err != nil {
			die(err)
		}

		// Select Applet Mahasiswa
		var cmd = []byte{0x00, 0xA4, 0x04, 0x00, 0x06, 0xF3, 0x60, 0x00, 0x00, 0x01, 0x01}
		rsp, err := card.Transmit(cmd)
		if err != nil {
			die(err)
		}

		// Ambil Nama
		cmd = []byte{0x90, 0x02, 0x01, 0x00, 0x00}
		rsp, err = card.Transmit(cmd)
		if err != nil {
			die(err)
		}

		// Parsing Nama
		rsp = rsp[:len(rsp)-2]
		nama := string(rsp)
		fmt.Println(nama)

		// Ambil NPM
		cmd = []byte{0x90, 0x02, 0x02, 0x00, 0x00}
		rsp, err = card.Transmit(cmd)
		if err != nil {
			die(err)
		}

		// Parsing NPM
		rsp = rsp[:len(rsp)-2]
		npm := string(rsp)
		fmt.Println(npm)

		// Ambil KODE ORG
		cmd = []byte{0x90, 0x08, 0x01, 0x00, 0x00}
		rsp, err = card.Transmit(cmd)
		if err != nil {
			die(err)
		}

		//Parsing KODE ORG
		rsp = rsp[:len(rsp)-2]
		kodeOrg := string(rsp)
		fmt.Println(kodeOrg)

		defer card.Disconnect(scard.ResetCard)

		body := strings.NewReader(`npm=` + npm + `&nama=` + nama)
		req, err := http.NewRequest("POST", "http://localhost/test.php", body)
		if err != nil {
			// handle err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}

		defer resp.Body.Close()
	}
}
