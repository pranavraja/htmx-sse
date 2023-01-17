package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"strings"
)

type question struct {
	q           string
	reveal      string
	onePoint    func(string) bool
	bonusPoints func(string) bool
}

// e.g. @45.2547719,19.8451607,3a,61.8y,47.6h,96.44t
func pano(name, pos, id string) string {
	var (
		lat      float64
		lng      float64
		zoom     int
		altitude float64
		heading  float64
		pitch    float64
	)
	if _, err := fmt.Sscanf(pos, `@%f,%f,%da,%fy,%fh,%ft`, &lat, &lng, &zoom, &altitude, &heading, &pitch); err != nil {
		log.Fatalf("couldn't scan %s: %s", pos, err)
	}
	return fmt.Sprintf(`<h2>Question {{ dec .Question }} - "%s"</h2>
<p>Click and drag the image, and use arrow keys to move and turn. If the map goes black, try refreshing 😬</p>
<div class="pano" id="map-{{ .Question }}"></div>
<script>
window.panorama = new google.maps.StreetViewPanorama(
    document.getElementById("map-{{ .Question }}"),
    {
	  position: { lat: %f, lng: %f },
	  pov: { heading: %f, pitch: %f },
	  pano: %q,
	  zoom: %d,
      addressControl: false,
      linksControl: false,
      clickToGo: false,
      showRoadLabels: false
    }
  );
</script>
<p>
{{ template "input" . }}
<button type="submit">I know this</button>
</p>`, name, lat, lng, heading, pitch-90, id, zoom-1)
}

func containsAny(terms ...string) func(string) bool {
	return func(answer string) bool {
		answer = strings.ToLower(answer)
		for _, t := range terms {
			if strings.Contains(answer, t) {
				return true
			}
		}
		return false
	}
}

