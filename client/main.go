package client

func main() {
	host, port, filePath, fileName, err := ParseFlags()
	if err != nil {
		return
	}

}
