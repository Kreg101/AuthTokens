package main

func main() {

	s := NewServer(":8080")
	err := s.Run()
	if err != nil {
		panic(err)
	}
}
