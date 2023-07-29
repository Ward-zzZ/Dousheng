package tools

import (
	"bufio"
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	DirtyMatchService = NewMatchService()
	AbuseWordsFile    = "../../shared/tools/abuseWords.txt"
)

// 前缀树
type Trie struct {
	End  bool           // 是否叶子节点
	Next map[rune]*Trie // 子节点
}

// 插入单词
func (n *Trie) AddWord(word string) {
	node := n
	chars := []rune(word)
	for index := range chars {
		node = node.AddChild(chars[index])
	}
	node.End = true
}

// 插入一个子节点
func (n *Trie) AddChild(c rune) *Trie {
	if n.Next == nil {
		n.Next = make(map[rune]*Trie)
	}
	if next, ok := n.Next[c]; ok {
		return next
	} else {
		n.Next[c] = &Trie{
			End:  false,
			Next: nil,
		}
		return n.Next[c]
	}
}

// 查找子节点
func (n *Trie) FindChild(c rune) *Trie {
	if n.Next == nil {
		return nil
	}

	if _, ok := n.Next[c]; ok {
		return n.Next[c]
	}
	return nil
}

type TrieTree struct {
	root *Trie
}

func NewTrieTreeMather() *TrieTree {
	return &TrieTree{
		root: &Trie{
			End: false,
		},
	}
}

// Build 构造TrieTree树
func (d *TrieTree) Build(words []string) {
	pinyinContents := HansCovertPinyin(words)
	for _, item := range words {
		d.root.AddWord(item)
	}
	// 添加拼音敏感词
	for _, item := range pinyinContents {
		d.root.AddWord(item)
	}
}

// Match 查找替换发现的敏感词
func (d *TrieTree) Match(text string, repl rune) (sensitiveWords []string, replaceText string) {
	if d.root == nil {
		return nil, text
	}

	textChars := []rune(text)
	textCharsCopy := make([]rune, len(textChars))
	copy(textCharsCopy, textChars)

	length := len(textChars)
	for i := 0; i < length; i++ {
		//root本身是没有key的，root的下面一个节点，才算是第一个；
		temp := d.root.FindChild(textChars[i])
		if temp == nil {
			continue
		}
		j := i + 1
		for ; j < length && temp != nil; j++ {
			if temp.End {
				sensitiveWords = append(sensitiveWords, string(textChars[i:j]))
				replaceRune(textCharsCopy, repl, i, j)
			}
			temp = temp.FindChild(textChars[j])
		}

		if j == length && temp != nil && temp.End {
			sensitiveWords = append(sensitiveWords, string(textChars[i:length]))
			replaceRune(textCharsCopy, repl, i, length)
		}
	}
	return sensitiveWords, string(textCharsCopy)
}

func replaceRune(chars []rune, replaceChar rune, begin int, end int) {
	for i := begin; i < end; i++ {
		chars[i] = replaceChar
	}
}

type MatchService struct {
	matcher *TrieTree
}

func NewMatchService() *MatchService {
	return &MatchService{
		matcher: &TrieTree{},
	}
}

func (m *MatchService) Build(words []string) {
	if len(words) > 0 {
		matcher := NewTrieTreeMather()
		matcher.Build(words)
		m.matcher = matcher
	}
}

func (m *MatchService) Match(text string, repl rune) (sensitiveWords []string, replaceText string) {
	sensitiveWords, replaceText = m.matcher.Match(text, repl)
	if len(sensitiveWords) > 0 {
		return
	}
	return
}

func InitDirtyMatchService() {
	words := make([]string, 0)
	file, err := os.Open(AbuseWordsFile)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n') //注意是字符
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		words = append(words, line)
	}
	DirtyMatchService.Build(words)
}

func Filter(text string) string {
	if len(text) == 0 {
		return text
	}
	if DirtyMatchService.matcher.root == nil {
		InitDirtyMatchService()
	}
	text = strings.ToLower(text)
	text = strings.Replace(text, " ", "", -1) // 去除空格

	// 过滤除中英文及数字以外的其他字符
	otherCharReg := regexp.MustCompile("[^\u4e00-\u9fa5a-zA-Z0-9]")
	text = otherCharReg.ReplaceAllString(text, "")

	_, replaceText := DirtyMatchService.Match(text, '*')
	return replaceText
}

// 中文汉字转拼音
func HansCovertPinyin(contents []string) []string {
	pinyinContents := make([]string, 0)
	for _, content := range contents {
		chineseReg := regexp.MustCompile("[\u4e00-\u9fa5]")
		if !chineseReg.Match([]byte(content)) {
			continue
		}

		// 只有中文才转
		pin := pinyin.New(content)
		pinStr, err := pin.Convert()
		if err == nil {
			pinyinContents = append(pinyinContents, pinStr)
		}
	}
	return pinyinContents
}
