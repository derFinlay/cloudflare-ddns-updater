package main

func main() {
	
}

func loadConfig() Config {

}

resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
if err != nil {
   log.Fatalln(err)
}