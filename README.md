# brief

brief data interchange format

XML with the minimal syntax of keywords and indented blocks (like python)

This repo contains readers and translators for the brief format.

# Why?

JSON syntax is fine for most specifications.  Converting to YAML or TOML makes writing the structure of JSON even easier.  But JSON is limited to objects with key-value pairs.  There is no notion of hierarchy, other than values being objects.

There are times when you want an object to have both keywords and a body.  For that, there is no substitute for XML.  In XML, elements contain key-value pairs and a body of sub-elements or text. 

If we combine XML structure with the simplicity of indented blocks we get an XML structure with more brevity.

# Examples

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

        div id:block class:myblock
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

    <div id="block" class="myblock">
        <p>the quick brown fox
jumped over the moon and ran into a cow</p>
    </div>
</body>
</html>
```

