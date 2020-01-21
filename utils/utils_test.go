package utils

import (
	"testing"
)

func TestP2PKHAddresses(t *testing.T) {
	h160 := "74d691da1574e6b3c192ecfb52cc8984ee7b6c56"
	want := "1BenRpVUFK65JFWcQSuHnJKzc4M8ZP8Eqa"
	get := H160ToP2PKH(h160, false)
	if want != get {
		t.Errorf("For h160 value of %s for mainnet, wanted the address %s but got %s", h160, want, get)
	}
	want = "mrAjisaT4LXL5MzE81sfcDYKU3wqWSvf9q"
	get = H160ToP2PKH(h160, true)
	if want != get {
		t.Errorf("For h160 value of %s for testnet, wanted the address %s but got %s", h160, want, get)
	}
}
func TestP2SHAddresses(t *testing.T) {
	h160 := "74d691da1574e6b3c192ecfb52cc8984ee7b6c56"
	want := "3CLoMMyuoDQTPRD3XYZtCvgvkadrAdvdXh"
	get := H160ToP2SH(h160, false)
	if want != get {
		t.Errorf("For h160 value of %s for mainnet, wanted the address %s but got %s", h160, want, get)
	}
	want = "2N3u1R6uwQfuobCqbCgBkpsgBxvr1tZpe7B"
	get = H160ToP2SH(h160, true)
	if want != get {
		t.Errorf("For h160 value of %s for testnet, wanted the address %s but got %s", h160, want, get)
	}
}
