### Usage
+ Entrypoint is in the file `helloMonzo.go`
+ Package `"github.com/PuerkitoBio/goquery"` is used to parse fetched HTML file

### Output
+ Crawled results are dumped into `helloMonzo.out`
+ A page identified by a unique url will have an entry `-> <index> <url> (outdegree: <count>)`
+ Under the entry described above, there is a list of url links found in the page identified by `<url>`, thus each item in the list contributes one to the value `<count>`