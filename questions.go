package main

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

type question struct {
	q      string
	reveal string
	check  func(answer string) bool
}

var questions = []question{
	{
		q: `<header>
	<h2>Welcome to Pmoney's impossible quiz</h2>
    <p><b>Rules</b>:</p>
    <ul>
        <li>1 point for the first correct answer</li>
        <li>Google all you want, guess as many times as you want</li>
        <li>Everyone can see all your guesses though</li>
        <li>If you answer correctly, you have to explain why. If you can't explain why, you lose a point.</li>
    </ul>
    <p>Tip: If something goes wrong, refreshing the page should fix it. You won't lose progress on the quiz, since the state is server-side.</p>
    <p></p>
	<p>Waiting for pmoney to start this once in a lifetime event...</p>
</header>`,
		check: func(answer string) bool {
			return true
		},
	},
	{
		q: `<h2>Question {{ .Question }} - for practice</h2>
	<p>What is the 9 letter anagram?</p>
    <p><img src="/static/images/target.png" alt="RDE VUT ANE"></p>
    {{ template "input" "the 9 letter word" }}
	<button type="submit">I solved it</button>`,
		check: func(answer string) bool {
			return strings.EqualFold(answer, "adventure")
		},
		reveal: `<h2>Answer</h2><p><b>ADVENTURE</b></p>`,
	},
	{
		q: `<h2>Question {{ .Question }}</h2>
    <p>What is this strange UI?</p>
    <p><img src="/static/images/mysteriousdevice.jpeg" alt="A mysterious ancient technology"></p>
    {{ template "input" "mysterious device" }}
	<button type="submit">Easy</button>`,
		check: func(answer string) bool {
			return strings.Contains(strings.ToLower(answer), "zune")
		},
		reveal: `<h2>Answer</h2><p>The <b><a href="https://en.wikipedia.org/wiki/Zune">Microsoft Zune</a></b>!</p>`,
	},
	{
		q: `<h2>Question {{ .Question }}</h2>
		<p>What has a head and a tail but no body?</p>
		{{ template "input" "answer" }}
		   <button type="submit">I eat riddles for breakfast</button>`,
		check: func(answer string) bool {
			return strings.Contains(strings.ToLower(answer), "coin")
		},
		reveal: `<h2>Answer</h2><p>A coin.</p>`,
	},
	{
		q: `<h2>Question {{ .Question }}</h2>
		<p>White and black take turns placing one stone on the board. White can win by surrounding black's stones.</p>
		<p>What is the coordinate of the correct next move to ensure victory?</p>
		<div class="jgoboard" data-jgostyle="JGO.BOARD.mediumBW">
			.o..x....
			xxxxxo.o.
			o...oo...
			.o.o.....
			.........
		</div>
		<script>
			JGO.auto.init(document, JGO);
		</script>
		<p>
			{{ template "input" "coordinate" }}
			<button type="submit">I am supremely confident</button>
		</p>
		<p>If you need a hint, check out the <a href="https://en.wikipedia.org/wiki/Rules_of_Go#Concise_statement" target="_blank">Concise rules of Go</a>.</p>`,
		check: func(answer string) bool {
			return strings.EqualFold(answer, "c5")
		},
		reveal: `<h2>Answer</h2>
		<div class="jgoboard" data-jgostyle="JGO.BOARD.mediumBW">
			.oo.x....
			xxxxxo.o.
			o...oo...
			.o.o.....
			.........
		</div>
		<script>
			JGO.auto.init(document.getElementById('answer-{{ .Data }}'), JGO);
		</script>
		<p>The answer is <b>c5</b>. After <b>f5 g5</b> black cannot escape (if black tries <b>d5</b> instead, that leads to <b>a5</b>, <b>c3</b> leads to <b>d3</b>).</p>
		<p>Note that playing <b>d5</b> first is incorrect as then black can play <b>c5</b>, surrounding d5 and preventing a5 or d5 in the future.</p>`,
	},
	{
		q: `<h2>Question {{ .Question }}</h2>
		<p>In which country is it traditional to eat KFC for Christmas dinner?</p>
		{{ template "input" "a country" }}
		<button type="submit">My travels have not been in vain</button>`,
		check: func(answer string) bool {
			return strings.EqualFold(answer, "japan")
		},
		reveal: `<h2>Answer</h2><p>Japan.</p>`,
	},
	{
		q: `<h2>Question {{ .Question }}</h2>
		<p>What is the missing ingredient in this recipe for apple pie?</p>
		<p><b>Ingredients (pie crust)</b><br>
			<ul>
				<li>2&frac12; cups flour</li>
				<li>&frac12; tbsp sugar</li>
				<li>&frac12; tsp salt</li>
				<li>200g unsalted butter</li>
			</ul>
		</p>
		<p><b>Ingredients (filling)</b><br>
		<ul>
			<li>1kg Granny Smith apples</li>
			<li>8 tbsp unsalted butter</li>
			<li>3 tbsp flour</li>
			<li>&frac14; cup water</li>
			<li>1 cup sugar</li>
			<li>1 egg + 1 tbsp water, for egg wash</li>
		</ul></p>
		{{ template "input" "the secret ingredient" }}
		<button type="submit">Yes chef</button>`,
		check: func(answer string) bool {
			return strings.Contains(strings.ToLower(answer), "cinnamon")
		},
		reveal: `<h2>Answer</h2><p>The missing ingredient is <b>cinnamon</b>. Essential for that classic apple pie flavour.</p>`,
	},
	{
		q: `<h2>Question {{ .Question }}</h2>
		<p>This is bars 5-8 of which piece?</p>
		<p><img src="/static/images/sheetmusic.png" alt="sheet music" width="500"></p>
		{{ template "input" "the piece" }}
		<button type="submit">Yes maestro</button>`,
		check: func(answer string) bool {
			answer = strings.ToLower(answer)
			return strings.Contains(answer, "beethoven") && (strings.Contains(answer, "5") || strings.Contains(answer, "fifth"))
		},
		reveal: `<h2>Answer</h2>
		<p>Beethoven's Fifth Symphony.</p>
		<p>
		<audio controls>
			<source src="/static/images/beethoven5.ogg" type="audio/ogg">
		  Your browser does not support the audio element.
		  </audio>
		</p>`,
	},
	{
		q: `<h2>You have reached the end of the pmoney quiz</h2>
		<p>Thanks for playing!</p>`,
	},
}

var (
	//go:embed templates/*.html
	templateFS embed.FS

	templates *template.Template
)

func init() {
	buf := new(strings.Builder)
	buf.WriteString(`{{ define "question" }}` + "\n")
	for num, conf := range questions {
		fmt.Fprintf(buf, `{{ if eq .Question %d }}%s{{ end }}`+"\n", num, conf.q)
	}
	buf.WriteString(`{{ end }}` + "\n")
	buf.WriteString(`{{ define "reveal" }}` + "\n")
	for num, conf := range questions {
		fmt.Fprintf(buf, `{{ if eq .Data "%d" }}%s{{ end }}`+"\n", num, conf.reveal)
	}
	buf.WriteString(`{{ end }}` + "\n")
	templates = template.Must(template.New("").Parse(buf.String()))
	templates = template.Must(templates.ParseFS(templateFS, "templates/*.html"))
}
