package truncate

import "testing"

func TestHTML(t *testing.T) {
	cases := []struct {
		in     string
		limit  int
		suffix string
		want   string
	}{
		{
			"",
			0,
			"",
			"",
		},
		{
			"",
			5,
			"",
			"",
		},
		{
			"123",
			0,
			"",
			"",
		},
		{
			"123",
			2,
			"",
			"12",
		},
		{
			"123",
			3,
			"",
			"123",
		},
		{
			"1234",
			3,
			"",
			"123",
		},
		{
			"<b>123</b>",
			5,
			"",
			"<b>123</b>",
		},
		{
			"<b>12345</b>",
			5,
			"",
			"<b>12345</b>",
		},
		{
			"<b>1234567</b>",
			5,
			"",
			"<b>12345</b>",
		},
		{
			"<b>Monty Python</b>",
			5,
			"",
			"<b>Monty</b>",
		},
		{
			"<img />",
			5,
			"",
			"<img />",
		},
		{
			"<img>",
			5,
			"",
			"<img>",
		},
		{
			"<h1><u>test<img blah blah>ing 1 2 3</u></h1>",
			5,
			"",
			"<h1><u>test<img blah blah>i</u></h1>",
		},
		{
			"123<h1><u> 456 <img blah blah> 789 012</u></h1>",
			7,
			"",
			"123<h1><u> 456 <img blah blah> 7</u></h1>",
		},
		{
			"<h1><u>😄u n i 😄 c😄o😄d😄e</u></h1>",
			5,
			"",
			"<h1><u>😄u n i 😄</u></h1>",
		},
		{
			"<h1><u>1234567</u></h1>",
			5,
			"...",
			"<h1><u>12345...</u></h1>",
		},
		{
			"<h1><u>1234 &copy; 1234</u></h1>",
			5,
			"",
			"<h1><u>1234 &copy;</u></h1>",
		},
		{
			"<h1><u>&copy;</u></h1>",
			1,
			"",
			"<h1><u>&copy;</u></h1>",
		},
		{
			"<h1><u>1234 &copy; 1234</u></h1>",
			6,
			"",
			"<h1><u>1234 &copy; 1</u></h1>",
		},
		{
			"<H1>UPPERCASE TAGS<BR><BR></H1>",
			100,
			"",
			"<H1>UPPERCASE TAGS<BR><BR></H1>",
		},
		{
			"🤨",
			500,
			"",
			"🤨",
		},
		{
			"🤨🤨🤨🤨🤨🤨🤨🤨🤨🤨🤨🤨",
			10,
			"",
			"🤨🤨🤨🤨🤨🤨🤨🤨🤨🤨",
		},
	}

	for _, c := range cases {
		out, err := HTML([]byte(c.in), c.limit, c.suffix)
		got := string(out)
		if err != nil {
			t.Errorf(`unexpected error calling TruncateHtml(%q, 5, ""): %v; want: %q`, c.in, err, c.want)
		}
		if got != c.want {
			t.Errorf("TruncateHtml(%q, %d, %q): got %q; want %q", c.in, c.limit, c.suffix, got, c.want)
		}
	}
}
