# brief

The brief specification format.

In short, the brief format is XML with minimal syntax and indented blocks like Python.

This repo contains readers and translators for the brief format written in Go.

# Why?

JSON syntax is fine for most specifications.  Converting to YAML or TOML makes writing the structure of JSON even easier.  But JSON is limited to objects with key-value pairs.  There is no notion of hierarchy, other than values being objects.

There are times when you want an object to have both keywords and a body.  For that, there is no substitute for XML.  In XML, elements contain key-value pairs and a body of sub-elements or text. 

If we combine XML structure with the simplicity of indented blocks we get an XML structure with more brevity.

# Example

No better example of XML format than the widely known HTML dialect.  HTML5 has some variations, but we will skip over them for our purposes.

An HTML page contains a single top-level structure and two sub-structures:

```
html
    head
    body
```

Sub-structures are indented to identify a sub-block.  The first identifier on each line is an element name or type.  The back-tic contains multiline text which forms the contents of an element.


```
html
    head
        title `My Web Page`
    body class:mybody
        h1 `My Web Page`

        div id:main class:myblock
            p `the quick brown fox
jumped over the moon and ran into a cow`
```

Here is the equivalent in HTML.

```
<html>
<head>
    <title>My Web Page</title>
</head>
<body class="mybody">
    <h1>My Web Page</h1>

    <div id="main" class="myblock">
        <p>the quick brown fox
jumped over the moon and ran into a cow</p>
    </div>
</body>
</html>
```

# Key Values

For key-value pairs with a value that is more than just a simple token double quotes are used.  Also, back-slash can be used to escape nested quotes or back-tics.

```
elem key:"value of key"  <->   <elem key="value of key"/>

elem key:"my brother \"Bill\""  <->  <elem key="my brother \"Bill\"">

title `time is an \`illusion\`` <->  <title>time is an `illusion`</title>
```

Simple tokens cannot contain brief syntactic characters:  space, colon, back-tic, double-quote.  This allows number formats to be simple tokens.

```
elem size:33 max:1.4e3
```

# Name Key

Because elements often have a "name" keyword to identify them in the document, we give them a special place.  The element type can be a keyword by adding a colon (:) to the end and so it can also have a value, which is the name.

```
body
    div:foo  <->  <div name="foo"/>
```

The purpose of this shorthand is to improve readability and standardize on elements having names.  Names are important in a written specification which is the primary purpose of the brief format.