var questions = []question{
	{
		q: `<header>
	<h2><s>city</s> <s>movie</s> city guesser</h2>
    <p><b>Rules</b>:</p>
    <ul>
        <li>You'll be shown an interactive panorama selected from Google Maps</li>
        <li><strong>1 point</strong> for guessing the country, <strong>3 points</strong> for guessing the city. Duplicate guesses will be ignored.</li>
        <li>At the end of the game, the person with the highest <code>points per guess</code> is the winner</li>
    </ul>
</header>`,
		bonusPoints: func(answer string) bool {
			return true
		},
	},
	{
		q: `<div><canvas id="canvas" width="600" height="600"></canvas></div>
	<p>Waiting for pmoney to start this once in a lifetime event...</p>
    <p><b>Rules reminder</b>:</p>
    <ul>
        <li>You'll be shown an interactive panorama selected from Google Maps</li>
        <li><strong>1 point</strong> for guessing the country, <strong>3 points</strong> for guessing the city. Duplicate guesses will be ignored.</li>
        <li>At the end of the game, the person with the highest <code>points per guess</code> is the winner</li>
    </ul>
                <script defer="defer">
    const animation = new rive.Rive({
        src: '/static/images/earth.riv',
        canvas: document.getElementById("canvas"),
        autoplay: true
    });
                </script>`,
		bonusPoints: func(answer string) bool {
			return true
		},
	},
	{
		q: `<script defer="defer">
	animation.cleanup();
		</script>` + pano("Traffic Light Tree", "@51.5067585,-0.0107938,3a,26.7y,70.37h,95.09t", ""),
		onePoint:    containsAny("united kingdom", "uk", "england"),
		bonusPoints: containsAny("london"),
		reveal:      `<p><b>London, UK</b></p>`,
	},
	{
		q:           pano("Camels Chilling", "@29.9772294,31.1308562,2a,63.6y,341.24h,84.06t", ""),
		bonusPoints: containsAny("giza"),
		onePoint:    containsAny("egypt"),
		reveal:      `<p><b>Giza, Egypt</b></p>`,
	},
	{
		q:           pano("Foreshore Freeway Bridge", "@-33.9152697,18.4221139,3a,73.7y,320.69h,102.33t", ""),
		onePoint:    containsAny("south africa"),
		bonusPoints: containsAny("cape town"),
		reveal:      `<p><b>Cape Town, South Africa</b></p>`,
	},
	{
		q:           pano("Umbrella Field", "@48.8817289,2.3823744,3a,64.3y,353.11h,87.01t", ""),
		onePoint:    containsAny("france"),
		bonusPoints: containsAny("paris"),
		reveal:      `<p><b>Paris, France</b></p>`,
	},
	{
		q:           pano("Penny Farthing", "@-32.010644,115.7522975,3a,49.2y,92.43h,76.43t", "EgQG9n-JKdPXrM5bFUl6ow"),
		onePoint:    containsAny("australia"),
		bonusPoints: containsAny("perth"),
		reveal:      `<p><b>Perth, Australia</b></p>`,
	},
	{
		q:           pano("Scott's Hut", "@-77.6360568,166.4179189,2a,90y,242.17h,89.15t", ""),
		onePoint:    containsAny("antarctica", "new zealand"),
		bonusPoints: containsAny("cape evans"),
		reveal:      `<p><b>Cape Evans, Antarctica/New Zealand</b></p>`,
	},
	{
		q:           pano("Swole Bus", "@50.0379256,14.4976932,3a,24.8y,51.89h,90.22t", ""),
		bonusPoints: containsAny("prague"),
		onePoint:    containsAny("czechia", "czech republic"),
		reveal:      `<p><b>Prague, Czechia</b></p>`,
	},
	{
		q:           pano("Chicken Church", "@27.7932314,-82.7902872,3a,20.5y,325.09h,99.02t", ""),
		onePoint:    containsAny("united states", "usa"),
		bonusPoints: containsAny("madeira", "florida"),
		reveal:      `<p><b>Madeira Beach, Florida, USA</b></p>`,
	},
	{
		q:           pano("Cave Restaurant", "@-4.3129601,39.576904,3a,75y,225.68h,86.34t", ""),
		onePoint:    containsAny("kenya"),
		bonusPoints: containsAny("mombasa"),
		reveal:      `<p><b>Mombasa, Kenya</b></p>`,
	},
	{
		q:           pano("Waterfall Restaurant", "@13.9940387,121.3450291,2a,86.4y,57.9h,84.04t", ""),
		onePoint:    containsAny("philippines"),
		bonusPoints: containsAny("san pablo"),
		reveal:      `<p><b>San Pablo City, Philippines</b></p>`,
	},
	{
		q:           pano("Carrying Bread", "@31.7766664,35.2278165,3a,55.5y,352.25h,80.62t", ""),
		onePoint:    containsAny("israel"),
		bonusPoints: containsAny("jerusalem"),
		reveal:      `<p><b>Jerusalem, Israel</b></p>`,
	},
	{
		q:           pano("New Employee", "@45.5028684,9.1850453,4a,22y,303.67h,111.59t", ""),
		onePoint:    containsAny("italy"),
		bonusPoints: containsAny("milan"),
		reveal:      `<p><b>Milan, Italy</b></p>`,
	},
	{
		q:           pano("Floating Baby", "@1.2805073,103.8623109,2a,47.1y,250.29h,94.84t", ""),
		onePoint:    containsAny("singapore"),
		bonusPoints: containsAny("singapore"),
		reveal:      `<p><b>Singapore</b></p>`,
	},
	{
		q:           pano("Pigeon Attack", "@-16.495704,-68.1337214,2a,75y,296.85h,62.81t", ""),
		onePoint:    containsAny("bolivia"),
		bonusPoints: containsAny("la paz"),
		reveal:      `<p><b>La Paz, Bolivia</b></p>`,
	},
	{
		q:           pano("Shark Attack", "@51.758722,-1.213069,3a,32.7y,330.41h,103.41t", ""),
		onePoint:    containsAny("england", "uk", "united kingdom"),
		bonusPoints: containsAny("oxford"),
		reveal:      `<p><b>Oxford, England</b></p>`,
	},
	{
		q:           pano("Great Parking", "@23.0361074,120.231361,3a,85.9y,178.97h,106.01t", ""),
		onePoint:    containsAny("taiwan"),
		bonusPoints: containsAny("tainan"),
		reveal:      `<p><b>Tainan City, Taiwan</b></p>`,
	},
	{
		q:           pano("Ice Lounge", "@-36.8421189,174.7648173,2a,90y,147.88h,87.89t", ""),
		onePoint:    containsAny("nz", "new zealand", "aotearoa"),
		bonusPoints: containsAny("auckland"),
		reveal:      `<p><b>Auckland, New Zealand</b></p>`,
	},
	{
		q: `<h2>You have reached the end of the pmoney quiz</h2>
		<p>Thanks for playing!</p>
		<p>If you enjoyed this and want more, see <a href="https://neal.fun/wonders-of-street-view/">Wonders of Street View</a>.</p>`,
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
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"dec": func(i int64) int64 {
			return i - 1
		},
	}).Parse(buf.String()))
	templates = template.Must(templates.ParseFS(templateFS, "templates/*.html"))
}
