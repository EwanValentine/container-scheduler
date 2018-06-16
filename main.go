package main

func main() {
	wait := make(chan bool)

	<-wait
}
