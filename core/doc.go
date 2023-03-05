// Package core can be used to obtain a config including credentials from different sources,
// and provide a simplified abstraction to handle events and commands.
//
// You just need to simply create a new instance of core.Configuration:
//
// var config = core.NewConfiguration("app.brightsec.com")
//
// Configuration can be customized, for instance, you can set a credentials manually:
//
// var cred, _ = credentials.Credentials("your API key")
// var config = core.NewConfiguration("app.brightsec.com", WithCredentials(cred))
package core
