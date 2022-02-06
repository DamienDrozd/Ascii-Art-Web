package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Profile struct {
	// Download string
	Background string
	Output     string
	Color      string
	Text       string
}

var Buffer = "white"

func main() {
	os.Remove("templates/data/ascii.txt")                                                                // We delete the downloaded file so we don't have the one from previous launch
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates/static")))) // The server will analyse the static folder to seach thes called files in the html
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("templates/img"))))
	http.Handle("/font/", http.StripPrefix("/font/", http.FileServer(http.Dir("templates/font"))))
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("templates/data")))) // The server will analyse the data folder to seach thes called files in the html
	http.HandleFunc("/", serveur)                                                                  // Starts the server from the server function on port localhost:8080/ascii-art

	err := http.ListenAndServe(":8080", nil) // Starts the server on port 8080
	if err != nil {
		log.Fatal(err)
	}

}

func httpstatuscode() {
	resp, err := http.Get("http://localhost:8080/") //Generate the status code of the response which is returned
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 { //if the status code is between 200 and 300 the response is succesful
		fmt.Println("HTTP Status is in the 2xx range")
	} else {
		fmt.Println("Argh! Broken")
	}
}

var run = true

func serveur(w http.ResponseWriter, r *http.Request) {
	Text := r.FormValue("Text")     // Get the user's text input to convert to ascii
	Font := r.FormValue("Fontlist") // Get the user input to chose the font
	Color := r.FormValue("Color")   // Get the user input to chose the colour

	if r.FormValue("BackGroundColor") != "" {
		Buffer = r.FormValue("BackGroundColor") // Get the user input to chose the background colour
	}

	if Text != "" {
		writefile(ascii(Text, Font)) // If the input is invalid, we do not create a downloadable file
	}

	p := Profile{
		Background: "/static/" + Buffer + ".css",                   // The Background takes the value of the css file to import it in the html file
		Output:     ascii(Text, Font),                              // The output is created with the text input and the ascii function
		Color:      "/static/textcolor/textcolor" + Color + ".css", // The css file takes the value that we want to display via a css file
		Text:       Text,                                           // The text has the value of the user input
	}

	// Creation of a new template instance
	t := template.New("Label de ma template")

	// Declaration of files to parse
	t = template.Must(t.ParseFiles("./templates/index.html"))

	if run == true {
		run = false
		httpstatuscode()

	}

	// Execution of the fusion and injection in the exit flux
	// The variable p wille be presented by "." in the layout
	err := t.ExecuteTemplate(w, "layout", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err) // If the executetemplate function cannot run, displays an error message
	}

}

func writefile(input string) {
	os.Remove("templates/data/ascii.txt")
	outputfile, _ := os.OpenFile("templates/data/ascii.txt", os.O_CREATE|os.O_WRONLY, 0600) // Creates the ascii.txt file with writing and reading parameters
	outputfile.WriteString(input)                                                           // Puts the input in the file
	outputfile.Close()
}

func ascii(input string, Font string) string {
	// Get user input
	wordtab := newline(input) // Our array containing all of our words
	var output string
	output += "\n"
	for _, word := range wordtab { // Ranges our array to display every word one by one
		output = findLine(word, output, Font)
		output += "\n"
	}
	return output
}

func findLine(mot, output, Font string) string {
	// Writing

	for i := 0; i < 8; i++ { // Count lines
		for _, lettre := range mot { // Divides words into letters
			lettre -= 32                           // Converts value from ASCII to local library, get rid of useless ASCII commands
			lettre *= 9                            // Finds the right letter in our library by multiplying by the amount of line per character
			ligne := int(lettre) + 2 + i           // Converts our rune into it's decimal value, "i" is the line value
			output = readfile(ligne, output, Font) // Calls our function to display the letter's line
		}

		if i < 7 {
			output += "\n"
		}

	}
	return output
}

// Reads our library

func readfile(ligne int, output string, Font string) string { // Checks for a librairy value
	var filename string
	filename = `templates/data/standard.txt` // Base value is standard.txt
	if Font == "standard" {
		filename = `templates/data/standard.txt`
	}
	if Font == "shadow" {
		filename = `templates/data/shadow.txt`
	}
	if Font == "thinkertoy" {
		filename = `templates/data/thinkertoy.txt`
	}
	file, _ := os.Open(filename)          // Opens our file
	fileScanner := bufio.NewScanner(file) // Scan our file
	i := 0
	for fileScanner.Scan() { // Reads lines one by one
		i++
		if i == ligne { // When it meets our "ligne" value, calls our "for"
			for _, r := range fileScanner.Text() { // Divides our line into characters
				if r != '\n' { // If our character isn't a "\n", prints it
					output += string(r)
				}
			}

			break // Breaks our loop
		}
	}
	file.Close() // Closes our file

	return output
}

func newline(input string) []string { // Separates words from '\n' and puts them in an array
	tablenght := 0
	var wordtab = make([]string, 1) // Creates tab
	index := 0
	for x := range input { // Ranges our input and look for a '\n'
		if (input[x] == 92 && input[x+1] == 110) || (input[x] == 10) { // 92 is '\', 110 is 'n' and 10 is the '\n' command
			wordtab[tablenght] = (input[index:x]) // x equals the value at which we encounter a '\n', stores the word prior to that in the array
			wordtab = append(wordtab, "0")        // Adds a slot in our array
			if input[x] == 10 {
				index = x + 1
			} else {
				index = x + 2
			}
			tablenght++
		}
	}
	wordtab[tablenght] = input[index:len(input)] // Store last word
	return wordtab
}
