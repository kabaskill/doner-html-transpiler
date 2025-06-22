package main

// Dictionary contains the German to HTML translations
type Dictionary struct {
	tags       map[string]string
	attributes map[string]string
}

// NewDictionary creates a new dictionary with German-to-HTML mappings
func NewDictionary() *Dictionary {
	return &Dictionary{
		tags: map[string]string{
			// Basic structure
			"döner":        "html",
			"dokument":     "html",     // alternative for html
			"kopf":         "head",     // head
			"head":         "head",     // allow English too
			"titel":        "title",    // title
			"title":        "title",    // allow English too
			"körper":       "body",     // body
			"body":         "body",     // allow English too
			"meta":         "meta",     // meta
			"beschreibung": "meta",     // description/meta
			"verknüpfung":  "link",     // link
			"stil":         "style",    // style
			// "skript":       "script",   // REMOVED: script tags are dangerous
			
			// Text content
			"überschrift1":      "h1",      // heading 1
			"hauptüberschrift":  "h1",      // main heading
			"überschrift2":      "h2",      // heading 2
			"überschrift3":      "h3",      // heading 3
			"überschrift4":      "h4",      // heading 4
			"überschrift5":      "h5",      // heading 5
			"überschrift6":      "h6",      // heading 6
			"absatz":            "p",       // paragraph
			"p":                 "p",       // allow English too
			"bereich":           "div",     // div
			"spanne":            "span",    // span
			"stark":             "strong",  // strong
			"betont":            "em",      // emphasized
			"fett":              "b",       // bold
			"kursiv":            "i",       // italic
			
			// Lists
			"ungeordnete_liste": "ul", // unordered list
			"liste":             "ul", // list (simple form)
			"geordnete_liste":   "ol", // ordered list
			"listenelement":     "li", // list item
			"li":                "li", // allow English too
			
			// Links and media
			"anker":        "a",       // anchor/link
			"bild":         "img",     // image
			"video":        "video",   // video
			"audio":        "audio",   // audio
			
			// Forms
			"formular":     "form",    // form
			"eingabe":      "input",   // input
			"beschriftung": "label",   // label
			"knopf":        "button",  // button
			"auswahl":      "select",  // select
			"option":       "option",  // option
			"textbereich":  "textarea", // textarea
			
			// Tables
			"tabelle":      "table",   // table
			"tabellenreihe": "tr",     // table row
			"tabellendaten": "td",     // table data
			"tabellenkopf":  "th",     // table header
			"tabellenkörper": "tbody", // table body
			"tabellenheader": "thead", // table head
			"tabellenfuß":   "tfoot",  // table foot
		},
		
		attributes: map[string]string{
			// Common attributes
			"klasse":       "class",    // class
			"identität":    "id",       // id
			"stil":         "style",    // style
			"titel":        "title",    // title
			"sprache":      "lang",     // language
			
			// Link attributes
			"href":         "href",     // href (keeping same)
			"ziel":         "target",   // target
			
			// Image attributes
			"quelle":       "src",      // source
			"alternativ":   "alt",      // alternative text
			"breite":       "width",    // width
			"höhe":         "height",   // height
			
			// Form attributes
			"typ":          "type",     // type
			"name":         "name",     // name (keeping same)
			"wert":         "value",    // value
			"platzhalter":  "placeholder", // placeholder
			"erforderlich": "required", // required
			"deaktiviert":  "disabled", // disabled
			
			// Event attributes
			"bei_klick":    "onclick",  // onclick
			"bei_laden":    "onload",   // onload
			"bei_änderung": "onchange", // onchange
		},
	}
}

// TranslateTag translates a German tag to HTML
func (d *Dictionary) TranslateTag(germanTag string) (string, bool) {
	htmlTag, exists := d.tags[germanTag]
	return htmlTag, exists
}

// TranslateAttribute translates a German attribute to HTML
func (d *Dictionary) TranslateAttribute(germanAttr string) (string, bool) {
	htmlAttr, exists := d.attributes[germanAttr]
	return htmlAttr, exists
}
