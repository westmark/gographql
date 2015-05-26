/*
Package gographql implements a parser and a scanner for Facebook's query language GraphQL.

(see https://facebook.github.io/react/blog/2015/05/01/graphql-introduction.html)

Given the example input

  s := `
  {
    foo,
    bar,
    user(id:1) {
      name,
      age
    },
    store(address:"First Street", zip: "1337", active: true) {
      city
    }
  }
  `
a new Parser can be created

  block, err := gographql.NewParser(strings.NewReader(s)).Parse()

The return value `block` can then be inspected to discover the structure of the query.

As GraphQL is not yet finalized, expect heavy changes to this project.
*/
package gographql
