# graft

graft is a lightweight hackable static file preprocessor  

it takes static files with include or templating commands, and outputs the same static files but with the commands resolved. so include commands are replaced with other files, extend commands will template into other files, etc. 

without graft, you might have many HTML (or other) files with repeated sections, all of which must be updated manually. with graft, you can include files into other files using commands like `<g-include header.html />`. after compilation, itll be just like you copy and pasted that file

## examples

graft is probably most useful for making websites following DRY principles. see `examples/website`

but its not only for websites, it can be used for anything static, such as txt files. see `examples/txt-files`and `examples/custom-commands`

for more examples see the `examples` directory

## notes

i havent made a full scale website to test it with yet but soon ill be starting a project which will use this, and ill probably find bugs that i will fix

btw if you have any questions, make an issue, literally nothing has documentation so if you need to know how something works then make an issue

## credits

- [mattn's go-shellwords](https://github.com/mattn/go-shellwords) - MIT
