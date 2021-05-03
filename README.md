# brief

The brief specification format.

In short, the brief format is XML with minimal syntax and indented blocks like Python.

This repo contains readers and translators for the brief format written in Go.

## Why?

JSON syntax is fine for most specifications.  Converting to YAML or TOML makes writing the structure of JSON even easier.  But JSON is limited to objects with key-value pairs.  There is no notion of hierarchy, other than values being objects.

There are times when you want an object to have both keywords and a body.  For that, there is no substitute for XML.  In XML, elements contain key-value pairs and a body of sub-elements or text.

If we combine XML structure with the simplicity of indented blocks we get an XML structure with more brevity.

## Example

No better example of XML format than the widely known HTML dialect.  HTML5 has some variations, but we will skip over them for our purposes.

An HTML page contains a single top-level structure and two sub-structures:

```brief
html
    head
    body
```

Sub-structures are indented to identify a sub-block.  The first identifier on each line is an element name or type.  The back-tic contains multiline text which forms the contents of an element.  The back-tic must be after any key-value pairs.


```brief
html
    head
        title `My Web Page`
    body class:mybody
        h1 `My Web Page`

        div id:main class:myblock
            p id:X `the quick brown fox
jumped over the moon and ran into a cow`
```

Here is the equivalent in HTML.

```html
<html>
<head>
    <title>My Web Page</title>
</head>
<body class="mybody">
    <h1>My Web Page</h1>

    <div id="main" class="myblock">
        <p id="X">the quick brown fox
jumped over the moon and ran into a cow</p>
    </div>
</body>
</html>
```

## Key Values

For key-value pairs with a value that is more than just a simple token double quotes are used.  

```brief
elem key:"value of key"  <->   <elem key="value of key"/>

elem key:"my brother \"Bill\""  <->  <elem key="my brother \"Bill\"">
```

Simple tokens cannot contain brief syntactic characters:  space, colon, back-tic, double-quote.  This allows number formats to be simple tokens.

```brief
elem size:33 max:1.4e3
```

## Name Key

Because elements often have a "name" keyword to identify them in the document, we give them a special place.  The element type can be a keyword by adding a colon (:) to the end and so it can also have a value, which is the name.

```brief
body
    div:foo  <->  <div name="foo"/>
```

The purpose of this shorthand is to improve readability and standardize on elements having names.  Names are important in a written specification which is the primary purpose of the brief format.

## Differences from XML

Line cannot start with a back-tic.  Body text requires an element.  Content back-tic is last feature on a line.
A line that starts with plus '+' is a continuation of the attributes on the line above.  May be followed by a space.

```brief
elem:foo bar:zed
    + more:true range:"3 to 5" `
  more content here`
```

# Brief Decoder

Since brief is a specification format, it reads into a simple Node struct.


```go
type Node struct {
	Type, Name string
	Keys       map[string]string
	Body       []*Node
	Content    string
	Indent     int
}
```

This is more efficient than using reflection to map to an arbitrary structure.  

```go
root, err := brief.Decode(reader)
```

The Node structure can be walked to convert to a custom struct, if required.
