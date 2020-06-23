package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	BulkLocation  = "https://api.scryfall.com/bulk-data/"
	BulkLog       = "ScryGoBulk.info"
	SaveLocation  = "D:/Scryfall/"
)

func main() {
	// Download all of the bulk data available from the Scryfall API locally to ~Make Life Easier (tm)~
	fmt.Println("Downloading Bulk Data elements...")
	err := DownloadBulkData()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Bulk Data Downloaded.")
	}

	err = test()

	time.Sleep(time.Second * 2)	

	err = CatRulings(SaveLocation + "Rulings.json")
}

func test() error {
	return nil
}

func CatRulings(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path doesn't exist, we continue onward
			return errors.New("The file you attempted to cat [" + path + "] does not exist.")
		} else {
			// File exists, but there is an error
			return err
		}
	}
	// File exists, check the most recent update
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the contents
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var rulings []Ruling
	err = json.Unmarshal(contents, &rulings)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, ruling := range rulings {
		fmt.Println(ruling)
	}

	return nil
}

func DownloadBulkData() error {

	// Check for the ScryGo.info file and look for the last update
	_, err := os.Stat(SaveLocation + BulkLog)
	if err != nil {
		if os.IsNotExist(err) {
			// Path doesn't exist, we continue onward
		} else {
			// File exists, but there is an error
			return err
		}
	} else {
		// File exists, check the most recent update
		file, err := os.Open(SaveLocation + BulkLog)
		if err != nil {
			return err
		}
		defer file.Close()

		// Get and check the date to see if it is within the last 24 hours
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		date, _ := time.Parse(time.UnixDate, string(contents))
		// Prints the found date in Unix date format
		// fmt.Println("Found containing: " + string(contents) + " translating to " + date.Format(time.UnixDate))
		// fmt.Println("Now: " + time.Now().Format(time.UnixDate))
		// fmt.Println(fmt.Sprintf("%0.2f", time.Now().Sub(date).Hours()))
		if time.Now().Sub(date).Hours() < 24 {
			return errors.New("The bulk data has been updated within the last 24 hours. You don't need to update it yet.")
		}
	}

	// Quick sanity check
	if !PromptYN("The Bulk Data files are likely over a gigabyte of data that will be downloaded. Do you want to update these files?") {
		return errors.New("Download aborted by user.")
	}

	// Get the bulk data items from Scryfall
	resp, err := http.Get(fmt.Sprintf("%s", BulkLocation))
	if err != nil {
		return err
	}

	// Parse the JSON blob
	bytearray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	blob, err := ParseJSONList(bytearray)
	if err != nil {
		return err
	}

	// Download each of the bulk data objects
	for _, data_element := range blob.Data {
		path := SaveLocation + data_element.Name + ".json"
		err = DownloadFile(path, data_element.Download_URI)
		if err != nil {
			return err
		}
	}

	// Open the log file
	file, err := os.Create(SaveLocation + BulkLog)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add the current date to the log file for future checking
	// now := time.Now()
	// year_string := fmt.Sprintf("%04d", now.Year())
	// month_string := fmt.Sprintf("%02d", int(now.Month()))
	// day_string := fmt.Sprintf("%02d", now.Day())
	// timezone_string, _ := now.Zone()
	// _, err = file.WriteString(string(year_string + month_string + day_string + " " + timezone_string))

	_, err = file.WriteString(time.Now().Format(time.UnixDate))

	return nil
}

func PromptYN(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " [y/n]: ")
	text, _ := reader.ReadString('\n')
	switch text[0] {
	case 'y':
		return true
	case 'n':
		return false
	}
	return PromptYN(prompt)
}

