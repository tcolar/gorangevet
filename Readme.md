# GoRangeVet

Usage :`gorangevet` or `gorangevet <packages>`

`gorangevet` checks the inside of range loops to see if a pointer to the key or value is taken, which can lead to some
very nasty bugs.

The standard `go tool vet` does check that range loop key and values are not used in functions or routines, however it
does not check whether a pointer to them is taken.

There was a proposal for that, see https://github.com/golang/go/issues/20725

It seem it will not be made part of `go tool vet`, the opinion is that there may be legitimate
case where doing so may be useful.

I haven't come up with a useful use case, on the other had I have several times
come across code where it was done by accident and caused some nasty bugs.

Note that the proposal issue mentions that it might become disallowed altogether in Go2, which at least admits it's probably a bad idea.

Example problematic code that will be reported :
```
	vals := []int{1, 2, 3, 4, 5, 6, 7}
	results := []*int{}
	for _, v := range vals {
		results = append(results, &v)
	}

	for _, r := range results {
	    fmt.Print(*r)
	}
```
Output:
`7 7 7 7 7 7 7`

https://play.golang.org/p/EK0W4iy7iu
