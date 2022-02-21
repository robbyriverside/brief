# Brief Specification Format

Version 1.1.0

In short, the brief format is XML with minimal syntax and indented blocks like Python.

This repo contains a decoder for the brief format written in Go.

## Brief Example

A sample html page written in brief.

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

Further, documentation of the format below.

## Why?

The unique feature of XML having both keyword parameters and nested content body is useful when writing a specification.  However, writing XML by hand can be tedious because of the verbose syntax.  The primary design goal of brief to have the same structure as XML but easy to write.

Brief is not intended to be a data interchange format.  However, it can be easily converted to XML.  

Brief is the primary input format for the [Brevity Code Generator](https://github.com/robbyriverside/brevity).

## Brief Library

### Brief Decoder

Parses the brief format and creates a slice of Node objects.

```go
type Node struct {
    Type, Name string
    Keys       map[string]string
    Body       []*Node
    Parent     *Node
    Content    string
    Indent     int
}
```

This is more efficient than using reflection to map to an arbitrary structure and the Node object has many helpful methods for writing templates.

```go
var in io.Reader
rootNodes, err := brief.Decode(in)
```

```go
rootNodes, err := brief.DecodeFile("spec.brief")
```

Multiple top-level forms are allowed and returned as an array of Nodes by the decoder.

### Brief Encoder

Writes the Node object in brief format.

```go
var node Node
var out io.Writer
err := node.Encode(out)
```

### Brief XML Output

Writes the Node object in XML format.

```go
var node Node
var out io.Writer
err := node.WriteXML(out)
```

XML output uses a template.  This serves as an example of using brief with a template.

Contents of templates/xmlout.tmpl:

```text/template
{{define "Node"}}
{{.IndentString}}<{{.Type}}{{if .Name}} name="{{.Name}}"{{end}}{{range $key, $val := .Keys}} {{$key}}="{{$val}}"{{end}}>
{{- if .Content}}{{.Content}}{{ if not .Body}}</{{.Type}}>{{end}}{{end}}
{{- if .Body}}{{.IndentString}}{{range .Body}}{{ template "Node" . }}{{end}}
{{.IndentString}}</{{.Type}}>{{end -}}
{{end}}
{{- template "Node" .}}
```

### Template Methods

One of the primary targets of the Brief format is use in go text/templates.  There are many helpful node methods to assist in template building.

#### Node Spec

A node spec is a node type or a type:name pair.  This is used to identify a node when searching for it.

Here are some examples with an element foo with a name bar:

"foo:bar"  match both type and name.
"foo"      matches only the type without considering name.
"foo:"     matches the type and requires the name to be empty.

#### Context

Find a Node in the elements that contain this Node.  

Context will walk up the Parent hierarchy seeking a node with matches the node spec.  

```text/template
{{with .Context "project" }}
    print .Keys.id
{{end}}
```

#### Find

Find is a node method that searches the children for a node with a specific node spec.

```text/template
{{ .Find "foo" }}        // search all the children for any node of type "foo"
{{ .Find "foo:bar" }}    // search all the children for any node of type "foo" whose name is "bar"
```

#### Child

Child is a node method that follows a path to a specific child node.
The path is a series of node specs which must match each node as it walks down the children.

```text/template
{{ .Child "foo" "zed" "x" }}  // return the "x" node child of "zed" node child of "foo" node child of the current node
{{ .Child "foo:bar" }}        // return a child of the current node which is of type "foo" and named "bar"
```

#### Value Spec

A value spec is a string that can be used to locate a key value or name in a context element.

A single name, refers to the Name of the context element.  {context}.Name

A dotted pair refers to a key value from the context element. {context}.{key}

#### Lookup

Lookup is a Node method which gets a context value from a value spec.

```text/template
{{ .Lookup "project.id" }}
{{ .Lookup "project" }}
```

#### Slice

Slice is a Node method which creates a slice of strings from a sequence of value specs.

```text/template
{{ .Slice "project.id" "project" }}
```

#### Join

Join is a Node method which combines sequence of strings using a separator from a sequence of value specs.

```text/template
{{ .Join "/" "project.id" "project" }}
```

#### Printf

Printf is a Node method which applies a sequence of strings using a format from a sequence of value specs.

```text/template
{{ .Printf "%s:%s" "project.id" "project" }}
```

## Brief Format

The first token on each line is the element type.  After the element type, is a series of key-value pairs, optionally followed by a text body.  Child elements are indented on the lines below the parent element.

### Example

No better example of XML format than the widely known HTML dialect.  HTML5 has some variations, but we will skip over them for our purposes.

An HTML page contains a single top-level structure and two sub-structures:

```brief
html
    head
    body
```

Sub-structures are indented to identify a sub-block.  The first identifier on each line is an element name or type.  The back-tic contains multiline text which forms the text contents of an element.  The back-tic must be after any key-value pairs.

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

### Key Values

For key-value pairs with a value that is more than just a simple token double quotes are used.  

```brief
elem key:"value of key"  <->   <elem key="value of key"/>

elem key:"my brother \"Bill\""  <->  <elem key="my brother \"Bill\"">
```

Simple tokens cannot contain brief syntactic characters:  space, colon, back-tic, double-quote.  This allows number formats to be simple tokens.

```brief
elem size:33 max:1.4e3
```

### Name Key

Because specification elements often have a "name" keyword to identify them in the document, we give them a special place.  The element type can be a keyword by adding a colon (:) to the end and so it can also have a value, which is the name.

```brief
body
    div:foo  <->  <div name="foo"/>
```

The purpose of this shorthand is to improve readability and standardize on elements having names.  Names are important in a written specification which is the primary purpose of the brief format.

### Multi-line and Content

Line cannot start with a back-tic.  Body text requires an element.  Content back-tic is last feature on a line.
A line that starts with plus '+' is a continuation of the attributes on the line above.  May be followed by a space.

Using simple back string (or rawstring):
```brief
elem:foo bar:zed
    + more:true range:"3 to 5" `
  more content here`
```

Or using special hash delimiter (#| |#, #@ @#, #$ $#, #% %#) for when you need a nested rawstring ``
```brief
elem:foo bar:zed
    + more:true range:"3 to 5" #|
  more content here
  |#
```

### Include files

To keep files modular the #include directive allows other brief files to be inserted.  The decoder handles indentation for you so each file can be indented naturally from zero.  Include directives insert files so they can be treated like any other sub-element.

```brief
html
    head
        title `include other brief files`
        #include `standard_headers.brf`
        link rel:stylesheet href=mystyle.css
    body
        h1 `include other brief files`
```

### Comments

In the brief format, comments are treated as whitespace.

```brief
// foo element
elem:foo bar:zed
   /* multiline
      comment */
    + more:true range:"3 to 5" `
  more content here`
```
