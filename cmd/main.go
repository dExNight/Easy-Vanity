package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"ton-generator/internal/generator"
)

func main() {
	suffix := flag.String("suffix", "", "Desired address suffix (2-4 characters)")
	caseSense := flag.Bool("case", true, "Case sensitive matching")
	workchain := flag.Int("workchain", 0, "Workchain ID (0 or -1)")
	flag.Parse()

	if len(*suffix) < 2 || len(*suffix) > 5 {
		fmt.Println("Error: suffix must be 2-4 characters long")
		return
	}

	if *workchain != 0 && *workchain != -1 {
		fmt.Println("Error: workchain must be 0 or -1")
		return
	}

	config := generator.GeneratorConfig{
		Suffix:    *suffix,
		CaseSense: *caseSense,
		Workchain: *workchain,
	}

	fmt.Printf("Searching for address with suffix '%s' (case sensitive: %v) in workchain %d\n",
		config.Suffix, config.CaseSense, config.Workchain)

	result := make(chan *generator.WalletData)
	stop := make(chan struct{})
	var attempts uint64

	numCPU := runtime.NumCPU()
	fmt.Println("Using", numCPU, "CPU cores")
	var wg sync.WaitGroup

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					wallet, err := generator.GenerateWalletData(config.Workchain)
					if err != nil {
						continue
					}

					atomic.AddUint64(&attempts, 1)
					if attempts%100 == 0 {
						fmt.Printf("Tried %d addresses...\n", attempts)
					}

					if generator.CheckAddressSuffix(wallet.Address, config) {
						result <- wallet
						close(stop)
						return
					}
				}
			}
		}()
	}

	wallet := <-result
	fmt.Printf("\nFound matching address after %d attempts!\n", attempts)
	fmt.Printf("Address: %s\n", wallet.Address)
	fmt.Printf("Mnemonic Phrase: %s\n", wallet.Mnemonic)
	fmt.Printf("Private Key: %s\n", wallet.PrivateKey)
	fmt.Printf("Public Key: %s\n", wallet.PublicKey)

	wg.Wait()
}
