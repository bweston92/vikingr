# vikingr

Handles organisation repositories to make sure their protected and require pull requests.

## Installation

```
$ go install github.com/bweston92/vikingr
```

## GitHub Tokens

In order to deal with branch permissions the token needs admin rights on organisations.

## Roadmap

I currently have this project to ensure the rules I prefer, but if anyone wants to contribute a few ideas listed below.

- [ ] Different rules based on repository name etc (reads yaml file or similar)
- [ ] Ability to have different protection rules.
- [ ] Ability to run on a WebHook.
- [x] Ability to run in a periodically.
- [ ] Add metric instrumentation.

# License

Please check the LICENSE file in the repository, just a standard MIT license.
