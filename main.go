package main

import (
	"log"

	"github.com/supabase/supabase-admin-api/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
