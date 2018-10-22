package main

func main() {
	bc, err := NewBlockchain()
	if err != nil {
		panic(err.Error())
	}
	if bc != nil {
		defer bc.db.Close()
	}

	cli := CLI{bc}
	cli.Run()
}
