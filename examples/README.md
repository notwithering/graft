note: normally you would not commit the dist/ folder in a real repo, but these are examples so they are there to show the compiled output  

# custom-commands

shows how you can create custom commands in `preset` package. has very simple last-built command which just says when the file was built

# cycle

shows the `preset`'s package ability to do cycle detection, it has a cool diagram showing the cycle 

```
panic: cycle detected:
        ┌→┌─/a.txt: {{include b.txt}}
        │ └→/b.txt: {{include c.txt}}
        └─└→/c.txt: {{include a.txt}}
```

# txt-files

shows how to use the `preset` package to have one `.txt` file include a other `.txt` file with `DoubleBraceSyntax`

# website

shows how to use the `preset` package to create a website with `SGMLTagSyntax` and how to use the `include` and `extend` commands, aswell as `extend`'s subcommands, `define` and `block`
