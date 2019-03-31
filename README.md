# truncate

A Go package that contains just one function that will truncate a given byte slice to a maximum of `maxlen` visible characters and optionally append a string before closing any open tags (e.g. for an ellipsis). HTML tags are automatically closed generating valid HTML.

## Usage

Import the package in your Go program.

    import (
        "github.com/andybons/truncate"
    )

Call `HTML` passing in the byte slice, the max len, and the string that should be appended to the truncated HTML before closing the open tags.

    func HTML(buf []byte, maxlen int, ellipsis string) ([]byte, error)

## License

The MIT license.
