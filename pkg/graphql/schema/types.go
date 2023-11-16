package graphql

import (
	"github.com/graphql-go/graphql"
)

var artistType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Artist",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "location": &graphql.Field{
                Type: graphql.String,
            },
            "soundcloudSetLink": &graphql.Field{
                Type: graphql.String,
            },
            // Note: Relationships like SocialMediaLinks might be handled differently
            // depending on your query and data fetching logic
        },
    },
)

var socialMediaType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "SocialMedia",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "platform": &graphql.Field{
                Type: graphql.String,
            },
            "link": &graphql.Field{
                Type: graphql.String,
            },
        },
    },
)
