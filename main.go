package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {

	// Read from a file passed as command argument (eg. ./rubaeparla file-to-convert.txt)
	b, err := ioutil.ReadFile(os.Args[1]) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	fileToread := string(b) // convert content to a 'string'

	// this is a sentence
	split := strings.SplitAfter(fileToread, ".")

	// Init destination audio empty file (see concatenateAudio for usage)
	destinationFile, err := os.Create("output.mp3")
	if err != nil {
		log.Fatal(err)
	}
	destinationFile.Close()

	// call convert function
	convert(split)
}

func convert(split []string) {

	// for every line (that ends with ".") convert to audio - Every sentence is an MP3
	for i := range split {

		// Instantiates a client.
		ctx := context.Background()

		client, err := texttospeech.NewClient(ctx)
		if err != nil {
			log.Fatal(err)
		}

		// Perform the text-to-speech request on the text input with the selected
		// voice parameters and audio file type.
		req := texttospeechpb.SynthesizeSpeechRequest{
			// Set the text input to be synthesized.
			Input: &texttospeechpb.SynthesisInput{
				InputSource: &texttospeechpb.SynthesisInput_Text{Text: split[i]},
			},
			// Build the voice request, select the language code ("en-US") and the SSML
			// voice gender ("neutral").
			Voice: &texttospeechpb.VoiceSelectionParams{
				LanguageCode: "en-GB",
				Name:         "en-GB-Wavenet-D",
				SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
			},
			// Select the type of audio file you want returned.
			// those settings are quite OK
			AudioConfig: &texttospeechpb.AudioConfig{
				AudioEncoding: texttospeechpb.AudioEncoding_MP3,
				SpeakingRate:  0.9,
				Pitch:         -2,
			},
		}

		resp, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			log.Fatal(err)
		}
		// Every sentence is an MP3 file
		//filename := strconv.Itoa(i) + ".mp3"
		//err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Printf("Splitting audio content to file: %v\n", filename)

		time.Sleep(500 * time.Millisecond) // avoid google throttling
		concatenateAudio(resp.AudioContent)

	}
}

func concatenateAudio(AudioContent []byte) {

	f, err := os.OpenFile("output.mp3", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.Write(AudioContent); err != nil {
		panic(err)
	}

}
