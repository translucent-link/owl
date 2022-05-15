package cmd

import (
	"log"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/translucent-link/owl/graph"
	"github.com/translucent-link/owl/graph/generated"
	"github.com/urfave/cli/v2"
)

func server(c *cli.Context) error {
	port := c.Int("port")
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

var ServerCommand = &cli.Command{
	Name:   "server",
	Usage:  "runs server process",
	Action: server,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Usage: "what port the server runs on",
		},
	},
}
