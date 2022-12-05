package main

import (
	"embed"
	"fmt"
	"reflect"
	"sort"
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
        <li>If you guess correctly, but can't explain your answer, you lose a point.</li>
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
		<p>What classic text adventure game begins like this?</p>
		<p><pre style="background: seagreen; color: white; display: inline-block; padding: 0.5rem; white-space: pre-wrap">You wake up. The room is spinning very gently round your head.
Or at least it would be if you could see it which you can't.

It is pitch black.

> <b>inventory</b>
You have:
  a splitting headache
  no tea

> <b>turn on light</b>
Good start to the day. Pity it's going to be the worst one of your life.
The light is now on.

> <b>open curtains</b>
As you part your curtains you see that it's a bright morning,
the sun is shining, the birds are singing, the meadows are blooming,
and a large yellow bulldozer is advancing on your home.

</pre></p>
    {{ template "input" "the game title" }}
	<button type="submit">A classic</button>`,
		reveal: `<h2>Answer</h2><p><i>The Hitchiker's Guide to the Galaxy</i>, developed by Infocom. You can still play it online at the <a href="https://www.bbc.co.uk/programmes/articles/1g84m0sXpnNCv84GpN2PLZG/the-game-30th-anniversary-edition" target="_blank">BBC website</a></p>`,
		check: func(answer string) bool {
			return strings.Contains(strings.ToLower(answer), "hitchhiker's guide")
		},
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
		<p>This grid of words is actually 4 groups of 4 related words. Find one of the groups.</p>
		<table style="text-align: left" border="1" cellpadding="5">
<tr><th><label><input type="checkbox" name="answer" value="Speck"> Speck</label></th> <th><label><input type="checkbox" name="answer" value="Bear off"> Bear off</label></th> <th><label><input type="checkbox" name="answer" value="Lechon"> Lechon</label></th> <th><label><input type="checkbox" name="answer" value="Trace"> Trace</label></th></tr>
<tr><th><label><input type="checkbox" name="answer" value="Autumn"> Autumn</label></th> <th><label><input type="checkbox" name="answer" value="Crumb"> Crumb</label></th> <th><label><input type="checkbox" name="answer" value="Iota"> Iota</label></th> <th><label><input type="checkbox" name="answer" value="Pip"> Pip</label></th></tr>
<tr><th><label><input type="checkbox" name="answer" value="Debris"> Debris</label></th> <th><label><input type="checkbox" name="answer" value="Prosciutto"> Prosciutto</label></th> <th><label><input type="checkbox" name="answer" value="Coup"> Coup</label></th> <th><label><input type="checkbox" name="answer" value="Gammon"> Gammon</label></th></tr>
<tr><th><label><input type="checkbox" name="answer" value="Smidgen"> Smidgen</label></th> <th><label><input type="checkbox" name="answer" value="Anchor"> Anchor</label></th> <th><label><input type="checkbox" name="answer" value="Morsel"> Morsel</label></th> <th><label><input type="checkbox" name="answer" value="Bacon"> Bacon</label></th></tr>
		</table>
		<p><button type="submit">I know this</button></p>`,
		check: func(answer string) bool {
			answers := strings.Split(answer, ",")
			sort.Strings(answers)
			if reflect.DeepEqual(answers, []string{"Autumn", "Coup", "Crumb", "Debris"}) {
				return true
			}
			if reflect.DeepEqual(answers, []string{"Bacon", "Lechon", "Prosciutto", "Speck"}) {
				return true
			}
			if reflect.DeepEqual(answers, []string{"Iota", "Morsel", "Smidgen", "Trace"}) {
				return true
			}
			if reflect.DeepEqual(answers, []string{"Anchor", "Bear off", "Gammon", "Pip"}) {
				return true
			}
			return false
		},
		reveal: `<h2>Answers</h2>
<p>Pork products - Bacon, Lechon, Prosciutto, Speck</p>
<p>Small amount - Iota, Morsel, Smidgen, Trace</p>
<p>Backgammon terms - Anchor, Bear off, Pip, Gammon</p>
<p>Ends with a silent letter - Autumn, Coup, Crumb, Debris</p>`,
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
		<button type="submit">Yes maestro</button>
		<p>It took me ages to write that sheet music, but you could also just listen to it as well I guess<br>
		<audio controls>
			<source src="/static/images/beethoven5.ogg" type="audio/ogg">
		  Your browser does not support the audio element. I guess you'll have to learn sheet music.
		  </audio>
		</p>`,
		check: func(answer string) bool {
			answer = strings.ToLower(answer)
			return strings.Contains(answer, "beethoven") && (strings.Contains(answer, "5") || strings.Contains(answer, "fifth"))
		},
		reveal: `<h2>Answer</h2><p>Beethoven's Fifth Symphony.</p>`,
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