type JSONElement struct {
	// Because I hate myself and haven't figured out a more elegant way to make this work the way I want it to
	ObjectType            string              `json:"object"`
	Status                int64               `json:"status"`
	Code                  string              `json:"code"`
	Details               string              `json:"details"`
	Type                  string              `json:"type"`
	Warnings              []string            `json:"warnings"`
	Data                  []JSONElement       `json:"data"`
	Has_More              bool                `json:"has_more"`
	Next_Page             string              `json:"next_page"`
	Total_Cards           int                 `json:"total_cards"`
	ID                    string              `json:"id"`
	MTGO_Code             string              `json:"mtgo_code"`
	Arena_Code            string              `json:"arena_code"`
	TCGPlayer_ID          string              `json:"tcgplayer_id"`
	Name                  string              `json:"name"`
	Set_Type              string              `json:"set_type"`
	Released_At           string              `json:"released_at"`
	Block_Code            string              `json:"block_code"`
	Block                 string              `json:"block"`
	Parent_Set_Code       string              `json:"parent_set_code"`
	Card_Count            int                 `json:"card_count"`
	Digital               bool                `json:"digital"`
	Foil_Only             bool                `json:"foil_only"`
	Nonfoil_Only          bool                `json:"nonfoil_only"`
	Scryfall_URI          string              `json:"scryfall_uri"`
	URI                   string              `json:"uri"`
	Icon_SVG_URI          string              `json:"icon_svg_uri"`
	Search_URI            string              `json:"search_uri"`
	Oracle_ID             string              `json:"oracle_id"`
	Source                string              `json:"source"`
	Published_At          string              `json:"published_at"`
	Comment               string              `json:"comment"`
	Symbol                string              `json:"symbol"`
	SVG_URI               string              `json:"svg_uri"`
	Loose_Variant         string              `json:"loose_variant"`
	Engligh               string              `json:"english"`
	Transposable          bool                `json:"transposable"`
	Represents_Mana       bool                `json:"represents_mana"`
	CMC                   float64             `json:"cmc"`
	Appears_In_Mana_Costs bool                `json:"appears_in_mana_costs"`
	Funny                 bool                `json:"funny"`
	Colors                []Color             `json:"colors"`
	Gatherer_Alternates   []string            `json:"gatherer_alternates"`
	Total_Values          int                 `json:"total_values"`
	Updated_At            string              `json:"updated_at"`
	Description           string              `json:"description"`
	Compressed_size       int                 `json:"compressed_size"`
	Download_URI          string              `json:"download_uri"`
	Content_Type          string              `json:"content_type"`
	Content_Encoding      string              `json:"content_encoding"`
	Arena_ID              int                 `json:"arena_id"`
	Language              string              `json:"lang"`
	MTGO_ID               int                 `json:"mtgo_id"`
	MTGO_Foil_ID          int                 `json:"mtgo_foil_id"`
	Multiverse_IDs        []int               `json:"multiverse_ids"`
	Prints_Search_URI     string              `json:"prints_search_uri"`
	Rulings_URI           string              `json:"rulings_uri"`
	All_Parts             []RelatedCardObject `json:"all_parts"`
	Card_Faces            []CardFaceObject    `json:"card_faces"`
	Keywords              []string            `json:"keywords"`
	Color_Identity        []Color             `json:"color_identity"`
	Color_Indicator       []Color             `json:"color_indicator"`
	EDHRec_Rank           int                 `json:"edhrec_rank"`
	Foil                  bool                `json:"foil"`
	Hand_Modifier         string              `json:"hand_modifier"`
	Layout                string              `json:"layout"`
	Legalities            LegalitiesObject    `json:"legalities"`
	Life_Modifier         string              `json:"life_modifier"`
	Loyalty               string              `json:"loyalty"`
	Mana_Cost             string              `json:"mana_cost"`
	Nonfoil               bool                `json:"nonfoil"`
	Oracle_Text           string              `json:"oracle_text"`
	Oversized             bool                `json:"oversized"`
	Power                 string              `json:"power"`
	Reserved              bool                `json:"reserved"`
	Toughness             string              `json:"toughness"`
	Type_Line             string              `json:"type_line"`
	Artist                string              `json:"artist"`
	Artist_IDs            []string            `json:"artist_ids"`
	Booster               bool                `json:"booster"`
	Border_Color          string              `json:"border_color"`
	Card_Back_ID          string              `json:"card_back_id"`
	Collector_Number      string              `json:"collector_number"`
	Content_Warning       bool                `json:"content_warning"`
	Flavor_Name           string              `json:"flavor_name"`
	Flavor_Text           string              `json:"flavor_text"`
	Frame_Effects         []string            `json:"frame_effects"`
	Frame                 string              `json:"frame"`
	Full_Art              bool                `json:"full_art"`
	Games                 []string            `json:"games"`
	Highres_Image         bool                `json:"highres_image"`
	Illustration_ID       string              `json:"illustartion_id"`
	Image_URIs            ImageURIsObject     `json:"image_uris"`
	Prices                PricesObject        `json:"prices"`
	Printed_Name          string              `json:"printed_name"`
	Printed_Text          string              `json:"printed_text"`
	Printed_Type_Line     string              `json:"printer_type_line"`
	Promo                 bool                `json:"promo"`
	Promo_Types           []string            `json:"promo_types"`
	Purchase_URIs         PurchaseURIsObject  `json:"purchase_uris"`
	Rarity                string              `json:"rarity"`
	Related_URIs          RelatedURIsObject   `json:"related_uris"`
	Reprint               bool                `json:"reprint"`
	Scryfall_Set_URI      string              `json:"scryfall_set_uri"`
	Set_Name              string              `json:"set_name"`
	Set_Search_URI        string              `json:"set_search_uri"`
	Set_URI               string              `json:"set_uri"`
	Set                   string              `json:"set"`
	Story_Spotlight       bool                `json:"story_spotlight"`
	Textless              bool                `json:"textless"`
	Variation             bool                `json:"variation"`
	Variation_Of          string              `json:"variation_of"`
	Watermark             string              `json:"watermark"`
	Preview_Previewed_At  string              `json:"preview.previewed_at"`
	Preview_Source_URI    string              `json:"preview.source_uri"`
	Preview_Source        string              `json:"preview.source"`
	Component             string              `json:"component"`
}

