# GO=GZIP64 #

A little golang console utility that reads a file and either:

    1) Encodes it - gzip compress followed by base64 encode
                    writes the result to stdout + <filename>.encoded

    2) Decodes it - base64 decode followed by gzip decompress
                    writes  the result to stdout + <filename>.decoded

