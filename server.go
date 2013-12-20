package main

import(
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
  "bufio"
  "os"
  "log"
  "io/ioutil"
  "bytes"
  "strings"
  "strconv"
  "sort"
)

type TrieNode struct {
  end int
  children map[rune]*TrieNode
  parent *TrieNode
  value rune
}

type T9Result struct {
  score int
  word string
}

type T9Results []T9Result

func NewTrieNode() *TrieNode {
  return &TrieNode{
    end: 0,
    children: make(map[rune]*TrieNode),
    parent: nil,
    value: 0,
  }
}

func TrieContains(root *TrieNode, s string) bool {
  for _, char := range s {
    if root.children[char] == nil {
      return false
    } else {
      root = root.children[char]
    }
  }
  return root.end > 0
}

func T9Words(root *TrieNode, digits string, results []T9Result) []T9Result {
  if len(digits) == 0 {
    if root.end > 0 {
      results = append(results, T9Result{root.end, WordForLeaf(root)})
    }
    return results
  }
  digit := digits[0]
  for _, char := range t9Mappings[int(digit)-48] { 
    if root.children[char] != nil {
      results = T9Words(root.children[char], digits[1:], results)
    }
  }
  return results
}

var wordTrie = NewTrieNode() 
var t9Mappings = map[int][]rune {
  2: []rune{'a','b','c'},
  3: []rune{'d','e','f'},
  4: []rune{'g','h','i'},
  5: []rune{'j','k','l'},
  6: []rune{'m','n','o'},
  7: []rune{'p','q','r','s'},
  8: []rune{'t','u','v'},
  9: []rune{'w','x','y','z'},
}

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func WordForLeaf(leaf *TrieNode) string {
  var buffer bytes.Buffer
  for {
    if leaf.parent != nil {
      buffer.WriteString(string(leaf.value))
      leaf = leaf.parent
    } else {
      break
    }
  }
  return Reverse(buffer.String())
}

func BuildTrie() {
  fi, err := os.Open("1_2_all_freq.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer fi.Close()

  scanner := bufio.NewScanner(fi)
  for scanner.Scan() {
    node := wordTrie
    wordline := strings.Fields(scanner.Text())
    for _, char := range wordline[0] {
      if node.children[char] != nil {
        node = node.children[char]
      } else {
        newChild := NewTrieNode() 
        newChild.parent = node
        newChild.value = char
        node.children[char] = newChild
        node = newChild
      }
    }
    node.end, _ = strconv.Atoi(wordline[2])
  }

  fmt.Printf("loaded")
}

// for sorting
func (results T9Results) Len() int {
  return len(results)
}

func (results T9Results) Swap(i, j int) {
  results[i], results[j] = results[j], results[i]
}

func (results T9Results) Less(i, j int) bool {
  return results[i].score > results[j].score
}

func main() {
  BuildTrie()
  rtr := mux.NewRouter()
  rtr.HandleFunc("/combinations", combinations).Methods("POST")
  rtr.HandleFunc("/", index).Methods("GET")

  http.Handle("/", rtr)
  http.ListenAndServe(":3000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
  body, _ := ioutil.ReadFile("public/index.html")
  w.Write(body)
}

func combinations(w http.ResponseWriter, r *http.Request) {
  word := r.FormValue("digits") 
  results := T9Words(wordTrie, word, make([]T9Result, 0))
  sort.Sort(T9Results(results))
  for _, result := range results {
    w.Write([]byte(fmt.Sprintf("%s: %s - %d\n", word, result.word, result.score)))
  }
}