type ScryfallError struct {
	ObjectType string   `json:"object"`
	Status     int64    `json:"status"`
	Code       string   `json:"code"`
	Details    string   `json:"details"`
	Type       string   `json:"type"`
	Warnings   []string `json:"warnings"`
}

type List struct {
	ObjectType  string        `json:"object"`
	Data        []JSONElement `json:"data"`
	Has_More    bool          `json:"has_more"`
	Next_Page   string        `json:"next_page"`
	Total_Cards int           `json:"total_cards"`
	Warnings    []string      `json:"warnings"`
}

type Set struct {
	ObjectType      string `json:"object"`
	ID              string `json:"id"`
	Code            string `json:"code"`
	MTGO_Code       string `json:"mtgo_code"`
	Arena_Code      string `json:"arena_code"`
	TCGPlayer_ID    string `json:"tcgplayer_id"`
	Name            string `json:"name"`
	Set_Type        string `json:"set_type"`
	Released_At     string `json:"released_at"`
	Block_Code      string `json:"block_code"`
	Block           string `json:"block"`
	Parent_Set_Code string `json:"parent_set_code"`
	Card_Count      int    `json:"card_count"`
	Digital         bool   `json:"digital"`
	Foil_Only       bool   `json:"foil_only"`
	Nonfoil_Only    bool   `json:"nonfoil_only"`
	Scryfall_URI    string `json:"scryfall_uri"`
	URI             string `json:"uri"`
	Icon_SVG_URI    string `json:"icon_svg_uri"`
	Search_URI      string `json:"search_uri"`
}

type Ruling struct {
	ObjectType   string `json:"object"`
	Oracle_ID    string `json:"oracle_id"`
	Source       string `json:"source"`
	Published_At string `json:"published_at"`
	Comment      string `json:"comment"`
}

type CardSymbol struct {
	ObjectType            string   `json:"object"`
	Symbol                string   `json:"symbol"`
	SVG_URI               string   `json:"svg_uri"`
	Loose_Variant         string   `json:"loose_variant"`
	Engligh               string   `json:"english"`
	Transposable          bool     `json:"transposable"`
	Represents_Mana       bool     `json:"represents_mana"`
	CMC                   float64  `json:"cmc"`
	Appears_In_Mana_Costs bool     `json:"appears_in_mana_costs"`
	Funny                 bool     `json:"funny"`
	Colors                []Color  `json:"colors"`
	Gatherer_Alternates   []string `json:"gatherer_alternates"`
}

