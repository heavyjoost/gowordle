package gowordle

import (
    "bufio"
    "fmt"
    "github.com/go-chat-bot/bot"
    "log"
    "math/rand"
    "os"
    "regexp"
    "runtime/debug"
    "strings"
)

const (
    // Pattern matches 5-letter lowercase words (to avoid proper nouns). To use 4 and 5 letter words, change {5} to {4,5}.
    pattern = `^[a-z]{5}[\r\n]*$`
    wordsPathVar = `GOWORDLE_WORDS_PATH`
)

var (
    re = regexp.MustCompile(pattern)
    words []string
    idx = 0
)

type ircformatting string
const (
    Reset ircformatting = "\x0f"
    Bold                = "\x02"
    Italic              = "\x1d"
    Underline           = "\x1f"
    Strikethrough       = "\x1e" // Might not be supported by client
    Monospace           = "\x11" // Might not be supported by client
    ReverseColor        = "\x16" // Might not be supported by client
    Color               = "\x03"
)

type irccolor string
const (
    Undefined irccolor = ""
    White              = "00"
    Black              = "01"
    Blue               = "02"
    Green              = "03"
    Red                = "04"
    Brown              = "05"
    Magenta            = "06"
    Purple             = "06"
    Orange             = "07"
    Yellow             = "08"
    LightGreen         = "09"
    Cyan               = "10"
    Teal               = "10"
    LightCyan          = "11"
    LightBlue          = "12"
    Pink               = "13"
    Grey               = "14"
    Gray               = "14"
    LightGrey          = "15"
    LightGray          = "15"
);

func die(err error) {
    if err != nil {
        // Only print stack if DEBUG variable has been defined
        if os.Getenv("DEBUG") != "" {
            debug.PrintStack()
        }
        log.Fatal(err)
    }
}

func colorify(text string, fg irccolor, bg irccolor) string {
    if bg == Undefined {
        return fmt.Sprintf("%s%s%s%s", Color, fg, text, Reset)
    } else {
        return fmt.Sprintf("%s%s,%s%s%s", Color, fg, bg, text, Reset)
    }
}

func gowordle(command *bot.PassiveCmd) (string, error) {
    input := strings.TrimRight(command.MessageData.Text, " \r\n")
    
    // We only care about what somebody says if it matches the length of the word
    if len(input) != len(words[idx]) {
        return "", nil
    }
    
    // User guesses word we're looking for
    if strings.EqualFold(input, words[idx]) {
        idx = (idx + 1) % len(words)
        return fmt.Sprintf("Congrats %s, %s is correct!\n%s", command.User.Nick, colorify(input, Black, Green), strings.Repeat("_ ", len(words[idx]))), nil
    }
    
    // Do complicated stuff when somebody makes a guess but it's wrong
    word := []rune(strings.ToLower(words[idx]))
    entry := []rune(strings.ToLower(input))
    output := make([]string, len(word))
    
    // Pass 1: look for letters that are in the right spot
    for i, _ := range entry {
        if entry[i] == word[i] {
            word[i] = '_'
            output[i] = colorify(string(entry[i]), Black, Green)
        }
    
    }
    
    // Pass 2: check remaining letters
    for i, _ := range entry {
        // Skip if handled in previous pass
        if word[i] == '_' {
            continue
        }
        
        // Mark letters in the wrong spot yellow. Replace letters not found in the word with an underscore.
        if j := strings.IndexRune(string(word), entry[i]); j >= 0 {
            word[j] = '?'
            output[i] = colorify(string(entry[i]), Black, Yellow)
        } else {
            output[i] = "_"
        }
    }
    
    return strings.Join(output, " "), nil
}

func init() {
    // Get path to words file from environment variable or use "words"
    path := os.Getenv(wordsPathVar)
    if path == "" {
        path = "words"
    }
    file, err := os.Open(path)
    die(err)
    defer file.Close()
    
    // Load words into memory
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        w := scanner.Text()
        if re.MatchString(w) {
            words = append(words, w)
        }
    }
    
    if words == nil {
        die(fmt.Errorf("No words loaded!"))
    }
    
    // Shuffle words
    rand.Shuffle(len(words), func(i, j int) {
        words[i], words[j] = words[j], words[i]
    })
    
    bot.RegisterPassiveCommand("gowordle", gowordle)
}
