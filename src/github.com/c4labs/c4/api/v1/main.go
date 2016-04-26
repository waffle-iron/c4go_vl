package main

import (
	"log"

	"github.com/c4labs/c4/api/v0001/client"
	"github.com/c4labs/c4/api/v0001/client/operations"
)

func main() {

  // load the swagger spec from URL or local file
   doc, err := spec.Load("https://raw.githubusercontent.com/go-swagger/go-swagger/master/examples/todo-list/swagger.yml")
   if err != nil {
     log.Fatal(err)
   }

   // create the transport
   transport := httptransport.New(doc)
   // configure the host
   if os.Getenv("TODOLIST_HOST") != "" {
     transport.Host = os.Getenv("TODOLIST_HOST")
   }

   // create the API client, with the transport
   client := apiclient.New(transport, strfmt.Default)

   // to override the host for the default client
   // apiclient.Default.SetTransport(transport)

   // make the request to get all items
   resp, err := client.Operations.All(operations.AllParams{})
   if err != nil {
     log.Fatal(err)
   }
   fmt.Printf("%#v\n", resp.Payload)


  url := "https://locahost:51483/"
  fmt.Println("URL:>", url)
	log.Print("Starting")
	c4apiClient := client.NewHTTPClient(nil)
  c4apiClient.Do(req)
	c4apiClient.Operations.FindAssets(operations.NewFindAssetsParams())
	log.Print("Ending")
}