type Catalog struct {
	ObjectType   string   `json:"object"`
	URI          string   `json:"uri"`
	Total_Values int      `json:"total_values"`
	Data         []string `json:"data"`
}

type BulkData struct {
	ObjectType       string `json:"object"`
	ID               string `json:"id"`
	Type             string `json:"type"`
	Updated_At       string `json:"updated_at"`
	URI              string `json:"uri"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Compressed_size  int    `json:"compressed_size"`
	Download_URI     string `json:"download_uri"`
	Content_Type     string `json:"content_type"`
	Content_Encoding string `json:"content_encoding"`
}

type Color string

type Card struct {
	// Core Fields
	Arena_ID          int    `json:"arena_id"`
	ID                string `json:"id"`
	Language          string `json:"lang"`
	MTGO_ID           int    `json:"mtgo_id"`
	MTGO_Foil_ID      int    `json:"mtgo_foil_id"`
	Multiverse_IDs    []int  `json:"multiverse_ids"`
	TCGPlayer_ID      int    `json:"tcgplayer_id"`
	ObjectType        string `json:"object"`
	Oracle_ID         string `json:"oracle_id"`
	Prints_Search_URI string `json:"prints_search_uri"`
	Rulings_URI       string `json:"rulings_uri"`
	Scryfall_URI      string `json:"scryfall_uri"`
	URI               string `json:"uri"`

	// Gameplay Fields
	All_Parts       []RelatedCardObject `json:"all_parts"`
	Card_Faces      []CardFaceObject    `json:"card_faces"`
	CMC             float64             `json:"cmc"`
	Colors          []Color             `json:"colors"`
	Keywords        []string            `json:"keywords"`
	Color_Identity  []Color             `json:"color_identity"`
	Color_Indicator []Color             `json:"color_indicator"`
	EDHRec_Rank     int                 `json:"edhrec_rank"`
	Foil            bool                `json:"foil"`
	Hand_Modifier   string              `json:"hand_modifier"`
	Layout          string              `json:"layout"`
	Legalities      LegalitiesObject    `json:"legalities"`
	Life_Modifier   string              `json:"life_modifier"`
	Loyalty         string              `json:"loyalty"`
	Mana_Cost       string              `json:"mana_cost"`
	Name            string              `json:"name"`
	Nonfoil         bool                `json:"nonfoil"`
	Oracle_Text     string              `json:"oracle_text"`
	Oversized       bool                `json:"oversized"`
	Power           string              `json:"power"`
	Reserved        bool                `json:"reserved"`
	Toughness       string              `json:"toughness"`
	Type_Line       string              `json:"type_line"`

	// Print Fields
	Artist               string             `json:"artist"`
	Artist_IDs           []string           `json:"artist_ids"`
	Booster              bool               `json:"booster"`
	Border_Color         string             `json:"border_color"`
	Card_Back_ID         string             `json:"card_back_id"`
	Collector_Number     string             `json:"collector_number"`
	Content_Warning      bool               `json:"content_warning"`
	Digital              bool               `json:"digital"`
	Flavor_Name          string             `json:"flavor_name"`
	Flavor_Text          string             `json:"flavor_text"`
	Frame_Effects        []string           `json:"frame_effects"`
	Frame                string             `json:"frame"`
	Full_Art             bool               `json:"full_art"`
	Games                []string           `json:"games"`
	Highres_Image        bool               `json:"highres_image"`
	Illustration_ID      string             `json:"illustartion_id"`
	Image_URIs           ImageURIsObject    `json:"image_uris"`
	Prices               PricesObject       `json:"prices"`
	Printed_Name         string             `json:"printed_name"`
	Printed_Text         string             `json:"printed_text"`
	Printed_Type_Line    string             `json:"printer_type_line"`
	Promo                bool               `json:"promo"`
	Promo_Types          []string           `json:"promo_types"`
	Purchase_URIs        PurchaseURIsObject `json:"purchase_uris"`
	Rarity               string             `json:"rarity"`
	Related_URIs         RelatedURIsObject  `json:"related_uris"`
	Released_At          string             `json:"released_at"`
	Reprint              bool               `json:"reprint"`
	Scryfall_Set_URI     string             `json:"scryfall_set_uri"`
	Set_Name             string             `json:"set_name"`
	Set_Search_URI       string             `json:"set_search_uri"`
	Set_Type             string             `json:"set_type"`
	Set_URI              string             `json:"set_uri"`
	Set                  string             `json:"set"`
	Story_Spotlight      bool               `json:"story_spotlight"`
	Textless             bool               `json:"textless"`
	Variation            bool               `json:"variation"`
	Variation_Of         string             `json:"variation_of"`
	Watermark            string             `json:"watermark"`
	Preview_Previewed_At string             `json:"preview.previewed_at"`
	Preview_Source_URI   string             `json:"preview.source_uri"`
	Preview_Source       string             `json:"preview.source"`
}

type LegalitiesObject struct {
	Standard  string `json:"standard"`
	Future    string `json:"future"`
	Historic  string `json:"historic"`
	Pioneer   string `json:"pioneer"`
	Modern    string `json:"modern"`
	Legacy    string `json:"legacy"`
	Pauper    string `json:"pauper"`
	Vintage   string `json:"vintage"`
	Penny     string `json:"penny"`
	Commander string `json:commander""`
	Brawl     string `json:"brawl"`
	Duel      string `json:"duel"`
	Oldschool string `json:"oldschool"`
}

type ImageURIsObject struct {
	Small       string `json:"small"`
	Normal      string `json:"normal"`
	Large       string `json:"large"`
	PNG         string `json:"png"`
	Art_Crop    string `json:"art_crop"`
	Border_Crop string `json:"border_crop"`
}

type PricesObject struct {
	USD      string `json:"usd"`
	USD_Foil string `json:"usd_foil"`
	EUR      string `json:"eur"`
	Tix      string `json:"tix"`
}

type PurchaseURIsObject struct {
	TCGPlayer   string `json:"tcgplayer"`
	CardMarket  string `json:"cardmarket"`
	CardHoarder string `json:"cardhoarder"`
}

type RelatedURIsObject struct {
	TCGPlayer_Decks string `json:"tcgplayer_decks"`
	EDHRec          string `json:"edhrec"`
	MTGTop8         string `json:"mtgtop8"`
}

type CardFaceObject struct {
	Artist            string          `json:"artist"`
	Color_Indicator   []Color         `json:"color_indicator"`
	Colors            []Color         `json:"colors"`
	Flavor_Text       string          `json:"flavor_text"`
	Illustration_ID   string          `json:"illustration_id"`
	Image_URIs        ImageURIsObject `json:"image_uris"`
	Loyalty           string          `json:"loyalty"`
	Mana_Cost         string          `json:"mana_cost"`
	Name              string          `json:"name"`
	ObjectType        string          `json:"object"` // Always "card_face"
	Oracle_Text       string          `json:"oracle_text"`
	Power             string          `json:"power"`
	Printed_Name      string          `json:"printed_name"`
	Printed_Text      string          `json:"printed_text"`
	Printed_Type_Line string          `json:"printed_type_line"`
	Toughness         string          `json:"toughness"`
	Type_Line         string          `json:"type_line"`
	Watermark         string          `json:"watermark"`
}

type RelatedCardObject struct {
	ID         string `json:"id"`
	ObjectType string `json:"object"`
	Component  string `json:"component"`
	Name       string `json:"name"`
	Type_Line  string `json:"type_line"`
	URI        string `json:"uri"`
}

func pretty_print(blob []byte) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, blob, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(prettyJSON.Bytes()))
	return nil
}

func ParseJSONList(response []byte) (List, error) {
	var blob List
	err := json.Unmarshal(response, &blob)
	if err != nil {
		return List{}, err
	}
	return blob, nil
}

func DownloadFile(path string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Close the Body after this code block has run
	defer resp.Body.Close()

	// Check if the destination filepath exists, and if not, create it
	exists, err := path_exists(path)
	if err != nil {
		return err
	}
	if !exists {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	// Again, close after the block is done
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func path_exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
