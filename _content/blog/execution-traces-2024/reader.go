//go:build OMIT

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/trace"
)

func main() {
	// START OMIT
	// Start reading from STDIN.
	r, err := trace.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var blocked int
	var blockedOnNetwork int
	for {
		// Read the event.
		ev, err := r.ReadEvent()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Process it.
		if ev.Kind() == trace.EventStateTransition {
			st := ev.StateTransition()
			if st.Resource.Kind == trace.ResourceGoroutine {
				from, to := st.Goroutine()

				// Look for goroutines blocking, and count them.
				if from.Executing() && to == trace.GoWaiting {
					blocked++
					if strings.Contains(st.Reason, "network") {
						blockedOnNetwork++
					}
				}
			}
		}
	}
	// Print what we found.
	p := 100 * float64(blockedOnNetwork) / float64(blocked)
	fmt.Printf("%2.3f%% instances of goroutines blocking were to block on the network\n", p)
	// END OMIT
}
