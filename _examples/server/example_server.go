package main

import (
	"os"

	"github.com/porthos-rpc/porthos-go"
	"log"
)

type input struct {
	Value int `json:"value" description:"Required"`
}

type output struct {
	Original int `json:"original_value"`
	Sum      int `json:"value_plus_one"`
}

func doSomething(req porthos.Request, res porthos.Response) {
	// nothing to do yet.
}

func doSomethingElseHandler(req porthos.Request, res porthos.Response) {
	m := make(map[string]int)
	_ = req.Bind(&m)
	log.Printf("doSomethingElse with value %f", m["value"])
}

func doSomethingThatReturnsValue(req porthos.Request, res porthos.Response) {
	var i input

	_ = req.Bind(&i)

	res.JSON(porthos.StatusOK, output{i.Value, i.Value + 1})
}

func doSomethingThatReturnsAList(req porthos.Request, res porthos.Response) {
	var i input

	_ = req.Bind(&i)

	l := make([]output, 1)
	l[0] = output{i.Value, i.Value + 1}

	res.JSON(porthos.StatusOK, l)
}

func doSomethingWithEmptyResponse(req porthos.Request, res porthos.Response) {
	var i input

	_ = req.Bind(&i)

	res.Empty(porthos.StatusOK)
}

func main() {
	b, err := porthos.NewBroker(os.Getenv("AMQP_URL"))
	defer b.Close()

	if err != nil {
		log.Print("Error creating broker")
		panic(err)
	}

	// create the RPC server.
	userService, err := porthos.NewServer(b, "UserService", porthos.Options{AutoAck: false})

	if err != nil {
		log.Print("Error creating server")
		panic(err)
	}

	defer userService.Close()

	ext, _ := porthos.NewMetricsShipperExtension(b, porthos.MetricsShipperConfig{
		BufferSize: 100,
	})

	// create and add the built-in metrics shipper.
	userService.AddExtension(ext)

	// create and add the access log extension.
	userService.AddExtension(porthos.NewAccessLogExtension())

	// create and add the specs shipper extension.
	userService.AddExtension(porthos.NewSpecShipperExtension(b))

	// dummy example procedure.
	userService.Register("doSomething", doSomething)

	// procedure with a json map spec.
	userService.RegisterWithSpec("doSomethingElse", doSomethingElseHandler, porthos.Spec{
		Description: "Here you can inform some description of your method",
		Request: porthos.ContentSpec{
			ContentType: "application/json",
			Body: porthos.BodySpecMap{
				"value": porthos.FieldSpec{Type: "float32", Description: "Required"},
			},
		},
	})

	// procedure with a json struct spec.
	userService.RegisterWithSpec("doSomethingThatReturnsValue", doSomethingThatReturnsValue, porthos.Spec{
		Request: porthos.ContentSpec{
			ContentType: "application/json",
			Body:        porthos.BodySpecFromStruct(input{}),
		},
		Response: porthos.ContentSpec{
			ContentType: "application/json",
			Body:        porthos.BodySpecFromStruct(output{}),
		},
	})

	// procedure with a json struct spec.
	userService.RegisterWithSpec("doSomethingWithEmptyResponse", doSomethingWithEmptyResponse, porthos.Spec{
		Request: porthos.ContentSpec{
			ContentType: "application/json",
			Body:        porthos.BodySpecFromStruct(input{}),
		},
	})

	// procedure with a json array spec.
	userService.RegisterWithSpec("doSomethingThatReturnsArray", doSomethingThatReturnsAList, porthos.Spec{
		Request: porthos.ContentSpec{
			ContentType: "application/json",
			Body:        porthos.BodySpecFromStruct(input{}),
		},
		Response: porthos.ContentSpec{
			ContentType: "application/json",
			Body:        porthos.BodySpecFromArray(output{}),
		},
	})

	userService.ListenAndServe()
}
