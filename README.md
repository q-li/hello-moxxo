### Build
+ Entrypoint is in the file `hello-monzo.go`
+ Package `"github.com/PuerkitoBio/goquery"` is used to parse fetched HTML file

### Rules
+ Try to find URLs from `<a>` tags' `href` attribute
+ Other possibilites could be `href` attribute of `<link>` tag, `src` attribute of `<img>` and `script` tags, `content` attribute of `<meta>` tag
+ Examines scheme of found "URL"s, to rule out several abnormalities

### Output
+ Crawled results are dumped into `hello-monzo.out`
+ A page identified by a unique url will have an entry `-> <index> <url> (outdegree: <count>)`
+ Under the entry described above, there is a list of url links found in the page identified by `<url>`, thus each item in the list contributes one to the value `<count>`